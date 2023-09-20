# Copyright (c) Mondoo, Inc.
# SPDX-License-Identifier: BUSL-1.1

packer {
  required_version = ">= 1.8.4"
  required_plugins {
    git = {
      version = ">= 0.3.2"
      source  = "github.com/ethanmdavidson/git"
    }
    vsphere = {
      version = ">= v1.1.1"
      source  = "github.com/hashicorp/vsphere"
    }
    cnspec = {
      version = ">= v6.1.4"
      source  = "github.com/mondoohq/cnspec"
    }
  }
}

variable "vsphere_endpoint" {
  type = string
}

variable "vsphere_username" {
  type = string
}

variable "vsphere_password" {
  type = string
}

variable "vsphere_host" {
  type = string
}

variable "vsphere_datacenter" {
  type = string
}

variable "vsphere_datastore" {
  type = string
}

variable "vsphere_network" {
  type        = string
}

locals {
  data_source_content = {
    "/ks.cfg" = templatefile("${abspath(path.root)}/http/kickstart.cfg", {})
  }
  data_source_command = "inst.ks=cdrom:/ks.cfg"
}

source "vsphere-iso" "centos7" {
  CPUs                 = 1
  RAM                  = 1024
  RAM_reserve_all      = true
  boot_command = [
    "<tab>",
    "text ${local.data_source_command}",
    "<enter><wait10><wait10><wait10>",
    "<wait10><wait10><wait10>",
  ]
  # boot_command         = ["<tab> text ks=http://{{ .HTTPIP }}:{{ .HTTPPort }}/kickstart.cfg<enter>", "<wait10><wait10><wait10>", "<wait10><wait10><wait10>"]
  boot_wait            = "10s"
  disk_controller_type = ["pvscsi"]
  guest_os_type        = "centos7_64Guest"
  cd_content           = local.data_source_content

  iso_checksum         = "sha256:d68f92f41ab008f94bd89ec4e2403920538c19a7b35b731e770ce24d66be129a"
  iso_url              = "http://ftp.halifax.rwth-aachen.de/centos/7.9.2009/isos/x86_64/CentOS-7-x86_64-Minimal-2207-02.iso"
    
  vm_name          = "example-centos"
  shutdown_command = "echo 'vagrant'| sudo -S /sbin/poweroff"
  ssh_password     = "vagrant"
  ssh_port         = 22
  ssh_timeout      = "10m"
  ssh_username     = "vagrant"

  network_adapters {
    network      = var.vsphere_network
    network_card = "vmxnet3"
  }

  storage {
    disk_size             = 32768
    disk_thin_provisioned = true
  }
  
  // vCenter Server Endpoint Settings and Credentials
  vcenter_server = "${var.vsphere_endpoint}"
  host           = "${var.vsphere_host}"
  username       = "${var.vsphere_username}"
  password       = "${var.vsphere_password}"
  insecure_connection  = "true"
  datacenter     = "${var.vsphere_datacenter}"
  datastore      = "${var.vsphere_datastore}"
}

build {
  description = "Build CentOS 7 x86_64"

  sources = ["source.vsphere-iso.centos7"]

  provisioner "shell" {
    scripts = ["scripts/prepare.sh"]
  }

  provisioner "cnspec" {
    use_proxy = true
  }
}
