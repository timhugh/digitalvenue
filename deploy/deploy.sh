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

echo "setting ${VERSION} as current"
rm -f ${current_root}
ln -sf ${version_root} ${current_root}
[ -f ${dv_root}/VERSION ] && cat ${dv_root}/VERSION > ${dv_root}/LAST_VERSION
echo "${VERSION}" > ${dv_root}/VERSION

echo "setting file permissions"
chown -R digitalvenue:digitalvenue /opt/digitalvenue
find /opt/digitalvenue -type d -exec chmod 750 {} \;
find /opt/digitalvenue -type f -exec chmod 640 {} \;
find /opt/digitalvenue/current/bin -type f -exec chmod 750 {} \;

echo "enabling and starting service"
ln -sf ${current_root}/dvserver.service /etc/systemd/system/dvserver.service
systemctl daemon-reload
systemctl enable dvserver
systemctl restart dvserver
