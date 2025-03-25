#!/usr/bin/env bash

if [ -z "$1" ]; then
  echo "Error: Must specify version to deploy"
  echo "Usage: $0 <version>"
  exit 1
fi

VERSION=$1

mkdir -p /opt/app/versions/${VERSION}
unzip -o /opt/app/versions/${VERSION}.zip -d /opt/app/versions/${VERSION}/
ln -sf /opt/app/versions/${VERSION} /opt/app/current

systemctl restart digital-venue
