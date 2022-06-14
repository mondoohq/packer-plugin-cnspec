# Packer Plugin Mondoo

![Mondoo Cloud-Native Security](assets/github.splash.png)

Mondoo makes it easy for you to build secure base images for your hosts using [Packer](https://www.packer.io). 

The **Packer Plugin Mondoo** tests Packer builds for vulnerabilities and misconfigurations by executing any security policies you have enabled for your environment in [Mondoo Platform](https://console.mondoo.com). Mondoo security policies cover a wide range of use cases including CIS Benchmarks, vulnerability scanning, secrets scanning, and more. If you are new to Mondoo you can get started by [signing up for a free account](https://mondoo.com/docs/tutorials/mondoo/account-setup/) today!

Mondoo supports Linux, Windows, and macOS, as well as Docker container builds. 

<!-- @import "[TOC]" {cmd="toc" depthFrom=1 depthTo=6 orderedList=false} -->

<!-- code_chunk_output -->

- [Packer Plugin Mondoo](#packer-plugin-mondoo)
  - [Usage](#usage)
  - [Configuration Reference](#configuration-reference)
    - [Required Parameters:](#required-parameters)
  - [Optional Parameters:](#optional-parameters)
  - [Mondoo Installation and Setup](#mondoo-installation-and-setup)
    - [Register for a Mondoo Platform account, and install & configure Mondoo Client](#register-for-a-mondoo-platform-account-and-install--configure-mondoo-client)
    - [Mondoo Client configuration - Local Setup](#mondoo-client-configuration---local-setup)
    - [Mondoo Client configuration - Service Account (CI/CD)](#mondoo-client-configuration---service-account-cicd)
    - [Configure Security Policies](#configure-security-policies)
  - [Install Packer](#install-packer)
  - [Install Packer Plugin Mondoo](#install-packer-plugin-mondoo)
    - [Linux Installation](#linux-installation)
    - [macOS Installation](#macos-installation)
    - [Windows Installation](#windows-installation)
    - [Compile from source](#compile-from-source)
  - [Examples](#examples)
    - [VirtualBox](#virtualbox)
    - [AWS AMI Image Build](#aws-ami-image-build)
      - [AWS AMI Packer Template Example](#aws-ami-packer-template-example)
      - [AWS AMI - Configuration](#aws-ami---configuration)
    - [Digital Ocean](#digital-ocean)
  - [Debugging:](#debugging)
  - [Uninstall](#uninstall)
  - [Author](#author)

<!-- /code_chunk_output -->

## Usage

The simplest setup is to add `mondoo` to your provisioners list (in packer's hcl syntax):

```hcl
build {
  sources = [
    "source.virtualbox-iso.alpine",
  ]

  provisioner "shell" {
    scripts = [
        "scripts/prepare.sh"
    ]
  }

  provisioner "mondoo" {
    on_failure =  "continue"
  }

  post-processor "vagrant" {
    vagrantfile_template = "Vagrantfile"
    output        = "output.box"
  }
}
```

Or in packer's json syntax:

```json
  "provisioners": [{
    "type": "shell",
    "scripts": [
      "scripts/install.sh",
    ],
    "override": {
      "virtualbox-iso": {
        "execute_command": "/bin/sh '{{.Path}}'"
      }
    }
  }, {
    "type": "mondoo",
    "on_failure": "continue",
  }],
```

The example above uses the `on_failure` option, which allows you to pass a packer build, even if vulnerabilities have been found.

If you're using WinRM as a connection type, you need to set the user and password manually (see [hashicorp/packer#7079](https://github.com/hashicorp/packer/issues/7079)):

```
    {
      "type": "mondoo",
      "winrm_user": "vagrant"
      "winrm_password": "vagrant",
    }
```

## Configuration Reference

### Required Parameters:

  * none

## Optional Parameters:

  * `on_failure` (string) - If on_failure is set to `continue` the build continues even if vulnerabilities have been found

  ```
  "on_failure": "continue",
  ```

  * `annotations` (map of string) - Custom annotations can be passed to mondoo. This eases searching for the correct asset report later.

  ```
  "annotations": {
    "mondoo.app/ami-name":  "{{user `ami_name`}}",
    "name":"Packer Builder",
    "custom_key":"custom_value"
  }
  ```

  * `asset_name` - Overwrite the asset name

  ```
  {
    "type": "mondoo",
    "asset_name": "my custom asset name"
  }
  ```

  * `sudo` - Use sudo to elevate permissions

  ```
  {
    "type": "mondoo",
    "sudo": {
      "active": true
    }
  }
  ```

## Mondoo Installation and Setup

To use the Packer Plugin Mondoo, you will need: to have Mondoo Client installed and configured, enable any security policies in your account, and install the Packer Plugin Mondoo plugin.

### Register for a Mondoo Platform account, and install & configure Mondoo Client

If you do not already have an account on Mondoo Platform, [sign up for a free account](https://mondoo.com/docs/tutorials/mondoo/account-setup/). 

### Mondoo Client configuration - Local Setup

The Packer Plugin Mondoo uses the local Mondoo Client configuration by default. Run `mondoo status` to validate your configuration:

```bash
➜  ~ mondoo status
→ loaded configuration from /etc/opt/mondoo/mondoo.yml using source default
→ Hostname:	some-hostname
→ IP:		192.168.22.22
→ Platform:	macos
→ Release:	12.3
→ Time:		2022-04-01T08:09:04-07:00
→ Version:	5.32.0 (API Version: 5)
→ API ConnectionConfig:	https://api.mondoo.app
→ API Status:	SERVING
→ API Time:	2022-04-01T15:09:05Z
→ API Version:	5
→ Space:	//captain.api.mondoo.app/spaces/my-mondoo-space
→ Client:	no managed client
→ Service Account:	//agents.api.mondoo.app/spaces/my-mondoo-space/serviceaccounts/abcdefg123456789
→ client is registered
→ client authenticated successfully
```

### Mondoo Client configuration - Service Account (CI/CD)

![Create a service account in Mondoo Platform](/assets/service_account.gif)

If you plan to run the Packer Plugin Mondoo as part of a pipeline, a service account can also be created to authenticate with Mondoo Platform to retrieve the policies you have enabled, and send results from your builds back to your account. 

To create a service account:

1. Log in to [Mondoo Platform](https://console.mondoo.com)
2. Navigate to the **SETTINGS**.
3. Click **Service Accounts**.
4. Click the **CREATE SERVICE ACCOUNT** button.
5. Click the **Base64-encoded** checkbox, then click **GENERATE NEW CREDENTIALS**.
6. Copy the Base64 encoded credentials to the clipboard.
7. Paste the credentials into a file `echo <base64_credentials> > mondoo.credentials`.

Use `MONDOO_CONFIG_PATH` environment variable to set the location of your credentials file:

```bash
export MONDOO_CONFIG_PATH=/path/to/mondoo.credentials
```

### Configure Security Policies

![Enable security policies in Mondoo Platform](/assets/enable_policies.gif).

Mondoo Platform has an ever increasing library of certified security policies and benchmarks that are production ready to test your builds. The policies are simple to enable and customize for your needs. 

To enable security policies:

1. Log in to [Mondoo Platform](https://console.mondoo.com)
2. Navigate to the **POLICY HUB**.
3. Click the **ADD POLICY** to view the all available policies.
4. Locate the policy, or policies you want to enable either by scrolling through the list of available policies, or by using the **Filter** search box.
5. Click the the checkbox next to any policy, and then click the **ENABLE** button.

Any changes will take effect immediately. Assets registered to the Space will automatically run applicable policies on their next scan.

If you need to customize the controls in a policy see [Customizing policies](https://mondoo.com/docs/tutorials/mondoo/policy-management/#customizing-policies).

## Install Packer

[Install packer](https://learn.hashicorp.com/tutorials/packer/get-started-install-cli) and verify the installation:

> Note: The plugin is tested with Packer 1.8.1

```bash
$ packer
usage: packer [--version] [--help] <command> [<args>]

Available commands are:
    build       build image(s) from template
    fix         fixes templates from old versions of packer
    inspect     see components of a template
    validate    check that a template is valid
    version     Prints the Packer version
```

## Install Packer Plugin Mondoo

To install the precompiled binary, download the appropriate package from [GitHub](https://github.com/mondoohq/packer-provisioner-mondoo/releases/latest) and place the binary in the Packer's plugin directory `~/.packer.d/plugins` (Linux, Mac) or `%USERPROFILE%/packer.d/plugins` (Windows). Other locations that Packer searches for are [documented on their website](https://www.packer.io/docs/extending/plugins.html#installing-plugins).

The following simplifies the installation:

### Linux Installation

```
mkdir -p ~/.packer.d/plugins
cd ~/.packer.d/plugins
curl -sSL https://github.com/mondoohq/packer-provisioner-mondoo/releases/latest/download/packer-provisioner-mondoo_linux_amd64.tar.gz | tar -xz > packer-provisioner-mondoo
chmod +x packer-provisioner-mondoo
```

### macOS Installation

```
mkdir -p ~/.packer.d/plugins
cd ~/.packer.d/plugins
curl -sSL https://github.com/mondoohq/packer-provisioner-mondoo/releases/latest/download/packer-provisioner-mondoo_darwin_amd64.tar.gz | tar -xz > packer-provisioner-mondoo
chmod +x packer-provisioner-mondoo
```

### Windows Installation

Download the binary from the [GitHub releases page](https://github.com/mondoohq/packer-provisioner-mondoo/releases) and put it in the same directory as your packer executable.

```powershell
# This script requires PowerShell
Invoke-WebRequest 'https://github.com/mondoohq/packer-provisioner-mondoo/releases/latest/download/packer-provisioner-mondoo_windows_amd64.zip' -O 'packer-provisioner-mondoo_windows_amd64.zip'

# extract zip and place it in the same path as packer
Expand-Archive -LiteralPath packer-provisioner-mondoo_windows_amd64.zip
Copy-Item ./packer-provisioner-mondoo_windows_amd64/packer-provisioner-mondoo.exe ((Get-Command packer).Source | Split-Path)

# clean up
Remove-Item -Recurse -Force .\packer-provisioner-mondoo_windows_amd64
Remove-Item packer-provisioner-mondoo_windows_amd64.zip
```

### Compile from source

If you wish to compile from source, you need to have [Go](https://golang.org/) installed and configured.

1. Clone the mondoo repository from GitHub into your $GOPATH:

```
$ mkdir -p $(go env GOPATH)/src/go.mondoo.io/packer-provisioner-mondoo && cd $_
$ git clone https://github.com/mondoohq/packer-provisioner-mondoo.git
```

2. Build the plugin for your current system and place the binary in the packer plugin directory

```
make install
```

## Examples

### VirtualBox

> NOTE The full example is located in [test/centos7](test/cents7/centos-7-x86_64.json)

***Packer Template***

The following configuration builds a CentOS 7 image with VirtualBox:

```
{
  "description": "Build CentOS 7 x86_64",
  "push": {
    "name": "centos7",
    "vcs": true
  },
  "variables": {},
  "provisioners": [
    {
      "type": "shell",
      "scripts": [
        "scripts/prepare.sh"
      ],
      "override": {
        "virtualbox-iso": {
          "execute_command": "echo 'vagrant'|sudo -S bash '{{.Path}}'"
        }
      }
    }, {
      "type": "mondoo",
      "debug": true
    }
  ],
  "builders": [
    {
      "type": "virtualbox-iso",
      "guest_additions_path": "VBoxGuestAdditions_{{.Version}}.iso",

      "guest_os_type": "Linux26_64",
      "headless": false,
      "disk_size": 10240,
      "http_directory": "http",

      "iso_url": "http://ftp.usf.edu/pub/centos/7.8.2003/isos/x86_64/CentOS-7-x86_64-Minimal-2003.iso",
      "iso_checksum": "659691c28a0e672558b003d223f83938f254b39875ee7559d1a4a14c79173193",
      "iso_checksum_type": "sha256",

      "ssh_username": "vagrant",
      "ssh_password": "vagrant",
      "ssh_port": 22,
      "ssh_wait_timeout": "10m",
      "shutdown_command": "echo 'vagrant'| sudo -S /sbin/poweroff",

      "boot_wait": "10s",
      "boot_command": [
        "<tab> text ks=http://{{ .HTTPIP }}:{{ .HTTPPort }}/kickstart.cfg<enter>",
        "<wait10><wait10><wait10>","<wait10><wait10><wait10>"
      ],

     "hard_drive_interface": "sata",
     "vboxmanage": [
        [ "modifyvm", "{{.Name}}", "--memory", "2048" ],
        [ "modifyvm", "{{.Name}}", "--cpus", "2" ]
      ]
    }
  ],
  "post-processors": [
    [{
      "type": "vagrant",
      "vagrantfile_template": "Vagrantfile",
      "output": "output.box"
    }]
  ]
}
```

### AWS AMI Image Build

This example illustrates the combination of Packer & Mondoo to build an AMI image. The full example is [available](https://github.com/mondoohq/packer-provisioner-mondoo/tree/master/examples/packer-aws).

#### AWS AMI Packer Template Example

The following packer templates is a simple example that builds on top of the official Ubuntu image, runs a shell provisioner and a mondoo vulnerability scan.

```json
{
  "variables": {
    "profile": "{{env `AWS_PROFILE`}}",
    "aws_region": "{{env `AWS_REGION`}}",
    "prefix": "{{env `PACKER_BUILD_PREFIX`}}",
    "timestamp": "{{isotime `2006-01-02`}}"
  },
  "builders": [{
    "type": "amazon-ebs",
    "region": "{{user `aws_region`}}",
    "source_ami_filter": {
      "filters": {
        "virtualization-type": "hvm",
        "name": "ubuntu/images/hvm-ssd/ubuntu-focal-20.04-amd64-server-*",
        "root-device-type": "ebs"
      },
      "owners": ["099720109477"],
      "most_recent": true
    },
    "instance_type": "t2.micro",
    "ssh_username": "ubuntu",
    "ami_name": "{{user `prefix`}}-ubuntu2004-base-{{user `timestamp`}}",
    "tags": {
      "Name": "Packer Builder - Ubuntu 20.04",
      "Base_AMI_Name": "{{ .SourceAMIName }}",
      "Source_AMI": "{{ .SourceAMI }}",
      "Source_AMI_Creation_Date": "{{ .SourceAMICreationDate }}"
    }
  }],
  "provisioners": [
    {
      "type": "shell",
      "inline":[
        "ls -l /home/ubuntu"
      ]
    },
    {
      "type": "mondoo",
      "on_failure": "continue",
      "labels": {
        "mondoo.app/ami-name":  "{{user `ami_name`}}",
        "name":"Packer Builder - Ubuntu 20.04",
        "asset_name": "Packer Build - Ubuntu 20.04",
        "created_on":"{{user `timestamp`}}"
      }
    }
  ]
}
```

The simplest configuration for mondoo would be:

```
{
  "type": "mondoo"
}
```

The additional `on_failure` allows Packer to continue, even if mondoo found vulnerabilities. Additional labels help you to identify the ami report on mondoo later. To verify the packer template, run packer `packer validate`:

```
$ packer validate example.json
Template validated successfully.
```

#### AWS AMI - Configuration

Once the packer template is verified, we are ready to build the image. In this case, we are going to build an AMI, therefore we need the AWS credentials to spin up a new instance. As shown above, the same will work with other cloud providers or Vagrant.

Set a prefix for the AMI name:

```bash
export PACKER_BUILD_PREFIX=mondoo
```

Set the AWS credentials using the `AWS_PROFILE` environment variable. 

```bash
export AWS_PROFILE=/path/to/aws/credentials
```

Set the `AWS_REGION` environment variable:

```bash
export AWS_REGION=us-east-1
```

Start the packer build:

```bash
$ packer build example.json

amazon-ebs output will be in this color.
==> amazon-ebs: Prevalidating AMI Name: mondoo-example 1562326441
    amazon-ebs: Found Image ID: ami-0cfee17793b08a293
==> amazon-ebs: Creating temporary keypair: packer_5d1f35a9-bf28-ad76-be7b-a7d1ba0b1a28
==> amazon-ebs: Creating temporary security group for this instance: packer_5d1f35ad-5e30-7a62-7142-05d3371896a9
==> amazon-ebs: Authorizing access to port 22 from [0.0.0.0/0] in the temporary security groups...
==> amazon-ebs: Launching a source AWS instance...
==> amazon-ebs: Adding tags to source instance
    amazon-ebs: Adding tag: "Name": "Packer Builder"
    amazon-ebs: Instance ID: i-077464c074ab682fe
==> amazon-ebs: Waiting for instance (i-077464c074ab682fe) to become ready...
==> amazon-ebs: Using ssh communicator to connect: 54.234.154.92
==> amazon-ebs: Waiting for SSH to become available...
==> amazon-ebs: Connected to SSH!
==> amazon-ebs: Provisioning with shell script: /var/folders/wb/12345/T/packer-shell496967260
    amazon-ebs: total 28
    amazon-ebs: drwxr-xr-x 4 ubuntu ubuntu 4096 Jul  5 11:34 .
    amazon-ebs: drwxr-xr-x 3 root   root   4096 Jul  5 11:34 ..
    amazon-ebs: -rw-r--r-- 1 ubuntu ubuntu  220 Aug 31  2015 .bash_logout
    amazon-ebs: -rw-r--r-- 1 ubuntu ubuntu 3771 Aug 31  2015 .bashrc
    amazon-ebs: drwx------ 2 ubuntu ubuntu 4096 Jul  5 11:34 .cache
    amazon-ebs: -rw-r--r-- 1 ubuntu ubuntu  655 May  9 20:20 .profile
    amazon-ebs: drwx------ 2 ubuntu ubuntu 4096 Jul  5 11:34 .ssh
==> amazon-ebs: Running mondoo vulnerability scan...
==> amazon-ebs: Executing Mondoo: [mondoo scan]
    amazon-ebs: Start vulnerability scan:
  →  detected automated runtime environment: Unknown CI
    amazon-ebs: 1:34PM INF ssh uses scp (beta) instead of sftp for file transfer transport=ssh
  →  verify platform access to ssh://chartmann@127.0.0.1:55661
  →  gather platform details................................
  →  detected ubuntu 16.04
  →  gather platform packages for vulnerability scan
  →  found 453 packages
    amazon-ebs:   →  analyse packages for vulnerabilities
    amazon-ebs: Advisory Report:
    amazon-ebs:   ■        PACKAGE                     INSTALLED               VULNERABLE (<)  ADVISORY
    amazon-ebs:   ■   9.8  linux-image-4.4.0-1087-aws  4.4.0-1087.98                           https://mondoo.app/advisories/
    amazon-ebs:   ├─  9.8  linux-image-4.4.0-1087-aws  4.4.0-1087.98                           https://mondoo.app/advisories/
    amazon-ebs:   ├─  8.8  linux-image-4.4.0-1087-aws  4.4.0-1087.98                           https://mondoo.app/advisories/
    amazon-ebs:   ├─  8.1  linux-image-4.4.0-1087-aws  4.4.0-1087.98                           https://mondoo.app/advisories/
    amazon-ebs:   ├─  7.8  linux-image-4.4.0-1087-aws  4.4.0-1087.98                           https://mondoo.app/advisories/
    amazon-ebs:   ├─  7.8  linux-image-4.4.0-1087-aws  4.4.0-1087.98                           https://mondoo.app/advisories/
    amazon-ebs:   ├─  7.8  linux-image-4.4.0-1087-aws  4.4.0-1087.98                           https://mondoo.app/advisories/
    amazon-ebs:   ├─  7.8  linux-image-4.4.0-1087-aws  4.4.0-1087.98                           https://mondoo.app/advisories/
    amazon-ebs:   ├─  7.8  linux-image-4.4.0-1087-aws  4.4.0-1087.98                           https://mondoo.app/advisories/
    amazon-ebs:   ├─  7.8  linux-image-4.4.0-1087-aws  4.4.0-1087.98                           https://mondoo.app/advisories/
...
  →  ■ found 70 advisories: 2 critical, 14 high, 26 medium, 3 low, 25 none, 0 unknown
  →  report is available at https://mondoo.app/v/serene-dhawan-599345/focused-darwin-833545/reports/1NakGz6ysD1MzEGT8hRJ6wow6ZQ
==> amazon-ebs: Stopping the source instance...
    amazon-ebs: Stopping instance
==> amazon-ebs: Waiting for the instance to stop...
==> amazon-ebs: Creating AMI mondoo-example 1562326441 from instance i-077464c074ab682fe
    amazon-ebs: AMI: ami-0cb9729eaa3f53209
==> amazon-ebs: Waiting for AMI to become ready...
```

As we see as the result, the mondoo scan found vulnerabilities but passed the build.

### Digital Ocean

> NOTE the full example is available as an [DigitalOcean example](https://github.com/mondoohq/packer-provisioner-mondoo/tree/master/examples/packer-digitalocean).

The following example is fully functional and builds and scans an image on DigitalOcean.

```
{
  "provisioners": [
    {
      "type": "mondoo"
    }
  ],

  "builders": [
    {
      "type": "digitalocean",
      "api_token": "DIGITALOCEAN_TOKEN",
      "image": "ubuntu-18-04-x64",
      "ssh_username": "root",
      "region": "nyc1",
      "size": "s-4vcpu-8gb"
    }
  ]
}
```

Replace the `api_token` with your own and run `packer`

```
packer build do-ubuntu.json
```

## Debugging:

To debug the mondoo scan, set the `debug` variable:

```
{
  "type": "mondoo",
  "debug": true
}
```

## Uninstall

You can easily uninstall the provisioner by removing the binary.

```
# linux & mac
rm ~/.packer.d/plugins/packer-provisioner-mondoo
```

## Author

Mondoo, Inc
