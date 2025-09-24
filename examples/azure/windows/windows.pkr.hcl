# Copyright (c) Mondoo, Inc.
# SPDX-License-Identifier: BUSL-1.1

packer {
  required_plugins {
    azure = {
      source  = "github.com/hashicorp/azure"
      version = ">= 2"
    }
    cnspec = {
      version = "~> 12"
      source  = "github.com/mondoohq/cnspec"
    }
  }
}

locals {
  random = uuidv4()
  date   = timestamp()
}

variable "location" {
  type = string
  description = "The Azure region to deploy to"
}

variable "resourceGroup" {
  type = string
  description = "The Azure resource group to deploy to"
}

variable "imageName" {
  type = string
  description = "The Azure Shared Image Gallery image name"
}

variable "imageVersion" {
  type = string
  description = "The Azure Shared Image Gallery image version"
}

source "azure-arm" "windows" {
  use_azure_cli_auth = true

  os_type         = "Windows"
  image_publisher = "MicrosoftWindowsServer"
  image_offer     = "WindowsServer"
  image_sku       = "2019-Datacenter"

  azure_tags = {
    packer   = "true",
    build-id = "${local.random}"
  }

  managed_image_name = "${var.imageName}-${var.imageVersion}"
  managed_image_resource_group_name = var.resourceGroup

  location = var.location
  vm_size  = "Standard_B4ms"

  communicator   = "winrm"
  winrm_use_ssl  = "true"
  winrm_insecure = "true"
  winrm_timeout  = "50m"
  winrm_username = "packer"
}

build {

  sources = ["sources.azure-arm.windows"]

  provisioner "cnspec" {
    asset_name      = "${var.imageName}-${var.imageVersion}"
    # score_threshold = 80
    on_failure      = "continue"
    debug           = false
    annotations     = {
      os-type       = "WindowsServer"
      os-version    = "2019-Datacenter"
      image-version = "${var.imageVersion}"
      build-time    = "${local.date}"
      build-id      = "${local.random}"
    }
  }

  provisioner "powershell" {
    inline = [
      "# If Guest Agent services are installed, make sure that they have started.",
      "foreach ($service in Get-Service -Name RdAgent, WindowsAzureTelemetryService, WindowsAzureGuestAgent -ErrorAction SilentlyContinue) { while ((Get-Service $service.Name).Status -ne 'Running') { Start-Sleep -s 5 } }",

      "& $env:SystemRoot\\System32\\Sysprep\\Sysprep.exe /oobe /generalize /quiet /quit /mode:vm",
      "while($true) { $imageState = Get-ItemProperty HKLM:\\SOFTWARE\\Microsoft\\Windows\\CurrentVersion\\Setup\\State | Select ImageState; if($imageState.ImageState -ne 'IMAGE_STATE_GENERALIZE_RESEAL_TO_OOBE') { Write-Output $imageState.ImageState; Start-Sleep -s 10  } else { break } }"
    ]
  }


}
