package model

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

type KeySetCertificateData struct {
	Kid types.String `tfsdk:"kid"`
	Exp types.Number `tfsdk:"exp"`
	E   types.String `tfsdk:"e"`
	X5c types.List   `tfsdk:"x5c"`
	Kty types.String `tfsdk:"kty"`
	N   types.String `tfsdk:"n"`
	X5t types.Number `tfsdk:"x5t"`
	Nbf types.Int64  `tfsdk:"nbf"`
}

type KeySetCertificateResource struct {
	KeySetId    types.String `tfsdk:"keyset_id"`
	Certificate types.String `tfsdk:"certificate"`
	Password    types.String `tfsdk:"password"`
	Data        types.Object `tfsdk:"data"`
}

func (k *KeySetCertificateResource) Consume(key models.TrustFrameworkKeyable) diag.Diagnostics {
	v, diags := buildCertificateData(key)
	if diags != nil {
		return diags
	}

	k.Data = *v

	return nil
}

type KeySetCertificate struct {
	KeySetId types.String `tfsdk:"keyset_id"`
	Data     types.Object `tfsdk:"data"`
}

func (k *KeySetCertificate) Consume(key models.TrustFrameworkKeyable) diag.Diagnostics {
	v, diags := buildCertificateData(key)
	if diags != nil {
		return diags
	}

	k.Data = *v

	return nil
}

func buildCertificateData(key models.TrustFrameworkKeyable) (*basetypes.ObjectValue, diag.Diagnostics) {
	t := map[string]attr.Type{
		"kid": types.StringType,
		"exp": types.Int64Type,
		"e":   types.StringType,
		"x5c": types.ListType{ElemType: types.StringType},
		"kty": types.StringType,
		"n":   types.StringType,
		"x5t": types.StringType,
		"nbf": types.Int64Type,
	}
	var x5cElements []attr.Value
	for _, x := range key.GetX5c() {
		x5cElements = append(x5cElements, types.StringValue(x))
	}
	x5c, diags := types.ListValue(types.StringType, x5cElements)
	if diags != nil {
		return nil, diags
	}

	d := map[string]attr.Value{
		"kid": types.StringPointerValue(key.GetKid()),
		"exp": types.Int64PointerValue(key.GetExp()),
		"e":   types.StringPointerValue(key.GetE()),
		"x5c": x5c,
		"kty": types.StringPointerValue(key.GetKty()),
		"n":   types.StringPointerValue(key.GetN()),
		"x5t": types.StringPointerValue(key.GetX5t()),
		"nbf": types.Int64PointerValue(key.GetNbf()),
	}

	v, diags := types.ObjectValue(t, d)

	return &v, diags
}
