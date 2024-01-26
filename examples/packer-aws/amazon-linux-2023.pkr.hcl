# Copyright (c) Mondoo, Inc.
# SPDX-License-Identifier: BUSL-1.1


packer {
  required_plugins {
    amazon = {
      version = ">= 1.2.0"
      source  = "github.com/hashicorp/amazon"
    }
    cnspec = {
      version = ">= 10.0.0"
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
  default = "mondoo-amazon-linux-2023-secure-base"
}

locals { timestamp = regex_replace(timestamp(), "[- TZ:]", "") }

source "amazon-ebs" "amazon2" {
  ami_name      = "${var.image_prefix}-${local.timestamp}"
  instance_type = "t2.micro"
  region        = var.aws_region
  source_ami_filter {
    filters = {
      name                = "al2023-ami-2023.3.20*-x86_64"
      root-device-type    = "ebs"
      virtualization-type = "hvm"
    }
    most_recent = true
    owners      = ["137112412989"]
  }
  ssh_username = "ec2-user"
  tags = {
    Base_AMI_Name = "{{ .SourceAMIName }}"
    Name          = "${var.image_prefix}-${local.timestamp}"
    Source_AMI    = "{{ .SourceAMI }}"
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
      Source_AMI    = "{{ .SourceAMI }}"
      Creation_Date = "{{ .SourceAMICreationDate }}"
    }
  }
}
