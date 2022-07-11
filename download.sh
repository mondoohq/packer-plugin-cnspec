#!/bin/bash
# Mondoo, Inc
# This script installs the latest version of the packer plugin for cases where the new HCL format cannot be used
# Please also have a look at the packer documentation https://www.packer.io/docs/plugins

set +x

os=linux
arch=amd64

# automatic download of latest version
version=$(curl https://api.github.com/repos/mondoohq/packer-plugin-mondoo/releases/latest | jq -r .name)
# alternative set the version manually
# version=v0.4.0
echo Download packer-plugin-mondoo "${version}" from
url="https://github.com/mondoohq/packer-plugin-mondoo/releases/download/${version}/packer-plugin-mondoo_${version}_x5.0_${os}_${arch}.zip"
echo "${url}"

curl -sSL "${url}" > packer-plugin-mondoo.zip
unzip packer-plugin-mondoo.zip
rm packer-plugin-mondoo.zip

mkdir -p ~/.packer.d/plugins
mv "packer-plugin-mondoo_${version}_x5.0_${os}_${arch}" ~/.packer.d/plugins/packer-plugin-mondoo
chmod +x ~/.packer.d/plugins/packer-plugin-mondoo

