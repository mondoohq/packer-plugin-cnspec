# Rocky Linux 8 

This example builds Rocky Linux 8 with vSphere.

Update the `variables.pkrvars.hcl` file with your settings. Then run packer build:

```
packer build -force -var-file variables.pkrvars.hcl .
```

Kudos: This example is based on [packer-examples-for-vsphere](https://github.com/vmware-samples/packer-examples-for-vsphere/tree/main/builds/linux/rocky/8)