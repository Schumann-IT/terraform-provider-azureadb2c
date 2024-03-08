package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccTrustframeworkKeySetKeyResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccTrustframeworkKeySetKeyResource("TestContainer"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("azureadb2c_trustframework_keyset_key.test", "keyset_id", "B2C_1A_TestContainer"),
					resource.TestCheckResourceAttr("azureadb2c_trustframework_keyset_key.test", "use", "sig"),
					resource.TestCheckResourceAttr("azureadb2c_trustframework_keyset_key.test", "type", "RSA"),
					resource.TestCheckResourceAttr("azureadb2c_trustframework_keyset_key.test", "data.use", "sig"),
					resource.TestCheckResourceAttr("azureadb2c_trustframework_keyset_key.test", "data.kty", "RSA"),
				),
			},
		},
	})
}

func testAccTrustframeworkKeySetKeyResource(name string) string {
	return fmt.Sprintf(`

resource "azureadb2c_trustframework_keyset" "test" {
  name = %[1]q
}

resource "azureadb2c_trustframework_keyset_key" "test" {
  keyset_id = azureadb2c_trustframework_keyset.test.id
  use = "sig"	
  type = "RSA"	
}
`, name)
}
