# Copyright 2023 Broadcom. All rights reserved.
# SPDX-License-Identifier: BSD-2

/*
    DESCRIPTION:
    VMware Photon OS 5 build definition.
    Packer Plugin for VMware vSphere: 'vsphere-iso' builder.
*/

//  BLOCK: packer
//  The Packer configuration.

packer {
  required_version = ">= 1.9.4"
  required_plugins {
    vsphere = {
      source  = "github.com/hashicorp/vsphere"
      version = ">= 1.2.1"
    }
    ansible = {
      source  = "github.com/hashicorp/ansible"
      version = ">= 1.1.0"
    }
    git = {
      source  = "github.com/ethanmdavidson/git"
      version = ">= 0.4.3"
    }
  }
}

//  BLOCK: data
//  Defines the data sources.

data "git-repository" "cwd" {}

//  BLOCK: locals
//  Defines the local variables.

locals {
  build_by          = "Built by: HashiCorp Packer ${packer.version}"
  build_date        = formatdate("YYYY-MM-DD hh:mm ZZZ", timestamp())
  build_version     = data.git-repository.cwd.head
  build_description = "Version: ${local.build_version}\nBuilt on: ${local.build_date}\n${local.build_by}"
  iso_paths         = ["[${var.common_iso_datastore}] ${var.iso_path}/${var.iso_file}"]
  iso_checksum      = "${var.iso_checksum_type}:${var.iso_checksum_value}"
  manifest_date     = formatdate("YYYY-MM-DD'T'hhmmss'Z'", timestamp())
  manifest_path     = "${path.cwd}/manifests/"
  manifest_output   = "${local.manifest_path}${local.manifest_date}.json"
  ovf_export_path   = "${path.cwd}/artifacts/${local.vm_name}"
  data_source_content = {
    "/ks.json" = templatefile("${abspath(path.root)}/data/ks.pkrtpl.hcl", {
      build_username           = var.build_username
      build_password           = var.build_password
      build_password_encrypted = var.build_password_encrypted
    })
  }
  data_source_command = var.common_data_source == "http" ? "ks=http://{{ .HTTPIP }}:{{ .HTTPPort }}/ks.json" : "ks=/dev/sr2:/ks.json"
  vm_name             = "${var.vm_name}"
  bucket_name         = replace("${var.vm_guest_os_family}-${var.vm_guest_os_name}-${var.vm_guest_os_version}", ".", "")
  bucket_description  = "${var.vm_guest_os_family} ${var.vm_guest_os_name} ${var.vm_guest_os_version}"
}

//  BLOCK: source
//  Defines the builder configuration blocks.

source "vsphere-iso" "linux-photon" {

  // vCenter Server Endpoint Settings and Credentials
  vcenter_server      = var.vsphere_endpoint
  username            = var.vsphere_username
  password            = var.vsphere_password
  insecure_connection = var.vsphere_insecure_connection

  // vSphere Settings
  datacenter                     = var.vsphere_datacenter
  cluster                        = var.vsphere_cluster
  host                           = var.vsphere_host
  datastore                      = var.vsphere_datastore
  folder                         = var.vsphere_folder
  resource_pool                  = var.vsphere_resource_pool
  set_host_for_datastore_uploads = var.vsphere_set_host_for_datastore_uploads

  // Virtual Machine Settings
  vm_name              = local.vm_name
  guest_os_type        = var.vm_guest_os_type
  firmware             = var.vm_firmware
  CPUs                 = var.vm_cpu_count
  cpu_cores            = var.vm_cpu_cores
  CPU_hot_plug         = var.vm_cpu_hot_add
  RAM                  = var.vm_mem_size
  RAM_hot_plug         = var.vm_mem_hot_add
  cdrom_type           = var.vm_cdrom_type
  disk_controller_type = var.vm_disk_controller_type
  storage {
    disk_size             = var.vm_disk_size
    disk_thin_provisioned = var.vm_disk_thin_provisioned
  }
  network_adapters {
    network      = var.vsphere_network
    network_card = var.vm_network_card
  }
  vm_version           = var.common_vm_version
  remove_cdrom         = var.common_remove_cdrom
  tools_upgrade_policy = var.common_tools_upgrade_policy
  notes                = local.build_description

  // Removable Media Settings
  iso_url      = var.iso_url
  iso_paths    = local.iso_paths
  iso_checksum = local.iso_checksum
  http_content = var.common_data_source == "http" ? local.data_source_content : null
  cd_content   = var.common_data_source == "disk" ? local.data_source_content : null

  // Boot and Provisioning Settings
  http_ip       = var.common_http_ip
  http_port_min = var.common_http_port_min
  http_port_max = var.common_http_port_max
  boot_order    = var.vm_boot_order
  boot_wait     = var.vm_boot_wait
  boot_command = [
    // This sends the "escape" key, waits, and then sends the "c" key. In the GRUB boot loader, this is used to enter command line mode.
    "<esc><wait>c",
    // This types a command to load the Linux kernel from the specified path, with the specified boot parameters.
    // The 'data_source_command' local variable is used to specify the kickstart data source configured in the common variables.
    "linux /isolinux/vmlinuz root=/dev/ram0 loglevel=3 insecure_installation=1 ${local.data_source_command} photon.media=cdrom",
    // This sends the "enter" key, which executes the command.
    "<enter>",
    // This types a command to load the initial RAM disk from the specified path.
    "initrd /isolinux/initrd.img",
    // This sends the "enter" key, which executes the command.
    "<enter>",
    // This types the "boot" command, which starts the boot process using the loaded kernel and initial RAM disk.
    "boot",
    // This sends the "enter" key, which executes the command.
    "<enter>"
  ]
  ip_wait_timeout   = var.common_ip_wait_timeout
  ip_settle_timeout = var.common_ip_settle_timeout
  shutdown_command  = "echo '${var.build_password}' | sudo -S -E shutdown -P now"
  shutdown_timeout  = var.common_shutdown_timeout

  // Communicator Settings and Credentials
  communicator       = "ssh"
  ssh_proxy_host     = var.communicator_proxy_host
  ssh_proxy_port     = var.communicator_proxy_port
  ssh_proxy_username = var.communicator_proxy_username
  ssh_proxy_password = var.communicator_proxy_password
  ssh_username       = var.build_username
  ssh_password       = var.build_password
  ssh_port           = var.communicator_port
  ssh_timeout        = var.communicator_timeout

  // Template and Content Library Settings
  convert_to_template = var.common_template_conversion
  dynamic "content_library_destination" {
    for_each = var.common_content_library_name != null ? [1] : []
    content {
      library     = var.common_content_library_name
      description = local.build_description
      ovf         = var.common_content_library_ovf
      destroy     = var.common_content_library_destroy
      skip_import = var.common_content_library_skip_export
    }
  }

  // OVF Export Settings
  dynamic "export" {
    for_each = var.common_ovf_export_enabled == true ? [1] : []
    content {
      name  = local.vm_name
      force = var.common_ovf_export_overwrite
      options = [
        "extraconfig"
      ]
      output_directory = local.ovf_export_path
    }
  }
}

//  BLOCK: build
//  Defines the builders to run, provisioners, and post-processors.

build {
  sources = ["source.vsphere-iso.linux-photon"]

  provisioner "shell" {
    inline = [
      "echo -e 'Starting patch cycle...'",
      "sudo tdnf update -y"
      ]
  }

  provisioner "shell" {
    inline = [
      # Amend the params set by default
      "sudo sed -i 's/^ClientAliveCountMax.*/ClientAliveCountMax 3/g' /etc/ssh/sshd_config",
      # Create a ssh banner
      "echo -e '\n* This system is for the use of authorized users only. *\n' | sudo tee /etc/issue.net > /dev/null",
      # Add additional sshd config
      "echo -e '\n#\n# Added via Packer...\n#' | sudo tee -a /etc/ssh/sshd_config > /dev/null",
      "echo 'KexAlgorithms curve25519-sha256,curve25519-sha256@libssh.org,diffie-hellman-group14-sha256,diffie-hellman-group16-sha512,diffie-hellman-group18-sha512,ecdh-sha2-nistp521,ecdh-sha2-nistp384,ecdh-sha2-nistp256,diffie-hellman-group-exchange-sha256' | sudo tee -a /etc/ssh/sshd_config > /dev/null",
      "echo 'Ciphers chacha20-poly1305@openssh.com,aes256-gcm@openssh.com,aes128-gcm@openssh.com,aes256-ctr,aes192-ctr,aes128-ctr' | sudo tee -a /etc/ssh/sshd_config > /dev/null",
      "echo 'MACs hmac-sha2-512-etm@openssh.com,hmac-sha2-256-etm@openssh.com,hmac-sha2-512,hmac-sha2-256' | sudo tee -a /etc/ssh/sshd_config > /dev/null",
      "echo 'PermitUserEnvironment no' | sudo tee -a /etc/ssh/sshd_config > /dev/null",
      "echo 'PermitEmptyPasswords no' | sudo tee -a /etc/ssh/sshd_config > /dev/null",
      "echo 'HostbasedAuthentication no' | sudo tee -a /etc/ssh/sshd_config > /dev/null",
      "echo 'maxstartups 10:30:60' | sudo tee -a /etc/ssh/sshd_config > /dev/null",
      "echo 'LogLevel INFO' | sudo tee -a /etc/ssh/sshd_config > /dev/null",
      "echo 'IgnoreRhosts yes' | sudo tee -a /etc/ssh/sshd_config > /dev/null",
      "echo 'DenyUsers root' | sudo tee -a /etc/ssh/sshd_config > /dev/null",
      "echo 'LoginGraceTime 60' | sudo tee -a /etc/ssh/sshd_config > /dev/null",
      "echo 'MaxSessions 4' | sudo tee -a /etc/ssh/sshd_config > /dev/null",
      "echo 'MaxAuthTries 4' | sudo tee -a /etc/ssh/sshd_config > /dev/null",
      "echo 'ClientAliveInterval 15' | sudo tee -a /etc/ssh/sshd_config > /dev/null",
      "echo 'Banner /etc/issue.net' | sudo tee -a /etc/ssh/sshd_config > /dev/null",
    ]
  }

  provisioner "shell" {
    expect_disconnect = true
    inline = [
      "sudo reboot now",
    ]
    pause_after  = "10s"
  }

  provisioner "cnspec" {
    #on_failure      = "continue"
    #risk_threshold = 85
    #mondoo_config_path = "/Path/To/Mondoo/config.yml"
    asset_name         = local.vm_name
    debug = true
    sudo {
      active = true
    }

    annotations = {
      build_date = local.build_date
      os_family  = var.vm_guest_os_family
      os_name    = var.vm_guest_os_name
      os_version = var.vm_guest_os_version
    }
  }

  post-processor "manifest" {
    output     = local.manifest_output
    strip_path = true
    strip_time = true
    custom_data = {
      build_username           = var.build_username
      build_date               = local.build_date
      build_version            = local.build_version
      common_data_source       = var.common_data_source
      common_vm_version        = var.common_vm_version
      vm_cpu_cores             = var.vm_cpu_cores
      vm_cpu_count             = var.vm_cpu_count
      vm_disk_size             = var.vm_disk_size
      vm_disk_thin_provisioned = var.vm_disk_thin_provisioned
      vm_firmware              = var.vm_firmware
      vm_guest_os_type         = var.vm_guest_os_type
      vm_mem_size              = var.vm_mem_size
      vm_network_card          = var.vm_network_card
      vsphere_cluster          = var.vsphere_cluster
      vsphere_host             = var.vsphere_host
      vsphere_datacenter       = var.vsphere_datacenter
      vsphere_datastore        = var.vsphere_datastore
      vsphere_endpoint         = var.vsphere_endpoint
      vsphere_folder           = var.vsphere_folder
    }
  }
}