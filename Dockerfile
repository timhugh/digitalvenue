# syntax=docker/dockerfile:1.4

FROM ubuntu:24.04 AS builder

RUN apt-get update && apt-get install -y \
  build-essential \
  git \
  libssl-dev \
  ninja-build \
  cmake

RUN  echo "GCC version: $(gcc --version | head -n 1)" && \
  echo "CMake version: $(cmake --version | head -n 1)" && \
  echo "Ninja version: $(ninja --version)"

WORKDIR /app
COPY CMakeLists.txt CMakePresets.json ./
COPY cmake ./cmake
COPY common ./common
COPY server ./server

RUN cmake --preset=release
RUN cmake --build --preset=release
RUN ctest --preset=release
RUN cpack --preset=release

FROM ubuntu:24.04 AS runner

# To mimic the environment of a digitalocean droplet:
#
# 1. We need to install gettext so the provision script can use envsubst
#
# 2. We need to stub systemctl and cloud-init so the provision and deploy scripts
#    can call them without failing
#
# We're also going to install nginx so the provision script won't have to do it
# every time the container starts

RUN apt-get update && \
  apt-get install -y \
  gettext \
  nginx

RUN echo "#!/bin/bash" > /usr/local/bin/systemctl && \
  chmod +x /usr/local/bin/systemctl
RUN echo "#!/bin/bash" > /usr/local/bin/cloud-init && \
  chmod +x /usr/local/bin/cloud-init

# Next we mimic the SCP step of our deployment process by placing
# the release tarball in the versions directory
RUN mkdir -p /opt/dvserver/versions
COPY deploy/* /opt/
COPY --from=builder /app/build/release.tar.gz /opt/dvserver/versions/release.tar.gz

# Finally we can create our entrypoint script, which will mimic the SSH
# calls in the deployment process
RUN <<EOF cat > /opt/entrypoint.sh
#!/bin/bash

set -e

DOMAIN_NAME=example.com \
IP_ADDRESS=192.168.1.1 \
APP_VERSION=release \
ENVIRONMENT=docker \
bash /opt/provision.sh

bash /opt/deploy.sh release

"\$@"
EOF
RUN chmod +x /opt/entrypoint.sh

ENTRYPOINT ["/opt/entrypoint.sh"]
CMD ["/opt/dvserver/current/bin/digitalvenue_server"]
