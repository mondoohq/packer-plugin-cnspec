# Use Packer and cnspec to protect your container image

## Pre-requisites

- [Packer](https://www.packer.io/)
- [Docker](https://www.docker.com/)
- Service account with Mondoo Platform

## Run Packer

To run Packer, you need to have a `.pkr.hcl` file. In this example, we have a `docker-ubuntu.pkr.hcl` file.

With the init command, Packer will download the cnspec plugin and install it in the `.packer.d/plugins` directory.

```shell
packer init docker-ubuntu.pkr.hcl
```

Configure Mondoo Platform service account credentials either via `cnspec login` or as environment variables.

```shell
export MONDOO_CONFIG_PATH=~/.config/mondoo/space-service-account.yml
```

Now, you can run the build command to create the Docker image.

```shell
packer build docker-ubuntu.pkr.hcl
```

```shell
packer build docker-ubuntu.pkr.hcl
mondoo-docker-ubuntu-2004-secure-base.docker.ubuntu: output will be in this color.

==> mondoo-docker-ubuntu-2004-secure-base.docker.ubuntu: Creating a temporary directory for sharing data...
==> mondoo-docker-ubuntu-2004-secure-base.docker.ubuntu: Pulling Docker image: ubuntu:jammy
    mondoo-docker-ubuntu-2004-secure-base.docker.ubuntu: jammy: Pulling from library/ubuntu
    mondoo-docker-ubuntu-2004-secure-base.docker.ubuntu: Digest: sha256:e9569c25505f33ff72e88b2990887c9dcf230f23259da296eb814fc2b41af999
    mondoo-docker-ubuntu-2004-secure-base.docker.ubuntu: Status: Image is up to date for ubuntu:jammy
    mondoo-docker-ubuntu-2004-secure-base.docker.ubuntu: docker.io/library/ubuntu:jammy
    mondoo-docker-ubuntu-2004-secure-base.docker.ubuntu: What's Next?
    mondoo-docker-ubuntu-2004-secure-base.docker.ubuntu: View a summary of image vulnerabilities and recommendations → docker scout quickview ubuntu:jammy
==> mondoo-docker-ubuntu-2004-secure-base.docker.ubuntu: Starting docker container...
    mondoo-docker-ubuntu-2004-secure-base.docker.ubuntu: Run command: docker run -v /Users/chris/.packer.d/tmp892498143:/packer-files -d -i -t --entrypoint=/bin/sh -- ubuntu:jammy
    mondoo-docker-ubuntu-2004-secure-base.docker.ubuntu: Container ID: b59be8cb5e2acf5e1abb02f4cec4896840e53e6489101d2506fbf216284b27cf
==> mondoo-docker-ubuntu-2004-secure-base.docker.ubuntu: Using docker communicator to connect: 172.17.0.2
==> mondoo-docker-ubuntu-2004-secure-base.docker.ubuntu: Provisioning with shell script: /var/folders/rw/y7r077vs25l2d43bqjbhq_r80000gn/T/packer-shell2801450226
==> mondoo-docker-ubuntu-2004-secure-base.docker.ubuntu: Running cnspec packer provisioner by Mondoo (Version: 10.0.3, Build: b2dd4a6)
    mondoo-docker-ubuntu-2004-secure-base.docker.ubuntu: detected packer container image build
    mondoo-docker-ubuntu-2004-secure-base.docker.ubuntu: load config from detected /Users/chris/.config/mondoo/optimistic-chebyshev-158847.yml
    mondoo-docker-ubuntu-2004-secure-base.docker.ubuntu: using service account credentials
    mondoo-docker-ubuntu-2004-secure-base.docker.ubuntu: scan packer build
    mondoo-docker-ubuntu-2004-secure-base.docker.ubuntu: scan completed successfully
==> mondoo-docker-ubuntu-2004-secure-base.docker.ubuntu: Committing the container
    mondoo-docker-ubuntu-2004-secure-base.docker.ubuntu: Image ID: sha256:2df9972417cd375642401b53f52c070f5f257498b41e709cbaa0f1aa9a49e2ef
==> mondoo-docker-ubuntu-2004-secure-base.docker.ubuntu: Killing the container: b59be8cb5e2acf5e1abb02f4cec4896840e53e6489101d2506fbf216284b27cf
Build 'mondoo-docker-ubuntu-2004-secure-base.docker.ubuntu' finished after 20 seconds 795 milliseconds.

==> Wait completed after 20 seconds 795 milliseconds

==> Builds finished. The artifacts of successful builds are:
--> mondoo-docker-ubuntu-2004-secure-base.docker.ubuntu: Imported Docker image: sha256:2df9972417cd375642401b53f52c070f5f257498b41e709cbaa0f1aa9a49e2ef
```

In this example we configured the junit export to `test-results.xml` and the `cnspec` plugin will export the scan results in junit format.

```hcl
provisioner "cnspec" {
  ...
  output        = "junit"
  output_target = "test-results.xml"
}
```

The exported junit file will contain the scan results and can be used in CI/CD pipelines.

```xml
<?xml version="1.0" encoding="UTF-8"?>
<testsuites>
	<testsuite name="Policy Report for mondoo-ubuntu-2004-secure-base-20240207121031" tests="89" failures="58" errors="0" id="0" time="">
		<testcase name="Ensure secure permissions on /etc/passwd are set" classname="score"></testcase>
		<testcase name="Ensure root group is empty" classname="score"></testcase>
		<testcase name="Ensure auditing for processes that start prior to auditd is enabled" classname="score">
			<failure message="results do not match" type="fail"></failure>
		</testcase>
...
```