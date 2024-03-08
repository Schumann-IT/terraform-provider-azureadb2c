provider "azureadb2c" {
  tenant_id     = "<tenant_id>"     # will be sourced from B2C_ARM_TENANT_ID if omitted
  client_id     = "<client_id>"     # will be sourced from B2C_ARM_CLIENT_ID if omitted
  client_secret = "<client_secret>" # will be sourced from B2C_ARM_CLIENT_SECRET if omitted
}
