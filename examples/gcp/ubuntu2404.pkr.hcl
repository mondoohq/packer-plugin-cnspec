# Copyright Mondoo, Inc. 2026
# SPDX-License-Identifier: BUSL-1.1

packer {
  required_plugins {
    googlecompute = {
      version = ">= 1.0.0"
      source  = "github.com/hashicorp/googlecompute"
    }
    cnspec = {
      version = "~> 12"
      source  = "github.com/mondoohq/cnspec"
    }
  }
}

variable "zone" {
  default = "europe-west1-d"
}

variable "project_id" {
  type        = string
  description = "The project ID that will be used to launch instances and store images"
}

variable "image_prefix" {
  type        = string
  description = "Prefix to be applied to image name"
  default     = "mondoo-gcp-ubuntu-2404-secure-base"
}

locals { timestamp = regex_replace(timestamp(), "[- TZ:]", "") }

source "googlecompute" "ubuntu2404" {
  image_name              = "${var.image_prefix}-${local.timestamp}"
  machine_type            = "e2-small"
  source_image            = "ubuntu-2404-noble-amd64-v20250624"
  ssh_username            = "packer"
  temporary_key_pair_type = "rsa"
  temporary_key_pair_bits = 2048
  zone                    = var.zone
  project_id              = var.project_id
}

build {
  name = "mondoo-gcp-ubuntu-2404-secure-base"

  sources = ["source.googlecompute.ubuntu2404"]
  provisioner "shell" {
    inline = [
      "sudo hostnamectl set-hostname ${var.image_prefix}-${local.timestamp}",
    ]
  }

  provisioner "cnspec" {
#    on_failure = "continue"
    asset_name = "${var.image_prefix}-${local.timestamp}"
    annotations = {
      Name          = "${var.image_prefix}-${local.timestamp}"
    }
  }
}
