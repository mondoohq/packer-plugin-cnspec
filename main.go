package main

//go:generate go run github.com/hashicorp/packer-plugin-sdk/cmd/packer-sdc mapstructure-to-hcl2 -type Config,SudoConfig

import (
	"fmt"
	"os"

	"github.com/hashicorp/packer-plugin-sdk/plugin"
	"go.mondoo.com/packer-plugin-mondoo/version"
)

func main() {
	pps := plugin.NewSet()
	pps.RegisterProvisioner(plugin.DEFAULT_NAME, new(Provisioner))
	pps.SetVersion(version.PluginVersion)
	err := pps.Run()

	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}
