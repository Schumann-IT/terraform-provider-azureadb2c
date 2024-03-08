resource "azureadb2c_trustframework_keyset" "example" {
  name = "example"
}

resource "azureadb2c_trustframework_keyset_key" "example" {
  keyset_id = azureadb2c_trustframework_keyset.example.id
  use       = "enc"
  type      = "RSA"
}
