# For full specification on the configuration of this file visit:
# https://github.com/hashicorp/integration-template#metadata-configuration
integration {
  name = "Mondoo"
  description = "TODO"
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
