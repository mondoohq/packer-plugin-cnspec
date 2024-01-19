# Packer Plugin for Mondoo cnspec

![packer-plugin-cnspec illustration](.github/social/preview.jpg)

Packer Plugin [cnspec](https://github.com/mondoohq/cnspec) by [Mondoo](https://mondoo.com) scans Linux and Windows [HashiCorp Packer](https://www.packer.io) builds for vulnerabilities and security misconfigurations. The plugin retrieves CVE data from Mondoo that is updated daily with the latest CVEs and advisories. Additionally, cnspec runs security scans using [cnspec-policies](https://github.com/mondoohq/cnspec-policies) to uncover common misconfigurations that open your hosts to the risk of attack. cnspec supports scanning Linux, Windows, and macOS, as well as Docker containers.

## Plugin modes

Packer Plugin cnspec is designed to work in one of two modes:

- **Unregistered** - In unregistered mode, the plugin works without being registered with Mondoo Platform, and is designed to provide baseline security scanning with minimal configuration. On Linux builds, the plugin runs the [Linux Security by Mondoo](https://github.com/mondoohq/cnspec-policies/blob/main/core/mondoo-linux-security.mql.yaml) policy. On Windows builds, the plugin runs the [Windows Security by Mondoo](https://github.com/mondoohq/cnspec-policies/blob/main/core/mondoo-windows-security.mql.yaml) policy. Each of these policies provides security hardening checks based on industry standards for Linux and Windows. Scan results display in STDOUT during the Packer run.

- **Registered** - In registered mode, the plugin is registered with your account in Mondoo Platform using a service account. This allows you to configure and customize any of the policies in Mondoo Platform, including CIS benchmarks and more. Scan results are shown in STDOUT and sent back to Mondoo Platform for your records.

## Tutorials

Check out the Packer tutorials on the Mondoo documentation site:

- [Build secure AMIs with Mondoo and Packer](https://mondoo.com/docs/cnspec/cnspec-aws/cnspec-aws-packer/)

- [Build secure VM images in Google Cloud with cnspec and HashiCorp Packer](https://mondoo.com/docs/cnspec/cnspec-gcp/cnspec-gcp-packer/)

# Install Packer plugin cnspec

You can install Packer Plugin cnspec using the `packer init` command, install it manually, or build it from source.

## Install using the packer init command

As of version 1.7, Packer's `packer init` command allows automatic installation of Packer plugins. For more information, read the [Packer documentation](https://www.packer.io/docs/commands/init).

To install Packer Plugin cnspec:

1. Copy and paste this code into your Packer configuration.

```hcl
packer {
  required_plugins {
    cnspec = {
      version = ">= 9.0.0"
      source  = "github.com/mondoohq/cnspec"
    }
  }
}
```

2. Run this command: `packer init`

### Install manually

You can find pre-built binary releases of the plugin [here](https://github.com/mondoohq/packer-plugin-cnspec/releases).

Once you have downloaded the latest archive corresponding to your target OS, uncompress it to retrieve the plugin binary file corresponding to your platform. To install the plugin, follow the Packer documentation on
[installing a plugin](https://www.packer.io/docs/extending/plugins/#installing-plugins).

### Build from source

If you prefer to build the plugin from source:

1. Clone this GitHub repository locally.

2. Run this command from the root directory: `go build`

3. After you successfully compile, the `packer-plugin-cnspec` plugin binary file is in the root directory. Copy the binary into `~/.packer.d/plugins/` by running this command: `make dev`

4. To install the compiled plugin, follow the Packer documentation on [installing plugins](https://developer.hashicorp.com/packer/docs/plugins/install-plugins).

After building the cnspec plugin successfully, use the latest version of Packer to build a machine and verify your changes. In the [example folder](https://github.com/mondoohq/packer-plugin-cnspec/blob/main/examples) we provide a basic template. To force Packer to use the development binary installed in the previous step, comment out the `packer {}` block.

To use the developer plugin, set the packer plugin environment variable:

```bash
export PACKER_PLUGIN_PATH=~/.packer.d/plugins
packer build amazon-linux-2.pkr.hcl
```

## Configure Packer Plugin cnspec

| **Name**             | **Description**                                                                                                                                                                                                                                                                                            | **Type**         | **Default** | **Required** |
| -------------------- | ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | ---------------- | ----------- | ------------ |
| `annotations`        | Apply custom annotations to Packer build assets to provide additional metadata for asset tracking.                                                                                                                                                                                                         | `map of strings` | None        | No           |
| `asset_name`         | Overwrite the asset name in Mondoo Platform.                                                                                                                                                                                                                                                               | `string`         | None        | No           |
| `on_failure`         | Set `on_failure = "continue"` to ignore build failures that do not meet any set `score_threshold`.                                                                                                                                                                                                         | `string`         | None        | No           |
| `score_threshold`    | Set a score threshold for Packer builds `[0-100]`. Any scans that fall below the `score_threshold` will fail unless `on_failure = "continue"`. To learn more, read [How Mondoo scores policies](https://mondoo.com/docs/platform/console/monitor/#how-mondoo-scores-policies) in the Mondoo documentation. | `int`            | None        | No           |
| `sudo`               | Use sudo to elevate permissions when running Mondoo scans.                                                                                                                                                                                                                                                 | `bool`           | None        | No           |
| `mondoo_config_path` | The path to the Mondoo's service account. Defaults to `$HOME/.config/mondoo/mondoo.yml`                                                                                                                                                                                                                    | `string`         | None        | No           |
| `output`             | Set output format: compact, csv, full, json, junit, report, summary, yaml (default "compact")                                                                                                                                                                                                              | `string`         | None        | No           |
| `output_target`      | Set output target. E.g. path to local file `result.xml`                                                                                                                                                                                                                                                    | `string`         | None        | No           |

### Example: Complete Configuration

A simple configuration where we set a score threshold of 85 and use sudo to elevate permissions when running the scans:

```hcl
provisioner "cnspec" {
  on_failure      = "continue"
  score_threshold = 85
  sudo {
    active = true
  }
}
```

The following configuration shows how to set the output format to JUnit and the output target to `test-results.xml`:

```hcl
provisioner "cnspec" {
  on_failure = "continue"
  output = "junit"
  output_target = "test-results.xml"
}
```

## Sample Packer Templates

You can find example Packer templates in the [examples](/examples/) directory in this repository. You can also find a [GitHub Action workflow example](/examples/github-actions/packer-build-scan.yaml) of how to use cnspec to test builds as part of a CI/CD pipeline.

## Get Started with cnspec

cnspec's benefits extend well beyond securing Packer builds! To start exploring, [download cnspec](https://mondoo.com/docs/cnspec/).

## Contributing

If you think you've found a bug in the code or you have a question about using this software, please reach out to us by opening an issue in this GitHub repository.

Contributions to this project are welcome! If you want to fix a bug, please do so by opening a pull request in this GitHub repository. If you want to add a feature, please start by opening an issue in this GitHub repository to discuss it with us beforehand.

### Join the community!

Join the [Mondoo Community GitHub Discussions](https://github.com/orgs/mondoohq/discussions) to collaborate on policy as code and security automation.
