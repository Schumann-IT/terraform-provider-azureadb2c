<a href="https://terraform.io">
    <img src=".github/tf.png" alt="Terraform logo" title="Terraform" align="left" height="50" />
</a>

# Terraform Provider for Azure Active Directory B2C

- [Terraform Website](https://www.terraform.io)
- [Azure AD B2C Provider Usage Examples](https://github.com/schumann-it/terraform-provider-azureadb2c/tree/main/examples)

## Usage Example

```
# Configure Terraform
terraform {
  required_providers {
    azuread = {
      source  = "schumann-it/azureadb2c"
      version = "~> 0.1.0"
    }
  }
}

# Configure the Azure Active Directory B2C Provider
provider "azureadb2c" {

  # NOTE: Environment Variables can also be used for Service Principal authentication
  # See official docs for more info: https://registry.terraform.io/providers/schumann-it/azureadb2c/latest/docs

  # client_id     = "..."
  # client_secret = "..."
  # tenant_id     = "..."
}

```

Further [usage documentation](https://registry.terraform.io/providers/schumann-it/azureadb2c/latest/docs) is available on the Terraform website.


## Developer Requirements

- [Terraform](https://www.terraform.io/downloads.html) 0.12.x or later
- [Go](https://golang.org/doc/install) 1.22.x (to build the provider plugin)

If you're building on Windows, you will also need:
- [Git Bash for Windows](https://git-scm.com/download/win)
- [Make for Windows](http://gnuwin32.sourceforge.net/packages/make.htm)

For *GNU32 Make*, make sure its bin path is added to your PATH environment variable.

For *Git Bash for Windows*, at the step of "Adjusting your PATH environment", please choose "Use Git and optional Unix tools from Windows Command Prompt".


## Developing the Provider

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (version 1.16+ is *required*). You'll also need to correctly setup a [GOPATH](http://golang.org/doc/code.html#GOPATH), as well as adding `$GOPATH/bin` to your `$PATH`.

Clone the repository to: `$GOPATH/src/github.com/schumann-it/terraform-provider-azureadb2c`

```sh
$ mkdir -p $GOPATH/src/github.com/terraform-providers; cd $GOPATH/src/github.com/terraform-providers
$ git clone https://github.com/schumann-it/terraform-provider-azureadb2c
```

Change to the clone directory and run `make tools` to install the dependent tooling needed to test and build the provider.

To compile the provider, run `make build`. This will build the provider and put the provider binary in the `$GOPATH/bin` directory.

```sh
$ make tools
...
$ make build
...
$ $GOPATH/bin/terraform-provider-azureadb2c
...
```

In order to test the provider, you can simply run `make test`.

```sh
$ make test
```

The majority of tests in the provider are Acceptance Tests - which provisions real resources in Azure. It's possible to run the entire acceptance test suite by running `make testacc` - however it's likely you'll want to run a subset, which you can do using a prefix, by running:

```
make testacc TESTARGS='-run=TestAccApplication'
```

The following ENV variables must be set in your shell prior to running acceptance tests:
- B2C_ARM_CLIENT_ID
- B2C_ARM_CLIENT_SECRET
- B2C_ARM_TENANT_ID

*NOTE:* Acceptance tests create real resources, and may cost money to run.
