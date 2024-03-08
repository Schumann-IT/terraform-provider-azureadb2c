package keyset

import (
	"context"
	"fmt"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/schumann-it/azure-b2c-sdk-for-go/msgraph"
)

var _ resource.Resource = &TrustframeworkKeySetResource{}

func NewTrustframeworkKeySetResource() resource.Resource {
	return &TrustframeworkKeySetResource{}
}

// TrustframeworkKeySetResource defines the resource implementation.
type TrustframeworkKeySetResource struct {
	client *msgraph.ServiceClient
}

func (r *TrustframeworkKeySetResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_trustframework_keyset"
}

func (r *TrustframeworkKeySetResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Trustframework KeySet",

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
						"must contain only alphanumeric characters",
					),
				},
			},
			"id": schema.StringAttribute{
				MarkdownDescription: "The id of the keyset",
				Computed:            true,
			},
			"metadata": schema.SingleNestedAttribute{
				Computed: true,
				Attributes: map[string]schema.Attribute{
					"odata_context": schema.StringAttribute{
						Computed: true,
					},
				},
			},
		},
	}
}

func (r *TrustframeworkKeySetResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
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

func (r *TrustframeworkKeySetResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data KeySet

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	ks, err := r.client.CreateKeySet(data.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("create keyset failed", err.Error())
		return
	}

	diags := data.Consume(ks)
	if diags != nil {
		resp.Diagnostics.Append(diags...)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *TrustframeworkKeySetResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data KeySet

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	ks, err := r.client.GetKeySet(data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("create keyset failed", err.Error())
		return
	}

	diags := data.Consume(ks)
	if diags != nil {
		resp.Diagnostics.Append(diags...)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *TrustframeworkKeySetResource) Update(_ context.Context, _ resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.Diagnostics.AddError("update not implemented", "keyset cannot be updated. please delete and create new.")
}

func (r *TrustframeworkKeySetResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data KeySet

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteKeySet(data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("delete keyset failed", err.Error())
	}
}
