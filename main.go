package main

//go:generate go run github.com/hashicorp/packer-plugin-sdk/cmd/packer-sdc mapstructure-to-hcl2 -type Config,SudoConfig

import (
	"github.com/hashicorp/packer-plugin-sdk/plugin"
)

func main() {
	server, err := plugin.Server()
	if err != nil {
		panic(err)
	}
	server.RegisterProvisioner(new(Provisioner))
	server.Serve()
}
