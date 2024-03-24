package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/resourcevalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/schumann-it/azure-b2c-sdk-for-go/msgraph"
	"github.com/schumann-it/terraform-provider-azureadb2c/internal/model"
)

var _ resource.Resource = &TrustframeworkKeySetCertificateResource{}
var _ resource.ResourceWithConfigValidators = &TrustframeworkKeySetCertificateResource{}
var _ resource.ResourceWithImportState = &TrustframeworkKeySetCertificateResource{}

func NewTrustframeworkKeySetCertificateResource() resource.Resource {
	return &TrustframeworkKeySetCertificateResource{}
}

type TrustframeworkKeySetCertificateResource struct {
	client *msgraph.ServiceClient
}

func (r *TrustframeworkKeySetCertificateResource) ConfigValidators(_ context.Context) []resource.ConfigValidator {
	return []resource.ConfigValidator{
		resourcevalidator.ExactlyOneOf(
			path.MatchRelative().AtName("key_set").AtName("id"),
			path.MatchRelative().AtName("key_set").AtName("name"),
		),
	}
}

func (r *TrustframeworkKeySetCertificateResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_trustframework_keyset_certificate"
}

func (r *TrustframeworkKeySetCertificateResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages certificates in Azure AD B2C Policy Keys. Example: [SAML application](https://learn.microsoft.com/en-us/azure/active-directory-b2c/saml-service-provider?tabs=windows&pivots=b2c-custom-policy#upload-the-certificate)",

		Attributes: map[string]schema.Attribute{
			"key_set": model.KeySetResourceSchema,
			"certificate": schema.StringAttribute{
				MarkdownDescription: "Base64 encoded pkcs12 certificate",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"password": schema.StringAttribute{
				MarkdownDescription: "The certificate passphrase",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
		},
	}
}

func (r *TrustframeworkKeySetCertificateResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *TrustframeworkKeySetCertificateResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var key model.KeySetCertificate
	var keySet model.KeySet

	resp.Diagnostics.Append(req.Plan.Get(ctx, &key)...)
	resp.Diagnostics.Append(key.KeySet.As(ctx, &keySet, basetypes.ObjectAsOptions{})...)
	if resp.Diagnostics.HasError() {
		return
	}

	set, err := r.client.UploadPkcs12(keySet.GetNameOrId(), key.Certificate.ValueString(), key.Password.ValueString())
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

func (r *TrustframeworkKeySetCertificateResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var key model.KeySetCertificate
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

	resp.Diagnostics.Append(resp.State.Set(ctx, &key)...)
}

func (r *TrustframeworkKeySetCertificateResource) Update(_ context.Context, _ resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.Diagnostics.AddError("cannot update", "keyset cannot be updated.. please delete and create new keyset.")
}

func (r *TrustframeworkKeySetCertificateResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var key model.KeySetCertificate
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

func (r *TrustframeworkKeySetCertificateResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("key_set").AtName("id"), req, resp)
}
