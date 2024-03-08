package model

import (
	_ "embed"
	"encoding/json"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

//go:embed TrustframeworkPatch.json
var trustframeworkApplicationPatchSource []byte
var trustframeworkApplicationPatch map[string]interface{}

func init() {
	_ = json.Unmarshal(trustframeworkApplicationPatchSource, &trustframeworkApplicationPatch)
}

type TrustframeworkApplicationPatch struct {
	ObjectId types.String `tfsdk:"object_id"`
	Data     types.Object `tfsdk:"data"`
}

func (a *TrustframeworkApplicationPatch) GetPatch() map[string]interface{} {
	return trustframeworkApplicationPatch
}

func (a *TrustframeworkApplicationPatch) Consume(app models.Applicationable) diag.Diagnostics {
	v, diags := buildApplicationData(app)
	if diags != nil {
		return diags
	}

	a.Data = *v

	return nil
}
