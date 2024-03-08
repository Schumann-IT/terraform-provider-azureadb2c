package model

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

type KeySetKeyData struct {
	Kid types.String `tfsdk:"kid"`
	Use types.String `tfsdk:"use"`
	Kty types.String `tfsdk:"kty"`
	N   types.String `tfsdk:"n"`
	E   types.String `tfsdk:"e"`
}

type KeySetKeyResource struct {
	KeySetId types.String `tfsdk:"keyset_id"`
	Use      types.String `tfsdk:"use"`
	Type     types.String `tfsdk:"type"`
	Data     types.Object `tfsdk:"data"`
}

func (k *KeySetKeyResource) Consume(key models.TrustFrameworkKeyable) diag.Diagnostics {
	v, diags := buildKeyData(key)
	if diags != nil {
		return diags
	}

	k.Data = *v

	return nil
}

type KeySetKey struct {
	KeySetId types.String `tfsdk:"keyset_id"`
	Data     types.Object `tfsdk:"data"`
}

func (k *KeySetKey) Consume(key models.TrustFrameworkKeyable) diag.Diagnostics {
	v, diags := buildKeyData(key)
	if diags != nil {
		return diags
	}

	k.Data = *v

	return nil
}

func buildKeyData(key models.TrustFrameworkKeyable) (*basetypes.ObjectValue, diag.Diagnostics) {
	t := map[string]attr.Type{
		"kid": types.StringType,
		"use": types.StringType,
		"kty": types.StringType,
		"n":   types.StringType,
		"e":   types.StringType,
	}
	d := map[string]attr.Value{
		"kid": types.StringPointerValue(key.GetKid()),
		"use": types.StringPointerValue(key.GetUse()),
		"kty": types.StringPointerValue(key.GetKty()),
		"n":   types.StringPointerValue(key.GetN()),
		"e":   types.StringPointerValue(key.GetE()),
	}
	v, diags := types.ObjectValue(t, d)
	if diags != nil {
		return nil, diags
	}

	return &v, nil
}
