package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/schumann-it/terraform-provider-azureadb2c/internal/acceptance"
)

func TestAccApplicationPatchResource(t *testing.T) {
	expected := acceptance.RandAlphanumericStrings(2, 10)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{
			"azuread": {
				Source: "hashicorp/azuread",
			},
		},
		Steps: []resource.TestStep{
			{
				Config: testAccTrustframeworkApplicationPatchResourceConfig(expected[0]),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("azureadb2c_application_patch.test", "data.display_name", expected[0]),
				),
			},
			{
				Config: testAccSamlApplicationPatchResourceConfig(expected[1], "https://metadata.example.com"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("azureadb2c_application_patch.saml", "data.display_name", expected[1]),
					resource.TestCheckResourceAttr("azureadb2c_application_patch.saml", "data.saml_metadata_url", "https://metadata.example.com"),
				),
			},
		},
	})
}

func testAccTrustframeworkApplicationPatchResourceConfig(name string) string {
	return fmt.Sprintf(`
resource "azuread_application" "test" {
  display_name = %[1]q

  api {
    mapped_claims_enabled          = null
    requested_access_token_version = 1
    known_client_applications = null
  }
}

resource "azureadb2c_application_patch" "test" {
  object_id = azuread_application.test.object_id
  patch_file = "./testdata/TrustframeworkApplicationPatch.json"
}
`, name)
}

func testAccSamlApplicationPatchResourceConfig(name, metadataUrl string) string {
	return fmt.Sprintf(`
resource "azuread_application" "saml" {
  display_name = %[1]q

  api {
    mapped_claims_enabled          = null
    requested_access_token_version = 2
    known_client_applications = null
  }
}

resource "azureadb2c_application_patch" "saml" {
  object_id = azuread_application.saml.object_id
  saml_metadata_url = %[2]q
  patch_file = "./testdata/SamlApplicationPatch.json"
}
`, name, metadataUrl)
}
