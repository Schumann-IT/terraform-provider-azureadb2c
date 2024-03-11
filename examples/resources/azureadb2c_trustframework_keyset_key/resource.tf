resource "azureadb2c_trustframework_keyset_key" "example" {
  key_set = {
    id = "B2C_1A_ExampleContainer"
  }

  use  = "sig" # or enc
  type = "RSA"
}
