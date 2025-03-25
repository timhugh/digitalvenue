#!/usr/bin/env bash

if [ -z "$1" ]; then
  echo "Error: Must specify version to deploy"
  echo "Usage: $0 <version>"
  exit 1
fi

VERSION=$1

systemctl stop digital-venue || true

ln -sf /opt/app/server-${VERSION}.jar /opt/app/server.jar
chmod 755 /opt/app/server.jar

systemctl daemon-reload
systemctl start digital-venue
systemctl enable digital-venue

systemctl status digital-venue
