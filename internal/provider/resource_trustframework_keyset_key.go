package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/schumann-it/azure-b2c-sdk-for-go/msgraph"
	"github.com/schumann-it/terraform-provider-azureadb2c/internal/model"
)

var _ resource.Resource = &TrustframeworkKeySetKeyResource{}

func NewTrustframeworkKeySetKeyResource() resource.Resource {
	return &TrustframeworkKeySetKeyResource{}
}

type TrustframeworkKeySetKeyResource struct {
	client *msgraph.ServiceClient
}

func (r *TrustframeworkKeySetKeyResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_trustframework_keyset_key"
}

func (r *TrustframeworkKeySetKeyResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Trustframework key",

		Attributes: map[string]schema.Attribute{
			"keyset_id": schema.StringAttribute{
				MarkdownDescription: "The id of the keyset",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"use": schema.StringAttribute{
				MarkdownDescription: "Set to 'sig' for signing and 'enc' for encryption.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.OneOfCaseInsensitive("sig", "enc"),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"type": schema.StringAttribute{
				MarkdownDescription: "The type parameter identifies the cryptographic algorithm family used with the key, The valid values are RSA, OCT.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.OneOfCaseInsensitive("RSA", "OCT"),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"data": schema.SingleNestedAttribute{
				MarkdownDescription: "Trustframework Key",
				Attributes: map[string]schema.Attribute{
					"kid": schema.StringAttribute{
						MarkdownDescription: "The id of the key",
						Computed:            true,
					},
					"use": schema.StringAttribute{
						MarkdownDescription: "What the key is used for.",
						Computed:            true,
					},
					"kty": schema.StringAttribute{
						MarkdownDescription: "The kty identifies the cryptographic algorithm family used with the key.",
						Computed:            true,
					},
					"n": schema.StringAttribute{
						MarkdownDescription: "The n",
						Computed:            true,
					},
					"e": schema.StringAttribute{
						MarkdownDescription: "The e",
						Computed:            true,
					},
				},
				Computed:  true,
				Sensitive: true,
			},
		},
	}
}

func (r *TrustframeworkKeySetKeyResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *TrustframeworkKeySetKeyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data model.KeySetKeyResource

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	key, err := r.client.GenerateKey(data.KeySetId.ValueString(), data.Use.ValueString(), data.Type.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("create key failed", err.Error())
		return
	}

	diags := data.Consume(key)
	if diags != nil {
		resp.Diagnostics.Append(diags...)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *TrustframeworkKeySetKeyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data model.KeySetKeyResource

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	set, err := r.client.GetKeySet(data.KeySetId.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("read key failed", err.Error())
		return
	}

	diags := data.Consume(set.GetKeys()[0])
	if diags != nil {
		resp.Diagnostics.Append(diags...)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *TrustframeworkKeySetKeyResource) Update(_ context.Context, _ resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.Diagnostics.AddWarning("cannot update", "key cannot be updated.. please delete and create new keyset.")
}

func (r *TrustframeworkKeySetKeyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
}
