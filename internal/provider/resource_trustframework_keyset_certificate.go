package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/schumann-it/azure-b2c-sdk-for-go/msgraph"
	"github.com/schumann-it/terraform-provider-azureadb2c/internal/model"
)

var _ resource.Resource = &TrustframeworkKeySetCertificateResource{}

func NewTrustframeworkKeySetCertificateResource() resource.Resource {
	return &TrustframeworkKeySetCertificateResource{}
}

type TrustframeworkKeySetCertificateResource struct {
	client *msgraph.ServiceClient
}

func (r *TrustframeworkKeySetCertificateResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_trustframework_keyset_certificate"
}

func (r *TrustframeworkKeySetCertificateResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Represents a JWK (JSON Web Key). TrustFrameworkKey is a JSON data structure that represents a cryptographic key. The structure of this resource follows the format defined in RFC 7517 Section 4.",

		Attributes: map[string]schema.Attribute{
			"keyset_id": schema.StringAttribute{
				MarkdownDescription: "Unique identifier of the trustframework keyset",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
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
			"data": schema.SingleNestedAttribute{
				MarkdownDescription: "Represents a JWK (JSON Web Key). TrustFrameworkKey is a JSON data structure that represents a cryptographic key. The structure of this resource follows the format defined in RFC 7517 Section 4.",
				Attributes: map[string]schema.Attribute{
					"kid": schema.StringAttribute{
						MarkdownDescription: "The unique identifier for the key.",
						Computed:            true,
					},
					"exp": schema.Int64Attribute{
						MarkdownDescription: "This value is a NumericDate as defined in RFC 7519 (A JSON numeric value representing the number of seconds from 1970-01-01T00:00:00Z UTC until the specified UTC date/time, ignoring leap seconds.)",
						Computed:            true,
					},
					"e": schema.StringAttribute{
						MarkdownDescription: "RSA Key - public exponent",
						Computed:            true,
					},
					"x5c": schema.ListAttribute{
						ElementType:         types.StringType,
						MarkdownDescription: "The x5c (X.509 certificate chain) parameter contains a chain of one or more PKIX certificates RFC 5280.",
						Computed:            true,
					},
					"kty": schema.StringAttribute{
						MarkdownDescription: "The kty (key type) parameter identifies the cryptographic algorithm family used with the key, The valid values are rsa, oct.",
						Computed:            true,
					},
					"n": schema.StringAttribute{
						MarkdownDescription: "RSA Key - modulus",
						Computed:            true,
					},
					"x5t": schema.StringAttribute{
						MarkdownDescription: "The x5t (X.509 certificate SHA-1 thumbprint) parameter is a base64url-encoded SHA-1 thumbprint (also known as digest) of the DER encoding of an X.509 certificate RFC 5280.",
						Computed:            true,
					},
					"nbf": schema.Int64Attribute{
						MarkdownDescription: "This value is a NumericDate as defined in RFC 7519 (A JSON numeric value representing the number of seconds from 1970-01-01T00:00:00Z UTC until the specified UTC date/time, ignoring leap seconds.)",
						Computed:            true,
					},
				},
				Computed:  true,
				Sensitive: true,
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
	var data model.KeySetCertificateResource

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	key, err := r.client.UploadPkcs12(data.KeySetId.ValueString(), data.Certificate.ValueString(), data.Password.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("create certificate failed", err.Error())
		return
	}

	diags := data.Consume(key)
	if diags != nil {
		resp.Diagnostics.Append(diags...)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

}

func (r *TrustframeworkKeySetCertificateResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data model.KeySetCertificateResource

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	set, err := r.client.GetKeySet(data.KeySetId.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("read certificate failed", err.Error())
		return
	}

	diags := data.Consume(set.GetKeys()[0])
	if diags != nil {
		resp.Diagnostics.Append(diags...)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *TrustframeworkKeySetCertificateResource) Update(_ context.Context, _ resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.Diagnostics.AddWarning("cannot update", "certificate cannot be updated.. please delete and create new keyset.")
}

func (r *TrustframeworkKeySetCertificateResource) Delete(_ context.Context, _ resource.DeleteRequest, _ *resource.DeleteResponse) {
}
