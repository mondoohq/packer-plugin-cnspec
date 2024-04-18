# Copyright (c) Mondoo, Inc.
# SPDX-License-Identifier: BUSL-1.1

packer {
  required_plugins {
    docker = {
      version = ">= 0.0.7"
      source  = "github.com/hashicorp/docker"
    }
    cnspec = {
      version = ">= 11.0.0"
      source  = "github.com/mondoohq/cnspec"
    }
  }
}

variable "image_prefix" {
  type        = string
  description = "Prefix to be applied to image name"
  default     = "mondoo-ubuntu-2004-secure-base"
}

locals { timestamp = regex_replace(timestamp(), "[- TZ:]", "") }

source "docker" "ubuntu" {
  image  = "ubuntu:jammy"
  commit = true
}

build {
  name    = "mondoo-docker-ubuntu-2004-secure-base"
  sources = [
    "source.docker.ubuntu"
  ]

  provisioner "shell" {
    inline = [
      "echo \"${var.image_prefix}-${local.timestamp}\" > /etc/hostname",
    ]
  }

  provisioner "cnspec" {
    on_failure  = "continue"
    asset_name  = "${var.image_prefix}-${local.timestamp}"
    annotations = {
      Name              = "${var.image_prefix}-${local.timestamp}"
      Description       = "Mondoo Secure Ubuntu 20.04 Base Image"
      # see https://developer.hashicorp.com/packer/docs/templates/hcl_templates/contextual-variables#build-variables
      Build             = "${ build.PackerRunUUID }"
      Source            = "${ source.name }"
      SourceImageDigest = "${ build.SourceImageDigest }"
    }
    output        = "junit"
    output_target = "test-results.xml"
  }
}
