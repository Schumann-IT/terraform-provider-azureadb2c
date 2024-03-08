package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccTrustframeworkApplicationPatchResource(t *testing.T) {
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
				Config: testAccTrustframeworkApplicationPatchResourceConfig("TestApplication"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("azureadb2c_trustframework_application_patch.test", "data.display_name", "TestApplication"),
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

resource "azureadb2c_trustframework_application_patch" "test" {
  object_id = azuread_application.test.object_id
}
`, name)
}
