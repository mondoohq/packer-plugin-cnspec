# CentOS 7 Example

## Build with vsphere plugin

Create the `centos-7-x86_64.vsphere.json` file with your credentials.

```json
{
  "vsphere_endpoint":"vsphere ip",
  "vsphere_username":"admin@vsphere.local",
  "vsphere_password":"password",
  "vsphere_host": "esxi host",
  "vsphere_datacenter": "datacenter",
  "vsphere_datastore" : "datastore-packer",
  "vsphere_network" : "VM Network"
}
```

Then run packer build:

```
packer build --force -var-file centos-7-x86_64.vsphere.json centos-7-x86_64.vsphere.pkr.hcl
```