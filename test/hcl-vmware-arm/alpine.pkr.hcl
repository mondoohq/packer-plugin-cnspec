# Copyright Mondoo, Inc. 2026
# SPDX-License-Identifier: BUSL-1.1

packer {
  required_plugins {
    vmware = {
      version = ">= 1.0.0"
      source  = "github.com/hashicorp/vmware"
    }
    cnspec = {
      version = ">= 10.0.0"
      source  = "github.com/mondoohq/cnspec"
    }
  }
}

source "vmware-iso" "alpine" {
  iso_urls = [
    "isos/alpine-virt-3.15.0-aarch64.iso",
    "http://dl-cdn.alpinelinux.org/alpine/v3.15/releases/aarch64/alpine-virt-3.15.0-aarch64.iso"
  ]
  iso_checksum =  "sha256:f302cf1b2dbbd0661b8f53b167f24131c781b86ab3ae059654db05cd62d3c39c"

  communicator = "ssh"
  ssh_username = var.sshusername
  ssh_password = var.sshpassword
  shutdown_command = "echo vagrant | sudo -S /sbin/poweroff"

  http_directory      =  "http"

  boot_wait = "10s"
  boot_command = [
    "root<enter><wait>",
    "ifconfig eth0 up && udhcpc -i eth0<enter><wait10>",
    "wget http://{{ .HTTPIP }}:{{ .HTTPPort }}/answers<enter><wait>",
    "setup-alpine -f $PWD/answers<enter><wait5>",
    "${var.rootpassword}<enter><wait>",
    "${var.rootpassword}<enter><wait>",
    "<wait10>y<enter>",
    "<wait10><wait10>",
    "reboot<enter>",
    "<wait10><wait10>",
    "root<enter><wait5>",
    "${var.rootpassword}<enter><wait5>",
    "echo http://dl-cdn.alpinelinux.org/alpine/v3.16/community/ >> /etc/apk/repositories<enter>",
    "apk add sudo<enter><wait5>",
    "echo 'Defaults env_keep += \"http_proxy https_proxy\"' > /etc/sudoers.d/wheel<enter>",
    "echo '%wheel ALL=(ALL) NOPASSWD:ALL' >> /etc/sudoers.d/wheel<enter>",
    "adduser ${var.sshusername}<enter><wait5>",
    "${var.sshpassword}<enter><wait>",
    "${var.sshpassword}<enter><wait>",
    "adduser ${var.sshusername} wheel<enter><wait5>",
    "apk add virtualbox-guest-additions virtualbox-guest-additions-openrc<enter>",
    "<wait10>"
  ]

  guest_os_type = "arm-other5xlinux-64"

  version = "19"
  disk_adapter_type = "nvme"
  network_adapter_type = "e1000e"

  vmx_data = {
    "usb_xhci:4.present" = "TRUE"
    "usb:1.present" = "TRUE"
  }
}

build {
  sources = [
    "source.vmware-iso.alpine",
  ]

  provisioner "shell" {
    scripts = [
        "scripts/prepare.sh"
    ]
  }

  provisioner "cnspec" {
    on_failure =  "continue"

    asset_name = "test-name"
    output = "json"
    risk_threshold = 50

    annotations = {
      name = "Packer Builder"
      custom_key = "custom_value"
    }
  }
}