# Copyright Mondoo, Inc. 2026
# SPDX-License-Identifier: BUSL-1.1

# End-to-end test: runs the locally built cnspec provisioner inside Packer
# against a docker container. Intentionally no `on_failure = "continue"`, so
# any scan error fails the build.
#
# The plugin binary must be built with version 1.0.0 and installed via
#   packer plugins install --path <binary> github.com/mondoohq/cnspec

packer {
  required_plugins {
    docker = {
      version = ">= 1.0.0"
      source  = "github.com/hashicorp/docker"
    }
    cnspec = {
      version = ">= 1.0.0"
      source  = "github.com/mondoohq/cnspec"
    }
  }
}

source "docker" "ubuntu" {
  image   = "ubuntu:24.04"
  discard = true
}

build {
  name    = "e2e"
  sources = ["source.docker.ubuntu"]

  provisioner "cnspec" {
    asset_name   = "packer-e2e-ubuntu"
    policybundle = "${path.root}/policy.mql.yaml"
    output       = "compact"

    annotations = {
      name = "Packer E2E"
    }
  }
}
