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
				Config: testAccTrustframeworkKeySetKeyResource("B2C_1A_TestContainer", "sig"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("azureadb2c_trustframework_keyset_key.test", "key_set.name", "TestContainer"),
					resource.TestCheckResourceAttr("azureadb2c_trustframework_keyset_key.test", "key_set.keys.#", "1"),
					resource.TestCheckResourceAttr("azureadb2c_trustframework_keyset_key.test", "key_set.keys.0.use", "sig"),
					resource.TestCheckResourceAttr("azureadb2c_trustframework_keyset_key.test", "key_set.keys.0.kty", "RSA"),
				),
			},
			{
				Config: testAccTrustframeworkKeySetKeyResourceByName("TestContainer", "enc"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("azureadb2c_trustframework_keyset_key.test", "key_set.id", "B2C_1A_TestContainer"),
					resource.TestCheckResourceAttr("azureadb2c_trustframework_keyset_key.test", "key_set.keys.#", "1"),
					resource.TestCheckResourceAttr("azureadb2c_trustframework_keyset_key.test", "key_set.keys.0.use", "enc"),
					resource.TestCheckResourceAttr("azureadb2c_trustframework_keyset_key.test", "key_set.keys.0.kty", "RSA"),
				),
			},
		},
	})
}

func testAccTrustframeworkKeySetKeyResource(id, enc string) string {
	return fmt.Sprintf(`
resource "azureadb2c_trustframework_keyset_key" "test" {
  key_set = {
	id = %[1]q 
  }
  use = %[2]q	
  type = "RSA"	
}
`, id, enc)
}

func testAccTrustframeworkKeySetKeyResourceByName(name, enc string) string {
	return fmt.Sprintf(`
resource "azureadb2c_trustframework_keyset_key" "test" {
  key_set = {
	name = %[1]q 
  }
  use = %[2]q	
  type = "RSA"	
}
`, name, enc)
}
