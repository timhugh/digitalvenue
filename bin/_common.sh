DV_ENV=${DV_ENV:-development}
DV_DEPLOY_ENV=${DV_DEPLOY_ENV:-production}
DV_DO_IMAGE=${DV_DO_IMAGE:-ubuntu-22-04-x64}
DV_DO_SIZE=${DV_DO_SIZE:-s-1vcpu-512mb-10gb}
DV_DO_REGION=${DV_DO_REGION:-sfo3}
DV_VERSION=$(git rev-parse --short HEAD)
DV_DO_TAG="env:${DV_DEPLOY_ENV}"
DV_REMOTE_ROOT=/opt/digitalvenue
DV_DO_DROPLET_NAME=dv-${DV_DEPLOY_ENV}

get_droplet_ip() {
  doctl compute droplet list --tag-name "${DV_DO_TAG}" --format PublicIPv4 --no-header
}
