#!/usr/bin/env bash

# Install java
apt-get update && apt-get install -y openjdk-21-jre-headless

# Install service definition for digital-venue server
cp /opt/digital-venue.service /etc/systemd/system/digital-venue.service
systemctl enable digital-venue.service

# Create directory for deployments
mkdir -p /opt/digital-venue/versions
