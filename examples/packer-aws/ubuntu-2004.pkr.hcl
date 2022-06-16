packer {
  required_plugins {
    amazon = {
      version = ">= 1.1.0"
      source  = "github.com/hashicorp/amazon"
    }
    mondoo = {
      version = ">= 0.2.1"
      source  = "github.com/mondoohq/mondoo"
    }
  }
}

variable "aws_profile" {
  type = string
  description = "AWS profile to use. Typically found in ~/.aws/credentials"
  default = "default"
}

variable "aws_region" {
  default = "us-east-1"
  type    = string
}

variable "image_prefix" {
  type        = string
  description = "Prefix to be applied to image name"
  default     = "mondoo-ubuntu-20.04-secure-base"
}

locals { timestamp = regex_replace(timestamp(), "[- TZ:]", "") }

source "amazon-ebs" "ubuntu2004" {
  profile       = var.aws_profile
  ami_name      = "${var.image_prefix}-${local.timestamp}"
  instance_type = "t2.micro"
  region        = var.aws_region
  source_ami_filter {
    filters = {
      name                = "ubuntu/images/hvm-ssd/ubuntu-focal-20.04-amd64-server-*"
      root-device-type    = "ebs"
      virtualization-type = "hvm"
    }
    most_recent = true
    owners      = ["099720109477"]
  }
  ssh_username = "ubuntu"
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
    "source.amazon-ebs.ubuntu2004"
  ]

  provisioner "shell" {
    inline = [
      "sudo hostnamectl set-hostname ${var.image_prefix}-${local.timestamp}",
      "sudo apt-get update -y",
      "sudo apt-get upgrade -y"
    ]
  }

  provisioner "mondoo" {
    on_failure = "continue"
    asset_name = "${var.image_prefix}-${local.timestamp}"

    annotations = {
      Name          = "${var.image_prefix}-${local.timestamp}"
      Base_AMI_Name = "{{ .SourceAMIName }}"
      Source_AMI    = "{{ .SourceAMI }}"
      Creation_Date = "{{ .SourceAMICreationDate }}"
    }
  }
}