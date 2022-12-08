#!/bin/bash
# Mondoo, Inc
# This script installs the latest version of the packer plugin for cases where the new HCL format cannot be used
# Please also have a look at the packer documentation https://www.packer.io/docs/plugins

set +x
set -e
set -o pipefail

os=""
case "$(uname -s)" in
	Linux)  os="linux" ;;
	Darwin) os="darwin" ;;
	DragonFly) os="dragonfly" ;;
	GNU/kFreeBSD) os="freebsd" ;;
	FreeBSD) os="freebsd" ;;
	OpenBSD) os="openbsd" ;;
	SunOS) os="solaris" ;;
	NetBSD) os="netbsd" ;;
	*)      fail "Cannot detect OS" ;;
esac

arch=""
case "$(uname -m)" in
	x86_64)  arch="amd64" ;;
	i386)    arch="386" ;;
	i686)    arch="386" ;;
	arm)     arch="arm" ;;
	aarch64) arch="arm64";;
	arm64)   arch="arm64";;
	*)       fail "Cannot detect architecture" ;;
esac

# automatic download of latest version
version=$(curl https://api.github.com/repos/mondoohq/packer-plugin-cnspec/releases/latest | jq -r .name)
# alternative set the version manually
# version=v0.4.0

archive="packer-plugin-cnspec_${version}_x5.0_${os}_${arch}.zip"
sha="packer-plugin-cnspec_${version}_SHA256SUMS"
echo Download "${archive}" from
url="https://github.com/mondoohq/packer-plugin-cnspec/releases/download/${version}"
echo "${url}"

curl -sSL "${url}/${archive}" > "${archive}"
curl -sSL "${url}/${sha}" > "${sha}"

# determine sha tool based on os
if [ $os = "darwin" ]; then
  sha256bin='shasum -a 256 -c'
else
  sha256bin='sha256sum -c'
fi

echo "Validating checksum ${sha} ..."
cat ${sha} | grep ${archive} | ${sha256bin}

unzip "${archive}"
rm "${archive}" "${sha}"

mkdir -p ~/.packer.d/plugins
mv "packer-plugin-cnspec_${version}_x5.0_${os}_${arch}" ~/.packer.d/plugins/packer-plugin-cnspec
echo "Marking executable..."
chmod +x ~/.packer.d/plugins/packer-plugin-cnspec

