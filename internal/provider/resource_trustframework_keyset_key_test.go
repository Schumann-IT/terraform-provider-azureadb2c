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
				Config: testAccTrustframeworkKeySetKeyResourceByIdSig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("azureadb2c_trustframework_keyset_key.testsig", "key_set.name", "TestContainerSig"),
					resource.TestCheckResourceAttr("azureadb2c_trustframework_keyset_key.testsig", "key_set.keys.#", "1"),
					resource.TestCheckResourceAttr("azureadb2c_trustframework_keyset_key.testsig", "key_set.keys.0.use", "sig"),
					resource.TestCheckResourceAttr("azureadb2c_trustframework_keyset_key.testsig", "key_set.keys.0.kty", "RSA"),
				),
			},
			{
				Config: testAccTrustframeworkKeySetKeyResourceByNameEnc(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("azureadb2c_trustframework_keyset_key.testenc", "key_set.id", "B2C_1A_TestContainerEnc"),
					resource.TestCheckResourceAttr("azureadb2c_trustframework_keyset_key.testenc", "key_set.keys.#", "1"),
					resource.TestCheckResourceAttr("azureadb2c_trustframework_keyset_key.testenc", "key_set.keys.0.use", "enc"),
					resource.TestCheckResourceAttr("azureadb2c_trustframework_keyset_key.testenc", "key_set.keys.0.kty", "RSA"),
				),
			},
			// ImportState testing
			{
				ResourceName:                         "azureadb2c_trustframework_keyset_key.testenc",
				ImportState:                          true,
				ImportStateId:                        "B2C_1A_TestContainerEnc",
				ImportStateVerify:                    true,
				ImportStateVerifyIdentifierAttribute: "key_set.id",
			},
		},
	})
}

func testAccTrustframeworkKeySetKeyResourceByIdSig() string {
	return fmt.Sprintf(`
resource "azureadb2c_trustframework_keyset_key" "testsig" {
  key_set = {
	id = "B2C_1A_TestContainerSig" 
  }
  use = "sig"
  type = "RSA"	
}
`)
}

func testAccTrustframeworkKeySetKeyResourceByNameEnc() string {
	return fmt.Sprintf(`
resource "azureadb2c_trustframework_keyset_key" "testenc" {
  key_set = {
	name = "TestContainerEnc" 
  }
  use = "enc"
  type = "RSA"	
}
`)
}
