# For full specification on the configuration of this file visit:
# https://github.com/hashicorp/integration-template#metadata-configuration
integration {
  name = "Mondoo"
  description = "Scans Linux and Windows HashiCorp Packer builds for vulnerabilities and security misconfigurations."
  identifier = "packer/mondoohq/cnspec"
  component {
    type = "provisioner"
    name = "Mondoo"
    slug = "mondoo"
  }
  component {
    type = "provisioner"
    name = "cnspec"
    slug = "cnspec"
  }
}
