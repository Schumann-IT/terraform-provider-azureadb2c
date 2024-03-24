package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccOrganizationalBrandingLocalizationResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccOrganizationalBrandingLocalizationWithBannerResource(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("azureadb2c_organizational_branding_localization.en", "background_color", "#ffffff"),
					resource.TestCheckResourceAttrSet("azureadb2c_organizational_branding_localization.en", "banner_logo_url"),
				),
			},
			{
				Config: testAccDefaultOrganizationalBrandingLocalizationResource(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("azureadb2c_organizational_branding_localization.default", "background_color", "#ffffff"),
					resource.TestCheckNoResourceAttr("azureadb2c_organizational_branding_localization.default", "banner_logo_url"),
				),
			},
			{
				Config: testAccOrganizationalBrandingLocalizationResource(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("azureadb2c_organizational_branding_localization.de", "background_color", "#00a075"),
					resource.TestCheckNoResourceAttr("azureadb2c_organizational_branding_localization.de", "banner_logo_url"),
				),
			},
		},
	})
}

func testAccOrganizationalBrandingLocalizationResource() string {
	return `
resource "azureadb2c_organizational_branding_localization" "de" {
  id = "de-DE"
  background_color = "#00a075"
  sign_in_page_text = "Hello"
}
`
}

func testAccDefaultOrganizationalBrandingLocalizationResource() string {
	return `
resource "azureadb2c_organizational_branding_localization" "default" {
  id = "0"
  background_color = "#ffffff"
  sign_in_page_text = "Default"
}
`
}

func testAccOrganizationalBrandingLocalizationWithBannerResource() string {
	return `
resource "azureadb2c_organizational_branding_localization" "en" {
  id = "en-US"
  background_color = "#ffffff"
  banner_logo = filebase64("./testdata/example.jpg")
}
`
}
