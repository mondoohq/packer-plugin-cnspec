# cnspec packer plugin with custom policy

This example demonstrates how the cnspec packer plugin leverages custom policy bundles.

1. Create a new policy like `custom-policy.mql.yaml`
2. Include the policy bundle in provider spec

```hcl
provisioner "cnspec" {
  policybundle = "custom-policy.mql.yaml" 
}
```

3. Run `packer init docker-ubuntu.pkr.hcl`
4. Run `packer build docker-ubuntu.pkr.hcl`