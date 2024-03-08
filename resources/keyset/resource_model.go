package keyset

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

type KeySet struct {
	Id       types.String `tfsdk:"id"`
	Name     types.String `tfsdk:"name"`
	Metadata types.Object `tfsdk:"metadata"`
}

func (ks *KeySet) Consume(keySet models.TrustFrameworkKeySetable) diag.Diagnostics {
	ks.Id = types.StringPointerValue(keySet.GetId())

	ad := keySet.GetAdditionalData()
	t := map[string]attr.Type{
		"odata_context": types.StringType,
	}
	d := map[string]attr.Value{
		"odata_context": types.StringPointerValue(ad["@odata.context"].(*string)),
	}
	v, diags := types.ObjectValue(t, d)
	if diags != nil {
		return diags
	}
	ks.Metadata = v

	return nil
}

type KeySetMetadata struct {
	Context types.String `tfsdk:"odata_context"`
}
