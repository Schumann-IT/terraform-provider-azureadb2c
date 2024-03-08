resource "azureadb2c_trustframework_keyset" "example" {
  name = "example"
}

resource "azureadb2c_trustframework_keyset_certificate" "example" {
  keyset_id   = azureadb2c_trustframework_keyset.example.id
  certificate = "<the certificate>"
  password    = "<the passphrase>"
}
