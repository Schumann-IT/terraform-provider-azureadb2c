package model

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
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
}
