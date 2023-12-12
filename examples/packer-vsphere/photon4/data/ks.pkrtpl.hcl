# Copyright (c) Mondoo, Inc.
# SPDX-License-Identifier: BUSL-1.1

{
    "hostname": "photon",
    "password":
        {
            "crypted": true,
            "text": "${build_password_encrypted}"
        },
    "disk": "/dev/sda",
    "partitions": [
        {"mountpoint": "/", "size": 0, "filesystem": "ext4"},
        {"mountpoint": "/boot", "size": 128, "filesystem": "ext4"},
        {"mountpoint": "/root", "size": 128, "filesystem": "ext4"},
        {"size": 128, "filesystem": "swap"}
    ],
    "bootmode": "efi",
    "packages": [
        "minimal",
        "linux-esx",
        "initramfs",
        "sudo",
        "vim",
        "cloud-utils",
        "parted"
    ],
    "postinstall": [
        "#!/bin/sh",
        "useradd -m -p '${build_password_encrypted}' -s /bin/bash ${build_username}",
        "usermod -aG sudo ${build_username}",
        "echo \"${build_username} ALL=(ALL) NOPASSWD: ALL\" >> /etc/sudoers.d/${build_username}",
        "chage -I -1 -m 0 -M 99999 -E -1 root",
        "chage -I -1 -m 0 -M 99999 -E -1 ${build_username}"
    ],
    "linux_flavor": "linux-esx",
    "network": {
        "type": "dhcp"
    }
}