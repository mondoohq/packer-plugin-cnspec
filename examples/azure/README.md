# Azure

This example shows how to build a Windows Server 2019 image in Azure. It uses
the [Azure RM Builder](https://www.packer.io/docs/builders/azure.html) to create a VM, install Windows, run a PowerShell
script to configure the VM, and then run cnspec packer plugin to assess the security. 

As a prerequisite, you need to have an Azure account. If you don't have one, you can create a free account. Then install [Packer](https://www.packer.io/downloads.html)
   and [Azure CLI](https://docs.microsoft.com/en-us/cli/azure/install-azure-cli)

```shell
az login
```

## Ubuntu

The Ubuntu example will build an Ubuntu 18.04 image in Azure.

```shell
cd ubuntu
```

Update the `vars.json` file with your Azure subscription ID and resource group. Then run the 
following command to build the image:

```shell
packer build -var-file=vars.json ubuntu-18.04.json
```


## Windows

The Windows example will build a Windows Server 2019 image in Azure.

```shell
cd windows
```

Update the variables.pkrvars.hcl file with your Azure subscription ID and tenant ID. Then install all the required
plugins

```shell
packer init windows.pkr.hcl
```

Now build the image with the following command:

```shell
packer build -var-file=variables.hcl windows.pkr.hcl
```
