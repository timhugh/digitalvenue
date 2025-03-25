variable "app_environment" {
  type = string
}

resource "digitalocean_droplet" "dv" {
  image  = "ubuntu-24-04-x64"
  name   = "digital-venue-${var.app_environment}"
  region = "sfo3"
  size   = "s-1vcpu-512mb-10gb"
  ssh_keys = [
    data.digitalocean_ssh_key.DigitalOceanKey.id
  ]
  tags = ["digital-venue", "env:${var.app_environment}"]

  connection {
    host        = self.ipv4_address
    user        = "root"
    type        = "ssh"
    private_key = file(var.ssh_key_path)
    timeout     = "2m"
  }

  user_data = <<-EOF
#!/usr/bin/env bash

mkdir -p /opt/app,
mkdir -p /opt/app/versions,
apt-get update && \
  apt-get install -y openjdk-21-jre-headless
EOF

  provisioner "file" {
    source      = "digital-venue.service"
    destination = "/etc/systemd/system/digital-venue.service"
  }

  provisioner "file" {
    source      = "deploy.sh"
    destination = "/opt/app/deploy.sh"
  }
}

output "droplet_ip" {
  value = digitalocean_droplet.dv.ipv4_address
}
