package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/schumann-it/terraform-provider-azureadb2c/internal/acceptance"
)

func TestAccTrustframeworkKeySetDataSource(t *testing.T) {
	var expected []map[string]string
	for _, n := range acceptance.RandAlphanumericStrings(2, 10) {
		expected = append(expected, map[string]string{
			"name": n,
			"id":   fmt.Sprintf("B2C_1A_%s", n),
		})
	}

	acctest.RandString(10)
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccTrustframeworkKeySetDataSourceById(expected[0]["name"]),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.azureadb2c_trustframework_keyset.test", "id", expected[0]["id"]),
					resource.TestCheckResourceAttr("data.azureadb2c_trustframework_keyset.test", "name", expected[0]["name"]),
				),
			},
			{
				Config: testAccTrustframeworkKeySetDataSourceByName(expected[1]["name"]),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.azureadb2c_trustframework_keyset.test", "id", expected[1]["id"]),
					resource.TestCheckResourceAttr("data.azureadb2c_trustframework_keyset.test", "name", expected[1]["name"]),
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
