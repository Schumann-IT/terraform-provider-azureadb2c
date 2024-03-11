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
		MarkdownDescription: `
This Provider provides a few resources that are missing in the [Azure Active Directory Provider](https://registry.terraform.io/providers/hashicorp/azuread/latest/docs)
but are required to fully automate [Azure AD B2C](https://learn.microsoft.com/en-us/azure/active-directory-b2c/overview). 

Detailed documentation regarding the Data Sources and Resources supported by the Azure AD B2C Provider can be found in the 
navigation to the left. 

Interested in the provider's latest features, or want to make sure you're up to date? Check out the 
[changelog](https://github.com/Schumann-IT/terraform-provider-azureadb2c/blob/main/CHANGELOG.md) for version information and release notes.

## Authenticating to Azure Active Directory

Authentication is currently only possible via service principal.  

The provided service principal must at least be granted the following permissions:

- Policy.Read.All
- Policy.ReadWrite.TrustFramework
- TrustFrameworkKeySet.Read.All
- TrustFrameworkKeySet.ReadWrite.All
- Application.Read.All
- Application.ReadWrite.All

Please see [Microsoft Graph permissions reference](https://learn.microsoft.com/en-us/graph/permissions-reference) and 
[Entra ID built-in roles](https://learn.microsoft.com/en-us/entra/identity/role-based-access-control/permissions-reference)
`,
		Attributes: map[string]schema.Attribute{
			"tenant_id": schema.StringAttribute{
				MarkdownDescription: "The Tenant ID which should be used. This can also be sourced from the `B2C_ARM_TENANT_ID` environment variable.",
				Optional:            true,
			},
			"client_id": schema.StringAttribute{
				MarkdownDescription: "The Client ID which should be used when authenticating as a service principal. This can also be sourced from the `B2C_ARM_CLIENT_ID` environment variable.",
				Optional:            true,
			},
			"client_secret": schema.StringAttribute{
				MarkdownDescription: "The application password to be used when authenticating using a client secret. This can also be sourced from the `B2C_ARM_CLIENT_SECRET` environment variable",
				Optional:            true,
			},
		},
	}
}

func (p *AzureadB2c) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data model.Provider

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	cred, diags := data.GetCredential()
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
		NewTrustframeworkKeySetCertificateResource,
		NewTrustframeworkKeySetKeyResource,
		NewApplicationPatchResource,
	}
}

func (p *AzureadB2c) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewTrustframeworkKeySetDataSource,
	}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &AzureadB2c{
			version: version,
		}
	}
}
