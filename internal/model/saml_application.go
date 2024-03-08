package model

import (
	_ "embed"
	"encoding/json"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

//go:embed SamlPatch.json
var samlApplicationPatchSource []byte
var samlApplicationPatch map[string]interface{}

func init() {
	_ = json.Unmarshal(samlApplicationPatchSource, &samlApplicationPatch)
}

type SamlApplicationPatch struct {
	ObjectId        types.String `tfsdk:"object_id"`
	SamlMetadataUrl types.String `tfsdk:"saml_metadata_url"`
	Data            types.Object `tfsdk:"data"`
}

func (a *SamlApplicationPatch) GetPatch() map[string]interface{} {
	p := samlApplicationPatch
	p["samlMetadataUrl"] = a.SamlMetadataUrl.ValueString()
	return p
}

func (a *SamlApplicationPatch) Consume(app models.Applicationable) diag.Diagnostics {
	v, diags := buildApplicationData(app)
	if diags != nil {
		return diags
	}

	a.Data = *v

	return nil
}
