package provider

import (
	"context"
	"fmt"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/schumann-it/azure-b2c-sdk-for-go/msgraph"
	"github.com/schumann-it/terraform-provider-azureadb2c/internal/model"
)

var _ resource.Resource = &TrustframeworkKeySetResource{}

func NewTrustframeworkKeySetResource() resource.Resource {
	return &TrustframeworkKeySetResource{}
}

type TrustframeworkKeySetResource struct {
	client *msgraph.ServiceClient
}

func (r *TrustframeworkKeySetResource) Metadata(_ context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_trustframework_keyset"
}

func (r *TrustframeworkKeySetResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Represents a trust framework keyset/policy key.",

		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the keyset",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(`^[a-zA-Z]+$`),
						"must only contain only alphanumeric characters",
					),
				},
			},
			"id": schema.StringAttribute{
				MarkdownDescription: "Unique identifier of the trustframework keyset",
				Computed:            true,
			},
			"metadata": schema.SingleNestedAttribute{
				MarkdownDescription: "metadata of the trustframework keyset",
				Computed:            true,
				Attributes: map[string]schema.Attribute{
					"odata_context": schema.StringAttribute{
						MarkdownDescription: "The context",
						Computed:            true,
					},
				},
			},
			"keys": schema.ListAttribute{
				ElementType: types.ObjectType{
					AttrTypes: map[string]attr.Type{
						"kid": types.StringType,
					},
				},
				Computed:  true,
				Sensitive: true,
			},
		},
	}
}

func (r *TrustframeworkKeySetResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*msgraph.ServiceClient)
	if !ok {
		resp.Diagnostics.AddError(
			"provider error",
			fmt.Sprintf("Expected *msgraph.ServiceClient, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}

func (r *TrustframeworkKeySetResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data model.KeySet

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	set, err := r.client.CreateKeySet(data.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("create keyset failed", err.Error())
		return
	}

	diags := data.Consume(set)
	if diags != nil {
		resp.Diagnostics.Append(diags...)
		return
	}

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *TrustframeworkKeySetResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data model.KeySet

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	set, err := r.client.GetKeySet(data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("get keyset failed", err.Error())
		return
	}

	diags := data.Consume(set)
	if diags != nil {
		resp.Diagnostics.Append(diags...)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *TrustframeworkKeySetResource) Update(_ context.Context, _ resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.Diagnostics.AddWarning("update not implemented", "keyset cannot be updated. please delete and create new.")
}

func (r *TrustframeworkKeySetResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data model.KeySet

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteKeySet(data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("delete keyset failed", err.Error())
	}
}
