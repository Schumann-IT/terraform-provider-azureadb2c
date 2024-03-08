package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccTrustframeworkKeySetResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccTrustframeworkKeySetResource("TestContainer"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("azureadb2c_trustframework_keyset.test", "name", "TestContainer"),
					resource.TestCheckResourceAttr("azureadb2c_trustframework_keyset.test", "id", "B2C_1A_TestContainer"),
				),
			},
		},
	})
}

func testAccTrustframeworkKeySetResource(name string) string {
	return fmt.Sprintf(`
resource "azureadb2c_trustframework_keyset" "test" {
  name = %[1]q
}
`, name)
}
