package provider

import (
	"context"
	_ "embed"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/schumann-it/azure-b2c-sdk-for-go/msgraph"
	"github.com/schumann-it/terraform-provider-azureadb2c/internal/model"
)

var _ resource.Resource = &ApplicationPatchResource{}

func NewApplicationPatchResource() resource.Resource {
	return &ApplicationPatchResource{}
}

type ApplicationPatchResource struct {
	client *msgraph.ServiceClient
}

func (r *ApplicationPatchResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_application_patch"
}

func (r *ApplicationPatchResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: `
Manages special requirements for application registration within Azure AD B2C when 
using [custom policies](https://learn.microsoft.com/en-us/azure/active-directory-b2c/user-flow-overview#custom-policies).

Please refer to the following examples:
- [Identity Experience Framework applications](https://learn.microsoft.com/en-us/azure/active-directory-b2c/tutorial-create-user-flows?pivots=b2c-custom-policy#register-identity-experience-framework-applications)
- [SAML applications](https://learn.microsoft.com/en-us/azure/active-directory-b2c/saml-service-provider?tabs=windows&pivots=b2c-custom-policy) 
- [Daemon applications](https://learn.microsoft.com/en-us/azure/active-directory-b2c/client-credentials-grant-flow?pivots=b2c-custom-policy) 

Other applications (like web and native apps) can still be configured via [Azure Active Directory Provider](https://registry.terraform.io/providers/hashicorp/azuread/latest/docs/resources/application) 
`,

		Attributes: map[string]schema.Attribute{
			"object_id": schema.StringAttribute{
				MarkdownDescription: "The object id of the application to be patched",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"patch_file": schema.StringAttribute{
				MarkdownDescription: "The path to the patch file. Must be an absolute path to a JSON file",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplaceIfConfigured(),
				},
			},
			"saml_metadata_url": schema.StringAttribute{
				MarkdownDescription: "The SAML metadata url",
				Optional:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplaceIfConfigured(),
				},
			},
			"data": schema.SingleNestedAttribute{
				MarkdownDescription: "identity experience framework app data",
				Computed:            true,
				Attributes: map[string]schema.Attribute{
					"id": schema.StringAttribute{
						MarkdownDescription: "The id of the application",
						Computed:            true,
					},
					"app_id": schema.StringAttribute{
						MarkdownDescription: "The application id (client id)",
						Computed:            true,
					},
					"saml_metadata_url": schema.StringAttribute{
						MarkdownDescription: "The saml metadata url",
						Computed:            true,
					},
					"display_name": schema.StringAttribute{
						MarkdownDescription: "The display name",
						Computed:            true,
					},
					"identifier_uris": schema.ListAttribute{
						ElementType:         types.StringType,
						MarkdownDescription: "The identifier uris",
						Computed:            true,
					},
				},
			},
		},
	}
}

func (r *ApplicationPatchResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*msgraph.ServiceClient)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *msgraph.ServiceClient, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}

func (r *ApplicationPatchResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data model.ApplicationPatch

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	p, diags := data.GetPatch()
	if diags != nil {
		resp.Diagnostics.Append(diags...)
		return
	}

	app, err := r.client.PatchApplication(data.ObjectId.ValueString(), p)
	if err != nil {
		resp.Diagnostics.AddError("patch application failed", err.Error())
		return
	}

	diags = data.Consume(app)
	if diags != nil {
		resp.Diagnostics.Append(diags...)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ApplicationPatchResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data model.ApplicationPatch

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	app, err := r.client.GetApplication(data.ObjectId.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("get application failed", err.Error())
		return
	}

	diags := data.Consume(app)
	if diags != nil {
		resp.Diagnostics.Append(diags...)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ApplicationPatchResource) Update(_ context.Context, _ resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.Diagnostics.AddError("update not implemented", "patch cannot be updated. please delete and create new.")
}

func (r *ApplicationPatchResource) Delete(_ context.Context, _ resource.DeleteRequest, _ *resource.DeleteResponse) {
}
