package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccTrustframeworkKeySetDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccTrustframeworkKeySetDataSourceById("SigningKeyContainer"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.azureadb2c_trustframework_keyset.test", "id", "B2C_1A_SigningKeyContainer"),
				),
			},
			{
				Config: testAccTrustframeworkKeySetDataSourceByName("SigningKeyContainer"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.azureadb2c_trustframework_keyset.test", "name", "SigningKeyContainer"),
				),
			},
		},
	})
}

func testAccTrustframeworkKeySetDataSourceById(name string) string {
	return fmt.Sprintf(`
resource "azureadb2c_trustframework_keyset_key" "test" {
  key_set = {
	name = %[1]q 
  }
  use = "enc"	
  type = "RSA"	
}

data "azureadb2c_trustframework_keyset" "test" {
  id = azureadb2c_trustframework_keyset_key.test.key_set.id
}
`, name)
}

func testAccTrustframeworkKeySetDataSourceByName(name string) string {
	return fmt.Sprintf(`
resource "azureadb2c_trustframework_keyset_key" "test" {
  key_set = {
	name = %[1]q 
  }
  use = "enc"	
  type = "RSA"	
}

data "azureadb2c_trustframework_keyset" "test" {
  id = azureadb2c_trustframework_keyset_key.test.key_set.id
}
`, name)
}
