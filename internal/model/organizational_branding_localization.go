package model

import (
	_ "embed"
	"encoding/base64"
	"fmt"

	"github.com/Azure/go-autorest/autorest/to"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

type OrganizationalBrandingLocalization struct {
	Id                 types.String `tfsdk:"id"`
	BackgroundColor    types.String `tfsdk:"background_color"`
	BackgroundImage    types.String `tfsdk:"background_image"`
	BackgroundImageUrl types.String `tfsdk:"background_image_url"`
	BannerLogo         types.String `tfsdk:"banner_logo"`
	BannerLogoUrl      types.String `tfsdk:"banner_logo_url"`
	SignInPageText     types.String `tfsdk:"sign_in_page_text"`
	SquareLogoLight    types.String `tfsdk:"square_logo_light"`
	SquareLogoLightUrl types.String `tfsdk:"square_logo_light_url"`
	SquareLogoDark     types.String `tfsdk:"square_logo_dark"`
	SquareLogoDarkUrl  types.String `tfsdk:"square_logo_dark_url"`
	UsernameHintText   types.String `tfsdk:"username_hint_text"`

	bi  []byte
	bl  []byte
	sll []byte
	sld []byte
}

func (o *OrganizationalBrandingLocalization) Consume(b models.OrganizationalBrandingLocalizationable) diag.Diagnostics {
	if to.String(b.GetBackgroundColor()) == "" {
		o.BackgroundColor = types.StringNull()
	} else {
		o.BackgroundColor = types.StringPointerValue(b.GetBackgroundColor())
	}

	if to.String(b.GetSignInPageText()) == "" {
		o.SignInPageText = types.StringNull()
	} else {
		o.SignInPageText = types.StringPointerValue(b.GetSignInPageText())
	}

	if to.String(b.GetUsernameHintText()) == "" {
		o.UsernameHintText = types.StringNull()
	} else {
		o.UsernameHintText = types.StringPointerValue(b.GetUsernameHintText())
	}

	cdn := b.GetCdnList()
	if len(cdn) > 0 && b.GetBannerLogoRelativeUrl() != nil {
		o.BannerLogoUrl = types.StringValue(fmt.Sprintf("https://%s/%s", b.GetCdnList()[0], to.String(b.GetBannerLogoRelativeUrl())))
	} else {
		o.BannerLogoUrl = types.StringNull()
	}
	if len(cdn) > 0 && b.GetBackgroundImageRelativeUrl() != nil {
		o.BackgroundImageUrl = types.StringValue(fmt.Sprintf("https://%s/%s", b.GetCdnList()[0], to.String(b.GetBackgroundImageRelativeUrl())))
	} else {
		o.BackgroundImageUrl = types.StringNull()
	}
	if len(cdn) > 0 && b.GetSquareLogoRelativeUrl() != nil {
		o.SquareLogoLightUrl = types.StringValue(fmt.Sprintf("https://%s/%s", b.GetCdnList()[0], to.String(b.GetSquareLogoRelativeUrl())))
	} else {
		o.SquareLogoLightUrl = types.StringNull()
	}
	if len(cdn) > 0 && b.GetSquareLogoDarkRelativeUrl() != nil {
		o.SquareLogoDarkUrl = types.StringValue(fmt.Sprintf("https://%s/%s", b.GetCdnList()[0], to.String(b.GetSquareLogoDarkRelativeUrl())))
	} else {
		o.SquareLogoDarkUrl = types.StringNull()
	}

	return nil
}

func (o *OrganizationalBrandingLocalization) Populate(b models.OrganizationalBrandingLocalizationable) diag.Diagnostics {
	var diags diag.Diagnostics

	if to.String(b.GetBackgroundColor()) != o.BackgroundColor.ValueString() {
		b.SetBackgroundColor(o.BackgroundColor.ValueStringPointer())
	}

	if to.String(b.GetSignInPageText()) != o.SignInPageText.ValueString() {
		b.SetSignInPageText(o.SignInPageText.ValueStringPointer())
	}

	if to.String(b.GetUsernameHintText()) != o.UsernameHintText.ValueString() {
		b.SetUsernameHintText(o.UsernameHintText.ValueStringPointer())
	}

	if !o.BannerLogo.IsNull() {
		t := make([]byte, base64.StdEncoding.DecodedLen(len(o.BannerLogo.ValueString())))
		_, err := base64.StdEncoding.Decode(t, []byte(o.BannerLogo.ValueString()))
		if err != nil {
			diags.AddError("failed to upload banner logo", err.Error())
		}
		o.bl = t
	} else {
		o.bl = nil
	}

	if !o.BackgroundImage.IsNull() {
		t := make([]byte, base64.StdEncoding.DecodedLen(len(o.BackgroundImage.ValueString())))
		_, err := base64.StdEncoding.Decode(t, []byte(o.BackgroundImage.ValueString()))
		if err != nil {
			diags.AddError("failed to upload background image", err.Error())
		}
		o.bi = t
	} else {
		o.bi = nil
	}

	if !o.SquareLogoLight.IsNull() {
		t := make([]byte, base64.StdEncoding.DecodedLen(len(o.SquareLogoLight.ValueString())))
		_, err := base64.StdEncoding.Decode(t, []byte(o.SquareLogoLight.ValueString()))
		if err != nil {
			diags.AddError("failed to upload square logo (light)", err.Error())
		}
		o.sll = t
	} else {
		o.sll = nil
	}

	if !o.SquareLogoDark.IsNull() {
		t := make([]byte, base64.StdEncoding.DecodedLen(len(o.SquareLogoDark.ValueString())))
		_, err := base64.StdEncoding.Decode(t, []byte(o.SquareLogoDark.ValueString()))
		if err != nil {
			diags.AddError("failed to upload square logo (dark)", err.Error())
		}
		o.sld = t
	} else {
		o.sld = nil
	}

	return diags
}

func (o *OrganizationalBrandingLocalization) GetBannerLogoBytes() []byte {
	return o.bl
}
func (o *OrganizationalBrandingLocalization) GetBackgroundImageBytes() []byte {
	return o.bi
}
func (o *OrganizationalBrandingLocalization) GetSquareLogoLightBytes() []byte {
	return o.sll
}
func (o *OrganizationalBrandingLocalization) GetSquareLogoDarkBytes() []byte {
	return o.sld
}
