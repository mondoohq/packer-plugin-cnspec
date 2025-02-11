Type: `cnspec`

Packer plugin [cnspec](https://github.com/mondoohq/cnspec) by [Mondoo](https://mondoo.com) scans Linux and Windows machine images for vulnerabilities and security misconfigurations. The plugin retrieves CVE data from Mondoo, which is updated daily with the latest CVEs and advisories. Additionally, cnspec runs security scans using [cnspec-policies](https://github.com/mondoohq/cnspec-policies) to uncover common misconfigurations that open your hosts to the risk of attack. 

## Basic Example

```hcl
provisioner "cnspec" {
  on_failure          = "continue"
  score_threshold     = 85
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

## Configuration Reference

Optional Parameters:
<!-- Code generated from the comments of the Config struct in provisioner/provisioner.go; DO NOT EDIT MANUALLY -->

- `host_alias` (string) - The alias by which the host should be known.
  Defaults to `default`.

- `user` (string) - The `user` set for your communicator. Defaults to the `user` set
  by packer.

- `local_port` (uint) - The port on which to attempt to listen for SSH
   connections. This value is a starting point. The provisioner will attempt
   listen for SSH connections on the first available of ten ports, starting at
   `local_port`. A system-chosen port is used when `local_port` is missing or
   empty.

- `ssh_host_key_file` (string) - The SSH key that will be used to run the SSH
   server on the host machine to forward commands to the target machine.
   packer connects to this server and will validate the identity of the
   server using the system known_hosts. The default behavior is to generate
   and use a onetime key.

- `ssh_authorized_key_file` (string) - The SSH public key of the packer `ssh_user`.
  The default behavior is to generate and use a onetime key.

- `use_sftp` (bool) - Deprecated: SFTP is now the default. To use SCP instead, set use_scp to true

- `use_scp` (bool) - Use SCP instead of SFTP. By default, SFTP is used since
  SCP communication can fail on Windows 2025 and SSH systems.

- `debug` (bool) - Sets the log level to `DEBUG`

- `asset_name` (string) - The asset name passed to Mondoo Platform. Defaults to the hostname
  of the instance.

- `on_failure` (string) - Configure behavior whether packer should fail if `scan_threshold` is
  not met. If `scan_threshold` configuration is omitted, the threshold
  is set to `0` and builds will pass regardless of what score is
  returned.
  If `score_threshold` is set to a value, and `on_failure = "continue"`
  builds will continue regardless of what score is returned.

- `labels` (map[string]string) - Configure an optional map of `key/val` labels for the asset in
  Mondoo Platform.

- `annotations` (map[string]string) - Configure an optional map of `key/val` annotations for the asset in
  Mondoo Platform.

- `incognito` (bool) - Configures incognito mode. By default it detects if a Mondoo service account
  is available. When set to false, scan results will not be sent to
  Mondoo Platform.

- `policies` ([]string) - A list of policies to be executed (will automatically activate incognito mode).

- `policybundle` (string) - A path to local policy bundle file.

- `sudo` (\*SudoConfig) - Runs scan with `--sudo`. Defaults to none.

- `winrm_user` (string) - Configure WinRM user. Defaults to `user` set by the packer communicator.

- `winrm_password` (string) - Configure WinRM user password. Defaults to `password` set by the packer
  communicator.

- `use_proxy` (bool) - Use proxy to connect to host to scan. This configuration will fall-back to
  packer proxy for cases where the provisioner cannot access the target directly

- `output` (string) - Set output format: compact, csv, full, json, junit, report, summary, yaml
  (default "compact")

- `output_target` (string) - Set output target. E.g. path to local file

- `score_threshold` (int) - An integer value to set the `score_threshold` of mondoo scans. Defaults to
  `0` which results in a passing score regardless of what scan results are
  returned.

- `mondoo_config_path` (string) - The path to the Mondoo's service account. Defaults to
  `$HOME/.config/mondoo/mondoo.yml`

<!-- End of code generated from the comments of the Config struct in provisioner/provisioner.go; -->


### SudoConfig
<!-- Code generated from the comments of the SudoConfig struct in provisioner/provisioner.go; DO NOT EDIT MANUALLY -->

- `active` (bool) - Active

<!-- End of code generated from the comments of the SudoConfig struct in provisioner/provisioner.go; -->


## Get Started with cnspec

If you are new to cnspec, check out [Get started with cnspec](https://mondoo.com/docs/cnspec/).

## Packer plugin cnspec tutorial

Check out the Packer tutorials on the Mondoo documentation site:

- [Building secure AMIs with Mondoo and Packer](https://mondoo.com/docs/cnspec/cnspec-aws/cnspec-aws-packer/) 
- [Building secure VM images in Google Cloud with cnspec and HashiCorp Packer](https://mondoo.com/docs/cnspec/cnspec-gcp/cnspec-gcp-packer/) 

## Sample Packer Templates

You can find example Packer templates in the [examples](https://github.com/mondoohq/packer-plugin-cnspec/tree/main/examples) directory in this repository.
