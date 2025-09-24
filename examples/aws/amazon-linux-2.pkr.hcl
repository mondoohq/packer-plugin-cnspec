# Copyright (c) Mondoo, Inc.
# SPDX-License-Identifier: BUSL-1.1


packer {
  required_plugins {
    amazon = {
      version = "~> 1"
      source  = "github.com/hashicorp/amazon"
    }
    cnspec = {
      version = "~> 12"
      source  = "github.com/mondoohq/cnspec"
    }
  }
}

variable "aws_region" {
  default = "us-east-1"
  type    = string
}

variable "image_prefix" {
  type        = string
  description = "Prefix to be applied to image name"
  default     = "mondoo-amazon-linux-2-secure-base"
}

locals {
  timestamp = formatdate("YYYYMMDDhhmmss", timestamp())
}

source "amazon-ebs" "amazon2" {
  ami_name      = "${var.image_prefix}-${local.timestamp}"
  instance_type = "t2.micro"
  region        = var.aws_region
  source_ami_filter {
    filters = {
      name                = "amzn2-ami-kernel-5.*-x86_64-gp2"
      root-device-type    = "ebs"
      virtualization-type = "hvm"
    }
    most_recent = true
    owners      = ["137112412989"]
  }
  ssh_username = "ec2-user"
  tags         = {
    Name          = "${var.image_prefix}-${local.timestamp}"
    Source_AMI    = "{{ .SourceAMI }}"
    Base_AMI_Name = "{{ .SourceAMIName }}"
    Creation_Date = "{{ .SourceAMICreationDate }}"
  }
}

build {
  name = "${var.image_prefix}-${local.timestamp}"

  sources = [
    "source.amazon-ebs.amazon2"
  ]

  provisioner "shell" {
    inline = [
      "sudo hostnamectl set-hostname ${var.image_prefix}-${local.timestamp}",
    ]
  }

  provisioner "cnspec" {
    on_failure = "continue"
    asset_name = "${var.image_prefix}-${local.timestamp}"
    sudo {
      active = true
    }
    annotations = {
      Name          = "${var.image_prefix}-${local.timestamp}"
      Base_AMI_Name = "${ build.SourceAMIName }"
      Source_AMI    = "${ build.SourceAMI }"
      Creation_Date = "${ build.SourceAMICreationDate }"
    }
  }
}

