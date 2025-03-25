#!/usr/bin/env bash

set -e

required_environment_variables="DOMAIN_NAME IP_ADDRESS APP_VERSION"
for variable in $required_environment_variables; do
  if [ -z "${!variable}" ]; then
    echo "Error: $variable is not set"
    exit 1
  fi
done

# Wait for cloud-init to finish
cloud-init status --wait || true

# Install java
apt-get update
apt-get install -y \
  openjdk-21-jre-headless \
  nginx

# Set up environment variables
envsubst < /opt/env.template > /opt/dvserver/.env

# Install nginx configuration
envsubst < /opt/nginx.conf.template > /etc/nginx/sites-available/digital-venue
ln -s /etc/nginx/sites-available/digital-venue /etc/nginx/sites-enabled/digital-venue
rm /etc/nginx/sites-enabled/default
nginx -t
systemctl restart nginx

# Install service definition for digital-venue server
cp /opt/digital-venue.service /etc/systemd/system/digital-venue.service
systemctl enable digital-venue.service

# Create directory for deployments
mkdir -p /opt/dvserver/versions
