package model

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
	Keys     types.List   `tfsdk:"keys"`
}

type KeySetMetadata struct {
	Context types.String `tfsdk:"odata_context"`
}

type KeySetData struct {
	Kid types.String `tfsdk:"kid"`
}

func (ks *KeySet) Consume(keySet models.TrustFrameworkKeySetable) diag.Diagnostics {
	ks.Id = types.StringPointerValue(keySet.GetId())

	diags := ks.ConsumeMetadata(keySet.GetAdditionalData())
	if diags != nil {
		return diags
	}

	return ks.ConsumeKeys(keySet.GetKeys())
}

func (ks *KeySet) ConsumeMetadata(data map[string]interface{}) diag.Diagnostics {
	t := map[string]attr.Type{
		"odata_context": types.StringType,
	}
	d := map[string]attr.Value{
		"odata_context": types.StringPointerValue(data["@odata.context"].(*string)),
	}
	v, diags := types.ObjectValue(t, d)
	if diags != nil {
		return diags
	}

	ks.Metadata = v

	return nil
}

func (ks *KeySet) ConsumeKeys(keys []models.TrustFrameworkKeyable) diag.Diagnostics {
	var l []attr.Value

	t := map[string]attr.Type{
		"kid": types.StringType,
	}
	for _, key := range keys {
		d := map[string]attr.Value{
			"kid": types.StringPointerValue(key.GetKid()),
		}
		v, diags := types.ObjectValue(t, d)
		if diags != nil {
			return diags
		}
		l = append(l, v)
	}

	k, diags := types.ListValue(types.ObjectType{AttrTypes: t}, l)
	if diags != nil {
		return diags
	}

	ks.Keys = k

	return nil
}
