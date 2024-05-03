package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/resourcevalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/models"
	"github.com/schumann-it/azure-b2c-sdk-for-go/msgraph"
	"github.com/schumann-it/terraform-provider-azureadb2c/internal/model"
)

var _ resource.Resource = &TrustframeworkKeySetKeyResource{}
var _ resource.ResourceWithConfigValidators = &TrustframeworkKeySetKeyResource{}
var _ resource.ResourceWithImportState = &TrustframeworkKeySetKeyResource{}

func NewTrustframeworkKeySetKeyResource() resource.Resource {
	return &TrustframeworkKeySetKeyResource{}
}

type TrustframeworkKeySetKeyResource struct {
	client *msgraph.ServiceClient
}

func (r *TrustframeworkKeySetKeyResource) ConfigValidators(_ context.Context) []resource.ConfigValidator {
	return []resource.ConfigValidator{
		resourcevalidator.ExactlyOneOf(
			path.MatchRelative().AtName("key_set").AtName("id"),
			path.MatchRelative().AtName("key_set").AtName("name"),
		),
	}
}

func (r *TrustframeworkKeySetKeyResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_trustframework_keyset_key"
}

func (r *TrustframeworkKeySetKeyResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages [Azure AD B2C Policy Keys](https://learn.microsoft.com/en-us/azure/active-directory-b2c/policy-keys-overview?pivots=b2c-custom-policy).",
		Attributes: map[string]schema.Attribute{
			"key_set": model.KeySetResourceSchema,
			"use": schema.StringAttribute{
				MarkdownDescription: "The use (public key use) parameter identifies the intended use of the public key. The use parameter is employed to indicate whether a public key is used for encrypting data or verifying the signature on data. Possible values are: sig (signature), enc (encryption)",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.OneOf("sig", "enc"),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"type": schema.StringAttribute{
				MarkdownDescription: "The kty (key type) parameter identifies the cryptographic algorithm family used with the key, Possible values are RSA, OCT.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.OneOf("RSA", "OCT"),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"secret": schema.StringAttribute{
				MarkdownDescription: "This is the field that is used to send the secret.",
				Optional:            true,
				Sensitive:           true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
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
	var key model.KeySetKey
	var keySet model.KeySet

	resp.Diagnostics.Append(req.Plan.Get(ctx, &key)...)
	resp.Diagnostics.Append(key.KeySet.As(ctx, &keySet, basetypes.ObjectAsOptions{})...)
	if resp.Diagnostics.HasError() {
		return
	}

	var set models.TrustFrameworkKeySetable
	var err error
	if key.Secret.IsNull() {
		set, err = r.client.GenerateKey(keySet.GetNameOrId(), key.Use.ValueString(), key.Type.ValueString())
	} else {
		set, err = r.client.UploadSecret(keySet.GetNameOrId(), key.Use.ValueString(), key.Secret.ValueString())
	}
	if err != nil {
		resp.Diagnostics.AddError("create keyset failed", err.Error())
		return
	}

	ksv, diags := keySet.GetObjectValue(set)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	key.KeySet = *ksv

	resp.Diagnostics.Append(resp.State.Set(ctx, &key)...)
}

func (r *TrustframeworkKeySetKeyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var key model.KeySetKey
	var keySet model.KeySet

	resp.Diagnostics.Append(req.State.Get(ctx, &key)...)
	resp.Diagnostics.Append(key.KeySet.As(ctx, &keySet, basetypes.ObjectAsOptions{})...)
	if resp.Diagnostics.HasError() {
		return
	}

	set, err := r.client.GetKeySet(keySet.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("read keyset failed", err.Error())
		return
	}

	ksv, diags := keySet.GetObjectValue(set)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	key.KeySet = *ksv

	diags = key.Consume(set)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &key)...)
}

func (r *TrustframeworkKeySetKeyResource) Update(_ context.Context, _ resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.Diagnostics.AddError("cannot update", "keyset cannot be updated.. please delete and create new keyset.")
}

func (r *TrustframeworkKeySetKeyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var key model.KeySetKey
	var keySet model.KeySet

	resp.Diagnostics.Append(req.State.Get(ctx, &key)...)
	resp.Diagnostics.Append(key.KeySet.As(ctx, &keySet, basetypes.ObjectAsOptions{})...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteKeySet(keySet.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("delete keyset failed", err.Error())
		return
	}

	_ = r.client.DeleteKeySet(fmt.Sprintf("%s.bak", keySet.Id.ValueString()))
}

func (r *TrustframeworkKeySetKeyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("key_set").AtName("id"), req, resp)
}
