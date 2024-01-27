# Azure

This example shows how to build a Windows Server 2019 image in Azure. It uses
the [Azure RM Builder](https://www.packer.io/docs/builders/azure.html) to create a VM, install Windows, run a PowerShell
script to configure the VM, and then run cnspec packer plugin to assess the security.

1. Install [Packer](https://www.packer.io/downloads.html)
   and [Azure CLI](https://docs.microsoft.com/en-us/cli/azure/install-azure-cli)

```shell
az login
```

Update the variables.pkrvars.hcl file with your Azure subscription ID and tenant ID.

2. Install all the required plugins

```shell
packer init windows.pkr.hcl
```

3. Build the image

```shell
packer build -var-file=variables.hcl windows.pkr.hcl
```
