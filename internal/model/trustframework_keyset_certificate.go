package model

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// KeySetCertificate represents a type that holds information about a certificate associated with a key set.
type KeySetCertificate struct {
	KeySet      types.Object `tfsdk:"key_set"`
	Certificate types.String `tfsdk:"certificate"`
	Password    types.String `tfsdk:"password"`
}
