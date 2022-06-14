.PHONY: prep/plugins install build test

PROVISIONER_BINARY_NAME=packer-provisioner-mondoo
PLUGINS_DIR=~/.packer.d/plugins

prep/plugins:
	mkdir -p ${PLUGINS_DIR}

build/generate:
	go generate ./...

build/snapshot:
	goreleaser --snapshot --skip-publish --rm-dist

build/dev:
	CGO_ENABLED=0 installsuffix=cgo go build -ldflags="-X 'main.Version=development'" -o ./dist/${PROVISIONER_BINARY_NAME}

install: prep/plugins build/dev
	rm ${PLUGINS_DIR}/${PROVISIONER_BINARY_NAME} || true
	cp ./dist/${PROVISIONER_BINARY_NAME} ${PLUGINS_DIR}/${PROVISIONER_BINARY_NAME}

test:
	go test -v

# use -debug & PACKER_LOG=1 to call step modus
test/packer/json:
	cd test/alpine3.11 && packer build -force alpine-3.11-x86_64.json

test/packer/hcl-virtualbox:
	cd test/hcl-virtualbox && packer build -force .

test/packer/hcl-docker:
	cd test/hcl-docker && packer build -force .

test/packer/proxy:
	cd test/centos7 && packer build centos-7-x86_64.json

test/packer/policybundle:
	cd test/policybundle && packer build centos-7.json