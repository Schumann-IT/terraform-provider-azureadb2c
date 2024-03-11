# authentication params are sources from B2C_* environment variables
provider "azureadb2c" {}

# provide authentication params directly
provider "azureadb2c" {
  tenant_id     = "<tenant_id>"
  client_id     = "<client_id>"
  client_secret = "<client_secret>"
}
