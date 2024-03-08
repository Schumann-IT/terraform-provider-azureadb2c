package model

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

type ApplicationData struct {
	Id              types.String `tfsdk:"id"`
	AppId           types.String `tfsdk:"app_id"`
	SamlMetadataUrl types.String `tfsdk:"saml_metadata_url"`
	DisplayName     types.String `tfsdk:"saml_metadata_url"`
	IdentifierUris  types.List   `tfsdk:"identifier_uris"`
}

func buildApplicationData(app models.Applicationable) (*basetypes.ObjectValue, diag.Diagnostics) {
	t := map[string]attr.Type{
		"id":                types.StringType,
		"app_id":            types.StringType,
		"display_name":      types.StringType,
		"saml_metadata_url": types.StringType,
		"identifier_uris":   types.ListType{ElemType: types.StringType},
	}

	var identifierUrisElements []attr.Value
	for _, iu := range app.GetIdentifierUris() {
		identifierUrisElements = append(identifierUrisElements, types.StringValue(iu))
	}
	identifierUris, diags := types.ListValue(types.StringType, identifierUrisElements)
	if diags != nil {
		return nil, diags
	}

	d := map[string]attr.Value{
		"id":                types.StringPointerValue(app.GetId()),
		"app_id":            types.StringPointerValue(app.GetAppId()),
		"saml_metadata_url": types.StringPointerValue(app.GetSamlMetadataUrl()),
		"display_name":      types.StringPointerValue(app.GetDisplayName()),
		"identifier_uris":   identifierUris,
	}

	v, diags := types.ObjectValue(t, d)

	return &v, diags
}
