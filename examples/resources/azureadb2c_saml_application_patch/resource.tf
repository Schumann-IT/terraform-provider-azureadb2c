resource "azuread_application" "example" {
  display_name = "example"

  lifecycle {
    ignore_changes = [
      api,
    ]
  }
}

resource "azureadb2c_saml_application_patch" "example" {
  object_id         = azuread_application.example.object_id
  saml_metadata_url = "https://metadata.example.com"
}
