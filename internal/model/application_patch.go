package model

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"os"
	"path"
	"path/filepath"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

type ApplicationPatch struct {
	ObjectId        types.String `tfsdk:"object_id"`
	SamlMetadataUrl types.String `tfsdk:"saml_metadata_url"`
	PatchFile       types.String `tfsdk:"patch_file"`
	Data            types.Object `tfsdk:"data"`
}

func (a *ApplicationPatch) GetPatch() (map[string]interface{}, diag.Diagnostics) {
	var diags diag.Diagnostics

	if a.PatchFile.IsNull() {
		return nil, diags
	}

	f := a.PatchFile.ValueString()
	if !path.IsAbs(f) {
		var err error
		f, err = filepath.Abs(f)
		if err != nil {
			diags.AddError("path must be absolute", fmt.Sprintf("expected absolute path, got: %s", f))
			return nil, diags
		}
	}

	b, err := os.ReadFile(f)
	if err != nil {
		diags.AddError("failed to read", fmt.Sprintf("cannot read patch file: %s", err.Error()))
		return nil, diags
	}

	var p map[string]interface{}
	err = json.Unmarshal(b, &p)
	if err != nil {
		diags.AddError("failed to parse json", fmt.Sprintf("cannot read patch file: %s", err.Error()))
		return nil, diags
	}

	if !a.SamlMetadataUrl.IsNull() {
		p["samlMetadataUrl"] = a.SamlMetadataUrl.ValueString()
	}

	return p, diags
}

func (a *ApplicationPatch) Consume(app models.Applicationable) diag.Diagnostics {
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
		return diags
	}

	d := map[string]attr.Value{
		"id":                types.StringPointerValue(app.GetId()),
		"app_id":            types.StringPointerValue(app.GetAppId()),
		"saml_metadata_url": types.StringPointerValue(app.GetSamlMetadataUrl()),
		"display_name":      types.StringPointerValue(app.GetDisplayName()),
		"identifier_uris":   identifierUris,
	}

	v, diags := types.ObjectValue(t, d)
	if diags != nil {
		return diags
	}

	a.Data = v

	return nil
}
