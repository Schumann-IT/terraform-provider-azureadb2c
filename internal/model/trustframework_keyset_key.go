package model

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type KeySetKey struct {
	KeySet types.Object `tfsdk:"key_set"`
	Use    types.String `tfsdk:"use"`
	Type   types.String `tfsdk:"type"`
}
