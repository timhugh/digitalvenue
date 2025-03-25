#!/usr/bin/env bash

set -e

if [ -z "$1" ]; then
  echo "Error: Must specify version to deploy"
  echo "Usage: $0 <version>"
  exit 1
fi

VERSION=$1

mkdir -p /opt/dvserver/versions/${VERSION}
unzip -o /opt/dvserver/versions/${VERSION}.zip -d /opt/dvserver/versions/${VERSION}/
ln -sf /opt/dvserver/versions/${VERSION} /opt/dvserver/current

systemctl restart digital-venue
