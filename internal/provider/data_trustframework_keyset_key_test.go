package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccTrustframeworkKeySetKeyDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccTrustframeworkKeySetKeyDataSource("SigningKeyContainer"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.azureadb2c_trustframework_keyset_key.test", "keyset_id", "B2C_1A_SigningKeyContainer"),
					resource.TestCheckResourceAttr("data.azureadb2c_trustframework_keyset_key.test", "data.use", "sig"),
					resource.TestCheckResourceAttr("data.azureadb2c_trustframework_keyset_key.test", "data.kty", "RSA"),
				),
			},
		},
	})
}

func testAccTrustframeworkKeySetKeyDataSource(name string) string {
	return fmt.Sprintf(`
resource "azureadb2c_trustframework_keyset" "test" {
  name = %[1]q
}

resource "azureadb2c_trustframework_keyset_key" "test" {
  keyset_id = azureadb2c_trustframework_keyset.test.id
  use = "sig"	
  type = "RSA"	
}

data "azureadb2c_trustframework_keyset_key" "test" {
  keyset_id = azureadb2c_trustframework_keyset_key.test.keyset_id
}
`, name)
}
