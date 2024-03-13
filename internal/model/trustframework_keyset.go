package model

import (
	"fmt"
	"strings"

	"github.com/Azure/go-autorest/autorest/to"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

const (
	KEY_SET_ID_PREFIX = "B2C_1A_%s"
)

var (
	keySetKeyType = map[string]attr.Type{
		"e":   types.StringType,
		"exp": types.Int64Type,
		"kid": types.StringType,
		"kty": types.StringType,
		"n":   types.StringType,
		"nbf": types.Int64Type,
		"x5c": types.ListType{ElemType: types.StringType},
		"x5t": types.StringType,
		"use": types.StringType,
	}

	keySetType = map[string]attr.Type{
		"id":   types.StringType,
		"name": types.StringType,
		"keys": types.ListType{ElemType: types.ObjectType{AttrTypes: keySetKeyType}},
	}
)

type (

	// KeySet struct defines the structure of a key set.
	KeySet struct {
		Id   types.String `tfsdk:"id"`
		Name types.String `tfsdk:"name"`
		Keys types.List   `tfsdk:"keys"`
	}

	// KeySetKeys struct defines the structure of a key set's keys, including the attributes: E, Exp, Kid, Kty, N, Nbf, X5c, X5t, and Use.
	KeySetKeys struct {
		E   types.String `tfsdk:"e"`
		Exp types.Number `tfsdk:"exp"`
		Kid types.String `tfsdk:"kid"`
		Kty types.String `tfsdk:"kty"`
		N   types.String `tfsdk:"n"`
		Nbf types.Number `tfsdk:"nbf"`
		X5c types.List   `tfsdk:"x5c"`
		X5t types.String `tfsdk:"x5t"`
		Use types.String `tfsdk:"use"`
	}
)

// GetId returns the value of the Id field if it is not null or unknown.
// Otherwise, the value of the Name field, prefixed with KEY_SET_ID_PREFIX is returned.
func (ks *KeySet) GetId() string {
	if ks.Id.IsNull() || ks.Id.IsUnknown() {
		return fmt.Sprintf(KEY_SET_ID_PREFIX, ks.Name.ValueString())
	}

	return ks.Id.ValueString()
}

// GetNameOrId returns the name of the KeySet if the Id field is null or unknown.
// Otherwise, it returns the value of the Id field.
func (ks *KeySet) GetNameOrId() string {
	if ks.Id.IsNull() || ks.Id.IsUnknown() {
		return ks.Name.ValueString()
	}

	return ks.Id.ValueString()
}

// Consume consumes a TrustFrameworkKeySetable and updates the KeySet object with the corresponding values.
func (ks *KeySet) Consume(keySet models.TrustFrameworkKeySetable) diag.Diagnostics {
	ks.Id = types.StringPointerValue(keySet.GetId())
	ks.Name = types.StringValue(strings.ReplaceAll(to.String(keySet.GetId()), KEY_SET_ID_PREFIX, ""))

	kv, diags := ks.buildKeysValue(keySet.GetKeys())
	if diags != nil {
		return diags
	}
	ks.Keys = *kv

	return diags
}

// GetObjectValue builds the value for the KeySetKey.KeySet and KeySetCertificate.KeySet fields.
func (ks *KeySet) GetObjectValue(keySet models.TrustFrameworkKeySetable) (*basetypes.ObjectValue, diag.Diagnostics) {
	kv, diags := ks.buildKeysValue(keySet.GetKeys())
	if diags != nil {
		return nil, diags
	}

	sd := map[string]attr.Value{
		"id":   types.StringPointerValue(keySet.GetId()),
		"name": types.StringValue(strings.ReplaceAll(to.String(keySet.GetId()), KEY_SET_ID_PREFIX, "")),
		"keys": kv,
	}

	sv, diags := types.ObjectValue(keySetType, sd)

	return &sv, diags
}

// buildKeysValue builds the value for the Keys field in the KeySet object based on the provided TrustFrameworkKeyables.
func (ks *KeySet) buildKeysValue(keys []models.TrustFrameworkKeyable) (*basetypes.ListValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	if len(keys) != 1 {
		diags.AddError("only one key allowed", fmt.Sprintf("a keyset can only contain one key, got %d", len(keys)))
		return nil, diags
	}

	var kl []attr.Value
	for _, key := range keys {
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
			"kty": types.StringPointerValue(key.GetKty()),
			"use": types.StringPointerValue(key.GetUse()),
			"n":   types.StringPointerValue(key.GetN()),
			"e":   types.StringPointerValue(key.GetE()),
			"exp": types.Int64PointerValue(key.GetExp()),
			"x5c": x5c,
			"x5t": types.StringPointerValue(key.GetX5t()),
			"nbf": types.Int64PointerValue(key.GetNbf()),
		}

		v, diags := types.ObjectValue(keySetKeyType, d)
		if diags != nil {
			return nil, diags
		}
		kl = append(kl, v)
	}
	kv, diags := types.ListValue(types.ObjectType{AttrTypes: keySetKeyType}, kl)

	return &kv, diags
}
