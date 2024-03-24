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
	Id              types.String `tfsdk:"id"`
	BackgroundColor types.String `tfsdk:"background_color"`
	SignInPageText  types.String `tfsdk:"sign_in_page_text"`
	BannerLogo      types.String `tfsdk:"banner_logo"`
	BannerLogoUrl   types.String `tfsdk:"banner_logo_url"`

	bl []byte
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

	cdn := b.GetCdnList()
	if len(cdn) > 0 && b.GetBannerLogoRelativeUrl() != nil {
		o.BannerLogoUrl = types.StringValue(fmt.Sprintf("https://%s/%s", b.GetCdnList()[0], to.String(b.GetBannerLogoRelativeUrl())))
	} else {
		o.BannerLogoUrl = types.StringNull()
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

	return diags
}

func (o *OrganizationalBrandingLocalization) GetBannerLogoBytes() []byte {
	return o.bl
}
