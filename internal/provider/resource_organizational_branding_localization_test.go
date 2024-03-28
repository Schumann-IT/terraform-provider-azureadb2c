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
					resource.TestCheckResourceAttr("azureadb2c_organizational_branding_localization.en", "background_color", "#008000"),
					resource.TestCheckResourceAttrSet("azureadb2c_organizational_branding_localization.en", "banner_logo_url"),
					resource.TestCheckResourceAttrSet("azureadb2c_organizational_branding_localization.en", "background_image_url"),
					resource.TestCheckResourceAttrSet("azureadb2c_organizational_branding_localization.en", "square_logo_light_url"),
					resource.TestCheckResourceAttrSet("azureadb2c_organizational_branding_localization.en", "square_logo_dark_url"),
					resource.TestCheckResourceAttr("azureadb2c_organizational_branding_localization.en", "username_hint_text", "Hint"),
				),
			},
			{
				Config: testAccDefaultOrganizationalBrandingLocalizationResource(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("azureadb2c_organizational_branding_localization.default", "background_color", "#ffffff"),
				),
			},
		},
	})
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
  background_color = "#008000"
  background_image = filebase64("./testdata/backgroundimage.png") 
  banner_logo = filebase64("./testdata/bannerlogo.jpg")
  square_logo_light = filebase64("./testdata/squarelogolight.jpg")
  square_logo_dark = filebase64("./testdata/squarelogodark.jpg")
  username_hint_text = "Hint"
}
`
}
