NAME=cnspec
BINARY=packer-plugin-${NAME}

COUNT?=1
TEST?=$(shell go list ./...)
HASHICORP_PACKER_PLUGIN_SDK_VERSION?=$(shell go list -m github.com/hashicorp/packer-plugin-sdk | cut -d " " -f2)

ifndef LATEST_VERSION_TAG
# echo "read LATEST_VERSION_TAG from git"
LATEST_VERSION_TAG=$(shell git describe --abbrev=0 --tags)
endif

ifndef MANIFEST_VERSION
# echo "read MANIFEST_VERSION from git"
MANIFEST_VERSION=$(shell git describe --abbrev=0 --tags)
endif

ifndef TAG
# echo "read TAG from git"
TAG=$(shell git log --pretty=format:'%h' -n 1)
endif

ifndef VERSION
# echo "read VERSION from git"
VERSION=${LATEST_VERSION_TAG}+$(shell git rev-list --count HEAD)
endif

ifndef CNSPEC_VERSION
CNSPEC_VERSION=$(shell go list -json -m go.mondoo.com/cnspec/v9 | jq -r ".Version")
endif

.PHONY: dev

build:
	CGO_ENABLED=0 go build -o ${BINARY} -ldflags="-X go.mondoo.com/cnquery/v9.Version=${CNSPEC_VERSION} -X go.mondoo.com/packer-plugin-cnspec/version.Version=${VERSION} -X go.mondoo.com/packer-plugin-cnspec/version.Build=${TAG}"

dev: build
	@mkdir -p ~/.packer.d/plugins/
	@mv ${BINARY} ~/.packer.d/plugins/${BINARY}

.PHONY: dev/linux
dev/linux: build
	@mkdir -p ~/.packer.d/plugins/github.com/mondoohq/cnspec/
	@mv ${BINARY} ~/.packer.d/plugins/github.com/mondoohq/cnspec/${BINARY}_${VERSION}_x5.0_linux_amd64
	@cat ~/.packer.d/plugins/github.com/mondoohq/cnspec/packer-plugin-cnspec_${VERSION}_x5.0_linux_amd64 | sha256sum -z --tag | cut -d"=" -f2 | tr -d " " > ~/.packer.d/plugins/github.com/mondoohq/cnspec/packer-plugin-cnspec_${VERSION}_x5.0_linux_amd64_SHA256SUM

.PHONY: dev/macos
dev/macos: build
	@mkdir -p ~/.packer.d/plugins/github.com/mondoohq/cnspec/
	@mv ${BINARY} ~/.packer.d/plugins/github.com/mondoohq/cnspec/${BINARY}_${VERSION}_macos_amd64
	@cat ~/.packer.d/plugins/github.com/mondoohq/cnspec/packer-plugin-cnspec_${VERSION}_macos_amd64 | shasum --tag | cut -d"=" -f2 | tr -d " " > ~/.packer.d/plugins/github.com/mondoohq/cnspec/packer-plugin-cnspec_${VERSION}_macos_amd64_SHA256SUM
	
test:
	@go test -race -count $(COUNT) $(TEST) -timeout=3m

test/golanglint:
	@golangci-lint run

install-packer-sdc: ## Install packer software development command
	@go install github.com/hashicorp/packer-plugin-sdk/cmd/packer-sdc@${HASHICORP_PACKER_PLUGIN_SDK_VERSION}

ci-release-docs: install-packer-sdc
	@packer-sdc renderdocs -src docs -partials docs-partials/ -dst docs/
	@/bin/sh -c "[ -d docs ] && zip -r docs.zip docs/"

plugin-check: install-packer-sdc build
	@packer-sdc plugin-check ${BINARY}

testacc: dev
	@PACKER_ACC=1 go test -count $(COUNT) -v $(TEST) -timeout=120m

generate: install-packer-sdc
	@go generate ./...
	packer-sdc renderdocs -src ./docs -dst ./.docs -partials ./docs-partials
	# checkout the .docs folder for a preview of the docs

# Copywrite Check Tool: https://github.com/hashicorp/copywrite
license: license/headers/check

license/headers/check:
	copywrite headers --plan

license/headers/apply:
	copywrite headers
