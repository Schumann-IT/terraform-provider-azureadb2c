resource "azureadb2c_trustframework_keyset_certificate" "example" {
  key_set = {
    id = "B2C_1A_ExampleContainer"
  }

  certificate = "<a base54 encoded pkcs12 certificate>"
  password    = "<the cert passphrase"
}
