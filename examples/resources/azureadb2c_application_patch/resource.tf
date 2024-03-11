data "azuread_application" "example" {
  display_name = "My First AzureAD Application"
}

resource "azureadb2c_application_patch" "example" {
  object_id  = azuread_application.example.object_id
  patch_file = format("%s/path/to/patch.json", path.module)
}
