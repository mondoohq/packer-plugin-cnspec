# Copyright (c) Mondoo, Inc.
# SPDX-License-Identifier: BUSL-1.1

packer {
  required_plugins {
    docker = {
      version = ">= 0.0.7"
      source  = "github.com/hashicorp/docker"
    }
    cnspec = {
      version = ">= 10.0.0"
      source  = "github.com/mondoohq/cnspec"
    }
  }
}

source "docker" "ubuntu" {
  image  = "ubuntu:xenial"
  commit = true
}

build {
  name = "learn-packer"
  sources = [
    "source.docker.ubuntu"
  ]

  provisioner "cnspec" {
    on_failure =  "continue"

    asset_name = "test-name"
    output = "compact"

    annotations = {
      name = "Packer Builder"
      custom_key = "custom_value"
    }
  }
}
