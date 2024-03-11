package model

import (
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type Provider struct {
	TenantId     types.String `tfsdk:"tenant_id"`
	ClientId     types.String `tfsdk:"client_id"`
	ClientSecret types.String `tfsdk:"client_secret"`
}

func (m Provider) GetCredential() (azcore.TokenCredential, diag.Diagnostics) {
	var d diag.Diagnostics

	tid := m.getWithDefault(m.TenantId, "B2C_ARM_TENANT_ID")
	if tid == "" {
		d.AddError("missing tenant_id", "must be configured or provided via B2C_ARM_TENANT_ID env var")
	}
	cid := m.getWithDefault(m.ClientId, "B2C_ARM_CLIENT_ID")
	if cid == "" {
		d.AddError("missing client_id", "must be configured or provided via B2C_ARM_CLIENT_ID env var")
	}
	cs := m.getWithDefault(m.ClientSecret, "B2C_ARM_CLIENT_SECRET")
	if cs == "" {
		d.AddError("missing client_secret", "must be configured or provided via B2C_ARM_CLIENT_SECRET env var")
	}

	cred, err := azidentity.NewClientSecretCredential(tid, cid, cs, nil)
	if err != nil {
		d.AddError("invalid azure ad b2c client credentials", err.Error())
	}

	return cred, d
}

func (m Provider) getWithDefault(attr types.String, env string) string {
	v := attr.ValueString()
	if v != "" {
		return v
	}

	return os.Getenv(env)
}
