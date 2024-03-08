package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/schumann-it/azure-b2c-sdk-for-go/msgraph"
	"github.com/schumann-it/terraform-provider-azureadb2c/internal/model"
)

var _ provider.Provider = &AzureadB2c{}

type AzureadB2c struct {
	version string
}

func (p *AzureadB2c) Metadata(ctx context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "azureadb2c"
	resp.Version = p.version
}

func (p *AzureadB2c) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"tenant_id": schema.StringAttribute{
				MarkdownDescription: "The Tenant ID of the B2C directory which should be used.",
				Optional:            true,
			},
			"client_id": schema.StringAttribute{
				MarkdownDescription: "The Client ID which should be used for service principal authentication",
				Optional:            true,
			},
			"client_secret": schema.StringAttribute{
				MarkdownDescription: "he application password to use when authenticating as a Service Principal using a Client Secret",
				Optional:            true,
			},
		},
	}
}

func (p *AzureadB2c) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data model.Provider

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	cred, diags := data.Credential()
	if diags != nil {
		resp.Diagnostics.Append(diags...)
		return
	}

	sc, err := msgraph.NewClientWithCredential(cred)
	if err != nil {
		resp.Diagnostics.AddError("invalid credentials", err.Error())
		return
	}

	resp.DataSourceData = sc
	resp.ResourceData = sc
}

func (p *AzureadB2c) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewTrustframeworkKeySetResource,
		NewTrustframeworkKeySetCertificateResource,
		NewTrustframeworkKeySetKeyResource,
		NewTrustframeworkApplicationPatchResource,
		NewSamlApplicationPatchResource,
	}
}

func (p *AzureadB2c) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewTrustframeworkKeySetDataSource,
		NewTrustframeworkKeySetKeyDataSource,
		NewTrustframeworkKeySetCertificateDataSource,
	}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &AzureadB2c{
			version: version,
		}
	}
}
