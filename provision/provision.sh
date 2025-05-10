#!/usr/bin/env bash

# NOTE: this script is used to provision new servers but it is also
# run again on every deployment, so everything in here must be idempotent.

set -e

required_environment_variables="DOMAIN_NAME IP_ADDRESS APP_VERSION ENVIRONMENT"
for variable in $required_environment_variables; do
  if [ -z "${!variable}" ]; then
    echo "Error: $variable is not set"
    exit 1
  fi
done

echo "waiting for cloud-init to complete"
cloud-init status --wait || true

if swapon --show | grep -q '/swapfile'; then
  echo "swap file already exists; skipping"
else
  echo "creating swap file"
  fallocate -l 4G /swapfile
  chmod 600 /swapfile
  mkswap /swapfile
  swapon /swapfile
  echo '/swapfile none swap sw 0 0' >> /etc/fstab
fi

echo "configuring firewall"
ufw allow http
ufw allow https
ufw allow ssh
ufw --force enable
ufw status

echo "installing apt packages"
apt-get update
apt-get install -y \
  nginx \
  python3-certbot-nginx

echo "setting up users"
# TODO: this assumes all failures are because the user already exists. It would be better to explicitly check before creating
# so we aren't swallowing other potential errors
adduser --disabled-password \
        --disabled-login \
        --home /opt/digitalvenue \
        --gecos "" \
        digitalvenue || true
usermod -aG digitalvenue www-data

echo "installing bashrc for digitalvenue user"
ln -sf /opt/digitalvenue/digitalvenue.bashrc /opt/digitalvenue/.bashrc

echo "setting up file permissions"
chown -R digitalvenue:digitalvenue /opt/digitalvenue
find /opt/digitalvenue -type d -exec chmod 750 {} \;
find /opt/digitalvenue -type f -exec chmod 640 {} \;
[ -d /opt/digitalvenue/current ] && find /opt/digitalvenue/current/bin -type f -exec chmod 750 {} \;

echo "installing nginx configuration"
envsubst '${DOMAIN_NAME} ${IP_ADDRESS}' < /opt/digitalvenue/nginx.conf.template > /etc/nginx/sites-available/digital-venue
ln -sf /etc/nginx/sites-available/digital-venue /etc/nginx/sites-enabled/digital-venue
rm -f /etc/nginx/sites-enabled/default
nginx -t

echo "generating and installing SSL cert"
certbot \
  --nginx \
  -d "${DOMAIN_NAME}" \
  --non-interactive \
  --agree-tos \
  --email "info@digital-venue.net"
echo "verifying nginx config"
nginx -t

echo "restarting nginx"
systemctl restart nginx
