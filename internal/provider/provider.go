// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/schumann-it/azure-b2c-sdk-for-go/msgraph"
)

// Ensure AzureadB2c satisfies various provider interfaces.
var _ provider.Provider = &AzureadB2c{}
var _ provider.ProviderWithFunctions = &AzureadB2c{}

// AzureadB2c defines the provider implementation.
type AzureadB2c struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// AzureadB2cProviderModel describes the provider data model.
type AzureadB2cProviderModel struct {
	TenantId     types.String `tfsdk:"tenant_id"`
	ClientId     types.String `tfsdk:"client_id"`
	ClientSecret types.String `tfsdk:"client_secret"`
}

func (m AzureadB2cProviderModel) AzureCredential(d *diag.Diagnostics) azcore.TokenCredential {
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
		d.AddError("invalid client credentials", err.Error())
	}

	return cred
}

func (m AzureadB2cProviderModel) getWithDefault(attr types.String, env string) string {
	v := attr.ValueString()
	if v != "" {
		return v
	}

	return os.Getenv(env)
}

func (p *AzureadB2c) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "azureadb2c"
	resp.Version = p.version
}

func (p *AzureadB2c) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"tenant_id": schema.StringAttribute{
				MarkdownDescription: "The tenant id of the B2C directory",
				Optional:            true,
			},
			"client_id": schema.StringAttribute{
				MarkdownDescription: "The client id of the service principal",
				Optional:            true,
			},
			"client_secret": schema.StringAttribute{
				MarkdownDescription: "The client secret of the service principal",
				Optional:            true,
			},
		},
	}
}

func (p *AzureadB2c) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data AzureadB2cProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	cred := data.AzureCredential(&resp.Diagnostics)
	sc, err := msgraph.NewClientWithCredential(cred)
	if err != nil {
		resp.Diagnostics.AddError("failed to initialize azure client", err.Error())
	}
	if resp.Diagnostics.HasError() {
		return
	}

	resp.DataSourceData = sc
	resp.ResourceData = sc
}

func (p *AzureadB2c) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewTrustframeworkKeySet,
		NewApplicationPatch,
	}
}

func (p *AzureadB2c) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{}
}

func (p *AzureadB2c) Functions(ctx context.Context) []func() function.Function {
	return []func() function.Function{}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &AzureadB2c{
			version: version,
		}
	}
}
