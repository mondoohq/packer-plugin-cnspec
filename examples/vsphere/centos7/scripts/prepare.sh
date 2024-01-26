# Copyright (c) Mondoo, Inc.
# SPDX-License-Identifier: BUSL-1.1

set -ux

# speeds up testing
# yum -y update

echo "UseDNS no" >> /etc/ssh/sshd_config

# prep vagrant
date > /etc/vagrant_box_build_time
mkdir -pm 700 /home/vagrant/.ssh
curl -o /home/vagrant/.ssh/authorized_keys 'https://raw.githubusercontent.com/mitchellh/vagrant/master/keys/vagrant.pub'
chown -R vagrant:vagrant /home/vagrant/.ssh
chmod -R go-rwsx /home/vagrant/.ssh

# virtual box
yum -y install bzip2 perl kernel-devel-`uname -r` dkms gcc
yum_rev=$(yum history stats | grep -E "^Transactions:" | cut -d : -f 2)

VBOX_VERSION=$(cat /home/vagrant/.vbox_version)
cd /tmp
mount -o loop /home/vagrant/VBoxGuestAdditions_$VBOX_VERSION.iso /mnt
sh /mnt/VBoxLinuxAdditions.run
umount /mnt
rm -rf /home/vagrant/VBoxGuestAdditions_*.iso

[[ -n "$yum_rev" ]] && yum -y history undo $yum_rev

# clean
yum -y clean all
dd if=/dev/zero of=/EMPTY bs=1M
rm -f /EMPTY
sync