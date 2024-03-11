package model

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type KeySetCertificate struct {
	KeySet      types.Object `tfsdk:"key_set"`
	Certificate types.String `tfsdk:"certificate"`
	Password    types.String `tfsdk:"password"`
}
