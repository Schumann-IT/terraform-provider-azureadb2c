// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package resources

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/schumann-it/terraform-provider-azureadb2c/internal"
)

func TestAccApplicationPatch(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { internal.testAccPreCheck(t) },
		ProtoV6ProviderFactories: internal.testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccApplicationPatchConfigTrustframework("696a421d-bb8d-43c2-88d4-98e1176fc030"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("azureadb2c_application_patch.test", "id", "696a421d-bb8d-43c2-88d4-98e1176fc030"),
				),
			},
			// Create and Read testing
			{
				Config: testAccApplicationPatchConfigSaml("f7dc9dbd-c36a-4116-8f19-ac4df1096ed6"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("azureadb2c_application_patch.test", "id", "f7dc9dbd-c36a-4116-8f19-ac4df1096ed6"),
				),
			},
		},
	})
}

func testAccApplicationPatchConfigTrustframework(id string) string {
	return fmt.Sprintf(`
resource "azureadb2c_application_patch" "test" {
  id = %[1]q
  type = "trustframework"
}
`, id)
}

func testAccApplicationPatchConfigSaml(id string) string {
	return fmt.Sprintf(`
resource "azureadb2c_application_patch" "test" {
  id = %[1]q
  type = "saml"
  saml_metadata_url = "http://example.com"
}
`, id)
}
