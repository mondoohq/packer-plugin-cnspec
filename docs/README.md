Packer plugin [cnspec](https://github.com/mondoohq/cnspec) by [Mondoo](https://mondoo.com) scans Linux and Windows [HashiCorp Packer](https://www.packer.io) builds for vulnerabilities and security misconfigurations. The plugin retrieves CVE data from Mondoo, which is updated daily with the latest CVEs and advisories. Additionally, cnspec runs security scans using [cnspec-policies](https://github.com/mondoohq/cnspec-policies) to uncover common misconfigurations that open your hosts to the risk of attack. cnspec supports scanning of Linux, Windows, and macOS, as well as Docker containers.

Packer plugin cnspec is designed to work in one of two modes:

- **Unregistered** - In unregistered mode, the plugin works without being registered to Mondoo Platform, and is designed to provide baseline security scanning with minimal configuration. The plugin runs either the [Linux Security by Mondoo](https://github.com/mondoohq/cnspec-policies/blob/main/core/mondoo-linux-security.mql.yaml) policy on Linux builds, or the [Windows Security by Mondoo](https://github.com/mondoohq/cnspec-policies/blob/main/core/mondoo-windows-security.mql.yaml) policy on Windows builds. Each of these policies provides security hardening checks based off of industry standards for Linux and Windows. Scan results are shown in STDOUT during the Packer run.
- **Registered** - In registered mode, the plugin is registered to your account in Mondoo Platform using a service account. Registered mode allows you to configure and customize any of the policies in Mondoo Platform including CIS benchmarks and more. Scan results are shown in STDOUT and sent back to Mondoo Platform for your records.

### Installation

To install this plugin, copy and paste this code into your Packer configuration, then run [`packer init`](https://www.packer.io/docs/commands/init).

```hcl
packer {
  required_plugins {
    cnspec = {
      version = ">= 10.0.0"
      source  = "github.com/mondoohq/cnspec"
    }
  }
}
```

Alternatively, you can use `packer plugins install` to manage installation of this plugin.

```sh
$ packer plugins install github.com/mondoohq/cnspec
```

### Components

#### Provisioners

- [cnspec](/packer/integrations/mondoohq/cnspec/latest/components/provisioner/cnspec) - Packer plugin [cnspec](https://github.com/mondoohq/cnspec) by [Mondoo](https://mondoo.com) scans 
Linux and Windows machine images for vulnerabilities and security misconfigurations. The plugin retrieves CVE data from Mondoo, which is updated daily with the latest CVEs and advisories. Additionally, cnspec runs policy-as-code security scans using [cnspec-policies](https://github.com/mondoohq/cnspec-policies) to uncover common misconfigurations that open your hosts to the risk of attack. 
- [mondoo](/packer/integrations/mondoohq/cnspec/latest/components/provisioner/mondoo) - Deprecated. Use the `cnspec` provisioner instead.

### Tutorials

Check out the Packer tutorials on the Mondoo documentation site:

- [Building secure AMIs with Mondoo and Packer](https://mondoo.com/docs/cnspec/cnspec-aws/cnspec-aws-packer/)
- [Building secure VM images in Google Cloud with cnspec and HashiCorp Packer](https://mondoo.com/docs/cnspec/cnspec-gcp/cnspec-gcp-packer/)

