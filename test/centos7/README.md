# CentOS Test

## Build with vagrant

Run `packer build centos-7-x86_64.json` from this directory. It will build the CentOS 7 image.

## Build with vsphere plugin

Create the `centos-7-x86_64.vsphere.variables` file with your credentials.

```json
{
    "vcenter_server":"your vSphere IP address",
    "username":"administrator@vsphere.local",
    "password":"vSphere password",
    "host":"ESXi IP address or dns",
}
```

Then run packer build:

```
packer build -var-file centos-7-x86_64.vsphere.variables centos-7-x86_64.vsphere.json
```