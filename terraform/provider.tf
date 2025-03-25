terraform {
  required_providers {
    digitalocean = {
      source  = "digitalocean/digitalocean"
      version = "~> 2.0"
    }
  }
}

variable "digitalocean_token" {
  description = "The DigitalOcean API token"
  type        = string
}

variable "ssh_key_path" {
  description = "The path to the private key file"
  type        = string
}

variable "ssh_key_name" {
  description = "The name of the SSH key in DigitalOcean"
  type        = string
}

provider "digitalocean" {
  token = var.digitalocean_token
}

data "digitalocean_ssh_key" "DigitalOceanKey" {
  name = var.ssh_key_name
}
