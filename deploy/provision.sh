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

# Create directory for deployments
mkdir -p /opt/dvserver/versions

# Set up environment variables
envsubst < /opt/env.template > /opt/dvserver/.env

# Wait for cloud-init to finish
cloud-init status --wait || true

# Install nginx
apt-get update
apt-get install -y nginx

# Install nginx configuration
envsubst '${DOMAIN_NAME} ${IP_ADDRESS}' < /opt/nginx.conf.template > /etc/nginx/sites-available/digital-venue
ln -sf /etc/nginx/sites-available/digital-venue /etc/nginx/sites-enabled/digital-venue
rm -f /etc/nginx/sites-enabled/default
nginx -t
systemctl restart nginx

# Install service definition for digital-venue server
cp /opt/digital-venue.service /etc/systemd/system/digital-venue.service
systemctl enable digital-venue.service

# Reload systemd daemon to pick up changes
systemctl daemon-reload
