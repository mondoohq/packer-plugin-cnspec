#!/bin/bash
# Mondoo, Inc
# This script installs the latest version of the packer plugin for cases where the new HCL format cannot be used
# Please also have a look at the packer documentation https://www.packer.io/docs/plugins

set +x
set -e
set -o pipefail

os=linux
arch=amd64

# automatic download of latest version
version=$(curl https://api.github.com/repos/mondoohq/packer-plugin-mondoo/releases/latest | jq -r .name)
# alternative set the version manually
# version=v0.4.0

archive="packer-plugin-mondoo_${version}_x5.0_${os}_${arch}.zip"
sha="packer-plugin-mondoo_${version}_SHA256SUMS"
echo Download "${archive}" from
url="https://github.com/mondoohq/packer-plugin-mondoo/releases/download/${version}"
echo "${url}"

curl -sSL "${url}/${archive}" > ${archive}
curl -sSL "${url}/${sha}" > ${sha}

echo "Validating checksum..."
cat ${sha} | grep ${archive} | sha256sum -c

unzip ${archive}
rm ${archive} ${sha}

mkdir -p ~/.packer.d/plugins
mv "packer-plugin-mondoo_${version}_x5.0_${os}_${arch}" ~/.packer.d/plugins/packer-plugin-mondoo
echo "Marking executable..."
chmod +x ~/.packer.d/plugins/packer-plugin-mondoo

