# Copyright (c) Mondoo, Inc.
# SPDX-License-Identifier: BUSL-1.1

packer {
  required_plugins {
    amazon = {
      version = ">= 1.1.0"
      source  = "github.com/hashicorp/amazon"
    }
    cnspec = {
      version = ">= 9.0.0"
      source  = "github.com/mondoohq/cnspec"
    }
  }
}

variable "aws_region" {
  default = "us-east-1"
  type    = string
}

variable "image_prefix" {
  type = string
  description = "Prefix to be applied to image name"
  default = "mondoo-windows2019-secure-base"
}

locals { timestamp = regex_replace(timestamp(), "[- TZ:]", "") }

source "amazon-ebs" "windows2019" {
  ami_name      = "${var.image_prefix}-${local.timestamp}"
  communicator  = "winrm"
  instance_type = "t2.micro"
  region        = var.aws_region
  source_ami_filter {
    filters = {
      name                = "Windows_Server-2019-English-Full-Base-*"
      root-device-type    = "ebs"
      virtualization-type = "hvm"
    }
    most_recent = true
    owners      = ["801119661308"]
  }
  user_data_file = "../scripts/bootstrap_win.txt"
  winrm_password = "SuperS3cr3t!!!!"
  winrm_username = "Administrator"
}

build {
  name    = "${var.image_prefix}-${local.timestamp}"
  sources = ["source.amazon-ebs.windows2019"]

  provisioner "cnspec" {
    on_failure = "continue"
    asset_name = "${var.image_prefix}-${local.timestamp}"
    annotations = {
      Source_AMI    = "{{ .SourceAMI }}"
      Creation_Date = "{{ .SourceAMICreationDate }}"
    }
  }
}