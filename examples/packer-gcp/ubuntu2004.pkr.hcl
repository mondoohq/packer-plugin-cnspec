packer {
  required_plugins {
    googlecompute = {
      version = ">= 1.0.0"
      source  = "github.com/hashicorp/googlecompute"
    }
    mondoo = {
      version = ">= 0.4.0"
      source  = "github.com/mondoohq/mondoo"
    }
  }
}

variable "zone" {
  default = "us-east5-a"
}

variable "mondoo_config_path" {
  type = string
  description = "The path to the config to be used when scanning"
  default = ""
}

variable "project_id" {
  type        = string
  description = "The project ID that will be used to launch instances and store images"
}

variable "image_prefix" {
  type        = string
  description = "Prefix to be applied to image name"
  default     = "mondoo-gcp-ubuntu-2004-secure-base"
}

locals { timestamp = regex_replace(timestamp(), "[- TZ:]", "") }

source "googlecompute" "ubuntu2004" {
  image_name              = "${var.image_prefix}-${local.timestamp}"
  machine_type            = "e2-small"
  source_image            = "ubuntu-pro-2004-focal-v20220627a"
  ssh_username            = "packer"
  temporary_key_pair_type = "rsa"
  temporary_key_pair_bits = 2048
  zone                    = var.zone
  project_id              = var.project_id
}

build {
  name = "mondoo-gcp-ubuntu-2004-secure-base"

  sources = ["source.googlecompute.ubuntu2004"]
  provisioner "shell" {
    inline = [
      "sudo hostnamectl set-hostname ${var.image_prefix}-${local.timestamp}",
    ]
  }

  provisioner "mondoo" {
    on_failure = "continue"
    asset_name = "${var.image_prefix}-${local.timestamp}"
    mondoo_config_path = "${var.mondoo_config_path}"
    annotations = {
      Name          = "${var.image_prefix}-${local.timestamp}"
    }
  }
}
