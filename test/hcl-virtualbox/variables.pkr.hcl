# Copyright (c) Mondoo, Inc.
# SPDX-License-Identifier: BUSL-1.1


variable "rootpassword" {
  description = "root password"
  default = "password"
  type = string
}

variable "sshusername" {
  description = "ssh username"
  default = "vagrant"
  type = string
}

variable "sshpassword" {
  description = "ssh password"
  default = "vagrant"
  type = string
}