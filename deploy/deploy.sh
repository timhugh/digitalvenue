#!/usr/bin/env bash

set -e

VERSION=$1
if [ -z "${VERSION}" ]; then
  echo "Error: Must specify version to deploy"
  echo "Usage: $0 <version>"
  exit 1
fi

dv_root=/opt/digitalvenue
current_root=${dv_root}/current
version_root=${dv_root}/versions/${VERSION}

if [ ! -d ${version_root} ]; then
  echo "${version_root} does not exist!"
  exit 1
fi

echo "deploying ${VERSION}"
rm -f ${current_root}
ln -sf ${version_root} ${current_root}

echo "enabling and starting service"
ln -sf ${current_root}/dvserver.service /etc/systemd/system/dvserver.service
systemctl daemon-reload
systemctl enable dvserver
systemctl restart dvserver
