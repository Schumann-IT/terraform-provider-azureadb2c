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

var _ resource.Resource = &TrustframeworkApplicationPatchResource{}

func NewTrustframeworkApplicationPatchResource() resource.Resource {
	return &TrustframeworkApplicationPatchResource{}
}

type TrustframeworkApplicationPatchResource struct {
	client *msgraph.ServiceClient
}

func (r *TrustframeworkApplicationPatchResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_trustframework_application_patch"
}

func (r *TrustframeworkApplicationPatchResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Patch the identity experience framework app",

		Attributes: map[string]schema.Attribute{
			"object_id": schema.StringAttribute{
				MarkdownDescription: "The object if of the application to be patched",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"data": schema.SingleNestedAttribute{
				MarkdownDescription: "identity experience framework app data",
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
				Computed:  true,
				Sensitive: true,
			},
		},
	}
}

func (r *TrustframeworkApplicationPatchResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *TrustframeworkApplicationPatchResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data model.TrustframeworkApplicationPatch

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.PatchApplication(data.ObjectId.ValueString(), data.GetPatch())
	if err != nil {
		resp.Diagnostics.AddError("patch application failed", err.Error())
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

func (r *TrustframeworkApplicationPatchResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data model.TrustframeworkApplicationPatch

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

func (r *TrustframeworkApplicationPatchResource) Update(_ context.Context, _ resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.Diagnostics.AddWarning("update not implemented", "patch cannot be updated. please delete and create new.")
}

func (r *TrustframeworkApplicationPatchResource) Delete(_ context.Context, _ resource.DeleteRequest, _ *resource.DeleteResponse) {
}
