// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package resources

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/schumann-it/terraform-provider-azureadb2c/internal"
)

func TestAccTrustframeworkKeySetResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { internal.testAccPreCheck(t) },
		ProtoV6ProviderFactories: internal.testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccTrustframeworkKeySetResource("TestContainer"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("azureadb2c_trustframework_keyset.test", "name", "TestContainer"),
				),
			},
		},
	})
}

func testAccTrustframeworkKeySetResource(id string) string {
	return fmt.Sprintf(`
resource "azureadb2c_trustframework_keyset" "test" {
  name = %[1]q
}
`, id)
}
