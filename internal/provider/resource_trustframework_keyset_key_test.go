package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/schumann-it/terraform-provider-azureadb2c/internal/acceptance"
)

func TestAccTrustframeworkKeySetKeyResource(t *testing.T) {
	var expected []map[string]string
	for _, n := range acceptance.RandAlphanumericStrings(2, 10) {
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
				Config: testAccTrustframeworkKeySetKeyResourceByIdSig(expected[0]["id"]),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("azureadb2c_trustframework_keyset_key.testsig", "key_set.name", expected[0]["name"]),
					resource.TestCheckResourceAttr("azureadb2c_trustframework_keyset_key.testsig", "key_set.keys.#", "1"),
					resource.TestCheckResourceAttr("azureadb2c_trustframework_keyset_key.testsig", "key_set.keys.0.use", "sig"),
					resource.TestCheckResourceAttr("azureadb2c_trustframework_keyset_key.testsig", "key_set.keys.0.kty", "RSA"),
				),
			},
			{
				Config: testAccTrustframeworkKeySetKeyResourceByNameEnc(expected[1]["name"]),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("azureadb2c_trustframework_keyset_key.testenc", "key_set.id", expected[1]["id"]),
					resource.TestCheckResourceAttr("azureadb2c_trustframework_keyset_key.testenc", "key_set.keys.#", "1"),
					resource.TestCheckResourceAttr("azureadb2c_trustframework_keyset_key.testenc", "key_set.keys.0.use", "enc"),
					resource.TestCheckResourceAttr("azureadb2c_trustframework_keyset_key.testenc", "key_set.keys.0.kty", "RSA"),
				),
			},
			// ImportState testing
			{
				ResourceName:                         "azureadb2c_trustframework_keyset_key.testenc",
				ImportState:                          true,
				ImportStateId:                        expected[1]["id"],
				ImportStateVerify:                    true,
				ImportStateVerifyIdentifierAttribute: "key_set.id",
			},
		},
	})
}

func testAccTrustframeworkKeySetKeyResourceByIdSig(id string) string {
	return fmt.Sprintf(`
resource "azureadb2c_trustframework_keyset_key" "testsig" {
  key_set = {
	id = %[1]q
  }
  use = "sig"
  type = "RSA"	
}
`, id)
}

func testAccTrustframeworkKeySetKeyResourceByNameEnc(name string) string {
	return fmt.Sprintf(`
resource "azureadb2c_trustframework_keyset_key" "testenc" {
  key_set = {
	name = %[1]q
  }
  use = "enc"
  type = "RSA"	
}
`, name)
}
