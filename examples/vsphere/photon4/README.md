# Photon OS 4

This example builds Photon OS 4 with vSphere.

Edit the `variables.pkrvars.hcl` file to configure the credentials for the default account on machine images.

```hcl title="variables.pkrvars.hcl"
build_username           = "example"
build_password           = "<plaintext_password>"
build_password_encrypted = "<sha512_encrypted_password>"
build_key                = "<public_key>"
```

Run the following command to generate a SHA-512 encrypted password for the `build_password_encrypted` using mkpasswd.

```shell
docker run -it --rm alpine:latest
mkpasswd -m sha512
```

Then run packer build:

```
packer build -force -var-file variables.pkrvars.hcl .
```

Kudos: This example is based on [packer-examples-for-vsphere](https://github.com/vmware-samples/packer-examples-for-vsphere/tree/main/builds/linux/photon/5)