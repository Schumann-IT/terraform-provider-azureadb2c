package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/schumann-it/terraform-provider-azureadb2c/internal/acceptance"
)

func TestAccTrustframeworkKeySetCertificateResource(t *testing.T) {
	names := acceptance.RandAlphanumericStrings(1, 10)

	var expected []map[string]string
	for _, n := range names {
		expected = append(expected, map[string]string{
			"name": n,
			"id":   fmt.Sprintf("B2C_1A_%s", n),
		})
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccTrustframeworkKeySetCertificateResource(expected[0]["name"]),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("azureadb2c_trustframework_keyset_certificate.test", "key_set.id", expected[0]["id"]),
					resource.TestCheckResourceAttr("azureadb2c_trustframework_keyset_certificate.test", "key_set.keys.0.kty", "RSA"),
				),
			},
			// ImportState testing
			{
				ResourceName:                         "azureadb2c_trustframework_keyset_certificate.test",
				ImportState:                          true,
				ImportStateId:                        expected[0]["id"],
				ImportStateVerify:                    true,
				ImportStateVerifyIdentifierAttribute: "key_set.id",
				// certificate and password are not returned by the api, so we need to ignore these fields
				ImportStateVerifyIgnore: []string{"certificate", "password"},
			},
		},
	})
}

func testAccTrustframeworkKeySetCertificateResource(name string) string {
	return fmt.Sprintf(`
resource "azureadb2c_trustframework_keyset_certificate" "test" {
  key_set = {
	name = %[1]q
  }
	
  certificate = "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"	
  password = "Trigger.07"	
}
`, name)
}
