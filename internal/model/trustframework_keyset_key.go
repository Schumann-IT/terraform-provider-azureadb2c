package model

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// KeySetKey represents a key associated with a key set.
// It contains the following fields:
// - KeySet: the key set to which the key belongs (tfsdk:"key_set").
// - Use: the intended use of the key (tfsdk:"use").
// - Type: the type of the key (tfsdk:"type").
type KeySetKey struct {
	KeySet types.Object `tfsdk:"key_set"`
	Use    types.String `tfsdk:"use"`
	Type   types.String `tfsdk:"type"`
	Secret types.String `tfsdk:"secret"`
}

func (ks *KeySetKey) Consume(keySet models.TrustFrameworkKeySetable) diag.Diagnostics {
	var diags diag.Diagnostics

	keys := keySet.GetKeys()
	if len(keys) != 1 {
		diags.AddError("unexpected resource state", fmt.Sprintf("key len must be 1, got: %d", len(keys)))
		return diags
	}

	u := *keys[0].GetUse()
	ks.Use = types.StringValue(strings.ToLower(u))
	t := *keys[0].GetKty()
	ks.Type = types.StringValue(strings.ToUpper(t))

	return diags
}
