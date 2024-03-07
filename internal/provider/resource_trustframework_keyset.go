// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/Azure/go-autorest/autorest/to"
	"github.com/hashicorp/terraform-plugin-framework-validators/resourcevalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/models"
	"github.com/schumann-it/azure-b2c-sdk-for-go/keyset"
	"github.com/schumann-it/azure-b2c-sdk-for-go/msgraph"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &TrustframeworkKeySet{}

func NewTrustframeworkKeySet() resource.Resource {
	return &TrustframeworkKeySet{}
}

// TrustframeworkKeySet defines the resource implementation.
type TrustframeworkKeySet struct {
	client *msgraph.ServiceClient
}

// TrustframeworkKeySetModel describes the resource data model.
type TrustframeworkKeySetModel struct {
	Id          types.String `tfsdk:"id"`
	Key         types.Object `tfsdk:"key"`
	Certificate types.Object `tfsdk:"certificate"`
}

type TrustframeworkKeySetKeyModel struct {
	Id      types.String `tfsdk:"id"`
	Use     types.String `tfsdk:"use"`
	KeyType types.String `tfsdk:"kty"`
}

type TrustframeworkKeySetCertificateModel struct {
	Id       types.String `tfsdk:"id"`
	Body     types.String `tfsdk:"body"`
	Password types.String `tfsdk:"password"`
}

func (r *TrustframeworkKeySet) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_trustframework_keyset"
}

func (r *TrustframeworkKeySet) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Trustframework KeySet",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The name of the keyset",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(`^B2C_1A_[a-zA-Z]+$`),
						"must start with B2C_1A_ and must contain only lowercase alphanumeric characters",
					),
				},
			},
			"key": schema.SingleNestedAttribute{
				MarkdownDescription: "Trustframework Key, conflicts with certificate",
				Attributes: map[string]schema.Attribute{
					"id": schema.StringAttribute{
						MarkdownDescription: "The id of the key",
						Computed:            true,
					},
					"use": schema.StringAttribute{
						MarkdownDescription: "Set to 'sig' for signing and 'enc' for encryption. Conflicts with 'certificate'",
						Required:            true,
						Validators: []validator.String{
							stringvalidator.OneOf("sig", "enc"),
						},
					},
					"kty": schema.StringAttribute{
						MarkdownDescription: "The kty (key type) parameter identifies the cryptographic algorithm family used with the key, The valid values are RSA, OCT.",
						Required:            true,
						Validators: []validator.String{
							stringvalidator.OneOf("RSA", "OCT"),
						},
					},
				},
				Optional: true,
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.RequiresReplaceIfConfigured(),
				},
			},
			"certificate": schema.SingleNestedAttribute{
				MarkdownDescription: "Trustframework Certificate, conflicts with key",
				Attributes: map[string]schema.Attribute{
					"id": schema.StringAttribute{
						MarkdownDescription: "The id of the certificate",
						Computed:            true,
					},
					"body": schema.StringAttribute{
						MarkdownDescription: "The pem encoded pkcs12 certificate body",
						Required:            true,
						Sensitive:           true,
					},
					"password": schema.StringAttribute{
						MarkdownDescription: "The certificate password",
						Required:            true,
						Sensitive:           true,
					},
				},
				Optional: true,
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.RequiresReplaceIfConfigured(),
				},
			},
		},
	}
}

func (r *TrustframeworkKeySet) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *TrustframeworkKeySet) ConfigValidators(context.Context) []resource.ConfigValidator {
	return []resource.ConfigValidator{
		resourcevalidator.ExactlyOneOf(
			path.MatchRoot("key"),
			path.MatchRoot("certificate"),
		),
	}
}

func (r *TrustframeworkKeySet) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data TrustframeworkKeySetModel
	var keyData *TrustframeworkKeySetKeyModel
	var certData *TrustframeworkKeySetCertificateModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	ks := keyset.NewKeySet(data.Id.ValueString())
	if !data.Key.IsNull() {
		resp.Diagnostics.Append(data.Key.As(ctx, &keyData, basetypes.ObjectAsOptions{})...)
		t := strings.ToUpper(keyData.KeyType.ValueString())
		if t != "RSA" {
			resp.Diagnostics.AddError("invalid key type", "only RSA is supported so far")
		}
		ks.WithRsaKey(keyData.Use.ValueString())
	}
	if !data.Certificate.IsNull() {
		resp.Diagnostics.Append(data.Certificate.As(ctx, &certData, basetypes.ObjectAsOptions{})...)
		ks.WithCertificate(&keyset.Certificate{
			Key:      certData.Body.ValueString(),
			Password: certData.Password.ValueString(),
		})
	}

	if resp.Diagnostics.HasError() {
		return
	}

	set, err := r.client.CreateKeySet(ks)
	if err != nil {
		resp.Diagnostics.AddError("create failed", err.Error())
		return
	}

	d := updateData(set, &data, certData)
	if d != nil {
		resp.Diagnostics.Append(d...)
		return
	}

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Trace(ctx, fmt.Sprintf("key set %s created", to.String(ks.Get().GetId())))

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *TrustframeworkKeySet) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data TrustframeworkKeySetModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	var certData *TrustframeworkKeySetCertificateModel
	resp.Diagnostics.Append(data.Certificate.As(ctx, &certData, basetypes.ObjectAsOptions{})...)

	if resp.Diagnostics.HasError() {
		return
	}

	set, err := r.client.GetKeySet(data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("create failed", err.Error())
		return
	}

	d := updateData(set, &data, certData)
	if d != nil {
		resp.Diagnostics.Append(d...)
		return
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *TrustframeworkKeySet) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.Diagnostics.AddError("cannot update", "key set cannot be updated.. please delete and create new.")
}

func (r *TrustframeworkKeySet) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data TrustframeworkKeySetModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	ks := keyset.NewKeySet(data.Id.ValueString())
	err := r.client.DeleteKeySet(ks)
	if err != nil {
		resp.Diagnostics.AddError("delete key set failed", err.Error())
	}
}

func updateData(set models.TrustFrameworkKeySetable, data *TrustframeworkKeySetModel, certData *TrustframeworkKeySetCertificateModel) diag.Diagnostics {
	for _, key := range set.GetKeys() {
		if !data.Key.IsNull() {
			elementTypes := map[string]attr.Type{
				"id":  types.StringType,
				"use": types.StringType,
				"kty": types.StringType,
			}
			elements := map[string]attr.Value{
				"id":  types.StringPointerValue(key.GetKid()),
				"use": types.StringPointerValue(key.GetUse()),
				"kty": types.StringPointerValue(key.GetKty()),
			}
			kd, diags := types.ObjectValue(elementTypes, elements)
			if diags != nil {
				return diags
			}
			data.Key = kd
		}
		if !data.Certificate.IsNull() {
			elementTypes := map[string]attr.Type{
				"id":       types.StringType,
				"body":     types.StringType,
				"password": types.StringType,
			}
			elements := map[string]attr.Value{
				"id":       types.StringPointerValue(key.GetKid()),
				"body":     types.StringValue(certData.Body.ValueString()),
				"password": types.StringValue(certData.Password.ValueString()),
			}
			cd, diags := types.ObjectValue(elementTypes, elements)
			if diags != nil {
				return diags
			}
			data.Certificate = cd
		}
	}

	return nil
}
