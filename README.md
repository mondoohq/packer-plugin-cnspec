# Packer Plugin for Mondoo cnspec

[cnspec](https://github.com/mondoohq/cnspec) scans [Packer](https://www.packer.io) builds for vulnerabilities and misconfigurations by executing security policies-as-code enabled in [Mondoo Platform](https://console.mondoo.com). Mondoo Platform comes stocked with an ever-increasing collection of certified security policies which can be easily customize to meet your needs. 

cnspec supports scanning of Linux, Windows, and macOS, as well as Docker containers.

## Get Started with cnspec

If you are new to cnspec you can [get started](https://mondoo.com/docs/cnspec/) today!

## Packer Plugin for Mondoo cnspec tutorial

Check out the [Building secure AMIs with Mondoo and Packer](https://mondoo.com/docs/tutorials/aws/build-secure-amis-packer/) tutorial on the Mondoo documentation site.

# Installation

## Using the packer init command
Starting from version 1.7, Packer supports a new `packer init` command allowing automatic installation of Packer plugins. Read the [Packer documentation](https://www.packer.io/docs/commands/init) for more information.

To install this plugin, copy and paste this code into your Packer configuration . Then, run `packer init`.

```hcl
packer {
  required_plugins {
    mondoo = {
      version = ">= 0.4.0"
      source  = "github.com/mondoohq/mondoo"
    }
  }
}
```

#### Manual installation

You can find pre-built binary releases of the plugin [here](https://github.com/mondoohq/packer-plugin-mondoo/releases).

Once you have downloaded the latest archive corresponding to your target OS, uncompress it to retrieve the plugin binary file corresponding to your platform. To install the plugin, please follow the Packer documentation on
[installing a plugin](https://www.packer.io/docs/extending/plugins/#installing-plugins).

### Build from source

If you prefer to build the plugin from sources, clone the GitHub repository locally and run the command `go build` from the root directory. Upon successful compilation, a `packer-plugin-mondoo` plugin binary file can be found in the root directory. To install the compiled plugin, please follow the official Packer documentation on [installing a plugin](https://www.packer.io/docs/extending/plugins/#installing-plugins).

## Configuration

| **Name** | **Description** | **Type** | **Default** | **Required** |
|---|---|------------------|-------------|--------------|
| `annotations`     | Custom annotations can be applied to Packer build assets to provide additional metadata for asset tracking.  | `map of strings` | None | No |
| `asset_name`      | Overwrite the asset name in Mondoo Platform. | `string` | None | No |
| `on_failure`      | Set `on_failure = "continue"` to ignore build failures that do not meet any set `score_threshold`.| `string` | None | No |
| `score_threshold` | Set a score threshold for Packer builds `[0-100]`. Any scans that fall below the `score_threshold` will fail unless `on_failure = "continue"`. For more information see [Policy Scoring](https://mondoo.com/docs/platform/policies/scoring/index.html) in the Mondoo documentation. | `int`            | None        | No           |
| `sudo`            | Use sudo to elevate permissions when running Mondoo scans. | `bool`         | None        | No           |
| `mondoo_config_path`            | The path to the configuration to be used when running Mondoo scans. | `string`         | None        | No           |


### Example: Complete Configuration

```hcl
provisioner "mondoo" {
  on_failure      = "continue"
  score_threshold = 85
  asset_name      = "example-secure-base-image"
  sudo {
    active = true
  }

  annotations = {
    Source_AMI    = "{{ .SourceAMI }}"
    Creation_Date = "{{ .SourceAMICreationDate }}"
  }
}
```

## Sample Packer Templates

You can find example Packer templates in the [examples](/examples/) directory in this repository.

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