package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccApplicationPatchResource(t *testing.T) {
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
				Config: testAccTrustframeworkApplicationPatchResourceConfig("IdentityExperienceFramework"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("azureadb2c_application_patch.test", "data.display_name", "IdentityExperienceFramework"),
				),
			},
			{
				Config: testAccSamlApplicationPatchResourceConfig("SAMLApplication", "https://metadata.example.com"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("azureadb2c_application_patch.saml", "data.display_name", "SAMLApplication"),
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
