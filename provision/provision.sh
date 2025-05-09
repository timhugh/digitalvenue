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

# Wait for cloud-init to finish
cloud-init status --wait || true

# create swapfile
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

# Install nginx
apt-get update
apt-get install -y nginx

# Install nginx configuration
envsubst '${DOMAIN_NAME} ${IP_ADDRESS}' < /opt/digitalvenue/nginx.conf.template > /etc/nginx/sites-available/digital-venue
ln -sf /etc/nginx/sites-available/digital-venue /etc/nginx/sites-enabled/digital-venue
rm -f /etc/nginx/sites-enabled/default
nginx -t
systemctl restart nginx

# TODO: certbot

# Reload systemd daemon to pick up changes
systemctl daemon-reload
