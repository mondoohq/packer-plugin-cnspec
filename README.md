# Packer Plugin for Mondoo cnspec

![packer-plugin-cnspec illustration](.github/social/preview.jpg)

Packer plugin [cnspec](https://github.com/mondoohq/cnspec) by [Mondoo](https://mondoo.com) scans Linux and Windows [HashiCorp Packer](https://www.packer.io) builds for vulnerabilities and security misconfigurations. The plugin retrieves CVE data from Mondoo, which is updated daily with the latest CVEs and advisories. Additionally, cnspec runs security scans using [cnspec-policies](https://github.com/mondoohq/cnspec-policies) to uncover common misconfigurations that open your hosts to the risk of attack. cnspec supports scanning of Linux, Windows, and macOS, as well as Docker containers.

## Plugin modes

Packer plugin cnspec is designed to work in one of two modes:

- **Unregistered** - In unregistered mode, the plugin works without being registered to Mondoo Platform, and is designed to provide baseline security scanning with minimal configuration. The plugin runs either the [Linux Security by Mondoo](https://github.com/mondoohq/cnspec-policies/blob/main/core/mondoo-linux-security.mql.yaml) policy on Linux builds, or the [Windows Security by Mondoo](https://github.com/mondoohq/cnspec-policies/blob/main/core/mondoo-windows-security.mql.yaml) policy on Windows builds. Each of these policies provides security hardening checks based off of industry standards for Linux and Windows. Scan results are shown in STDOUT during the Packer run.  
- **Registered** - In registered mode, the plugin is registered to your account in Mondoo Platform using a service account. Registered mode allows you to configure and customize any of the policies in Mondoo Platform including CIS benchmarks and more. Scan results are shown in STDOUT and sent back to Mondoo Platform for your records.



## Tutorials

Check out the Packer tutorials on the Mondoo documentation site:

- [Building secure AMIs with Mondoo and Packer]([https://mondoo.com/docs/cnspec/cnspec-aws/cnspec-aws-packer/)
- [Building secure VM images in Google Cloud with cnspec and HashiCorp Packer](https://mondoo.com/docs/cnspec/cnspec-gcp/cnspec-gcp-packer/) 

# Installation

## Using the packer init command
Starting from version 1.7, Packer supports a new `packer init` command allowing automatic installation of Packer plugins. Read the [Packer documentation](https://www.packer.io/docs/commands/init) for more information.

To install this plugin, copy and paste this code into your Packer configuration . Then, run `packer init`.

```hcl
packer {
  required_plugins {
    cnspec = {
      version = ">= 6.1.3"
      source  = "github.com/mondoohq/cnspec"
    }
  }
}
```

#### Manual installation

You can find pre-built binary releases of the plugin [here](https://github.com/mondoohq/packer-plugin-cnspec/releases).

Once you have downloaded the latest archive corresponding to your target OS, uncompress it to retrieve the plugin binary file corresponding to your platform. To install the plugin, please follow the Packer documentation on
[installing a plugin](https://www.packer.io/docs/extending/plugins/#installing-plugins).

### Build from source

If you prefer to build the plugin from source, clone the GitHub repository locally and run the command `go build` from the root directory. Upon successful compilation, a `packer-plugin-cnspec` plugin binary file can be found in the root directory. To install the compiled plugin, please follow the official Packer documentation on [installing a plugin](https://www.packer.io/docs/extending/plugins/#installing-plugins).

By using `make dev`, the binary is copied into `~/.packer.d/plugins/` after the build.

## Configuration

| **Name** | **Description** | **Type** | **Default** | **Required** |
|---|---|------------------|-------------|--------------|
| `annotations`     | Custom annotations can be applied to Packer build assets to provide additional metadata for asset tracking.  | `map of strings` | None | No |
| `asset_name`      | Overwrite the asset name in Mondoo Platform. | `string` | None | No |
| `on_failure`      | Set `on_failure = "continue"` to ignore build failures that do not meet any set `score_threshold`.| `string` | None | No |
| `score_threshold` | Set a score threshold for Packer builds `[0-100]`. Any scans that fall below the `score_threshold` will fail unless `on_failure = "continue"`. To learn more, read [How Mondoo scores policies](https://mondoo.com/docs/platform/console/monitor/#how-mondoo-scores-policies) in the Mondoo documentation. | `int`            | None        | No           |
| `sudo`            | Use sudo to elevate permissions when running Mondoo scans. | `bool`         | None        | No           |
| `mondoo_config_path`            | The path to the configuration to be used when running Mondoo scans. | `string`         | None        | No           |


### Example: Complete Configuration

```bash
provisioner "cnspec" {
  on_failure      = "continue"
  score_threshold = 85
  sudo {
    active = true
  }
}
```

## Sample Packer Templates

You can find example Packer templates in the [examples](/examples/) directory in this repository.

## Get Started with cnspec

If you want to use cnspec outside of packer, you can [get started](https://mondoo.com/docs/cnspec/) today!

## Contributing

* If you think you've found a bug in the code or you have a question regarding
  the usage of this software, please reach out to us by opening an issue in
  this GitHub repository.
* Contributions to this project are welcome: if you want to add a feature or a
  fix a bug, please do so by opening a Pull Request in this GitHub repository.
  In case of feature contribution, we kindly ask you to open an issue to
  discuss it beforehand.

### Join the community!

Join the [Mondoo Community GitHub Discussions](https://github.com/orgs/mondoohq/discussions) to collaborate on policy as code and security automation.
