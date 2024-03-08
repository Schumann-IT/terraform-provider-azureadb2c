// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package resources

import (
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/schumann-it/azure-b2c-sdk-for-go/msgraph"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &ApplicationPatch{}

//go:embed ../internal/model/TrustframeworkPatch.json
var trustframeworkApplicationPatch []byte

//go:embed application/SamlPatch.json
var samlApplicationPatch []byte

func NewApplicationPatch() resource.Resource {
	return &ApplicationPatch{}
}

// ApplicationPatch defines the resource implementation.
type ApplicationPatch struct {
	client *msgraph.ServiceClient
}

// ApplicationPatchModel describes the resource data model.
type ApplicationPatchModel struct {
	Id              types.String `tfsdk:"id"`
	Type            types.String `tfsdk:"type"`
	SamlMetadataUrl types.String `tfsdk:"saml_metadata_url"`
}

func (r *ApplicationPatch) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_application_patch"
}

func (r *ApplicationPatch) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "ApplicationPatch patch",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The ID of the application",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"type": schema.StringAttribute{
				MarkdownDescription: "A map that represents the patch",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.OneOf("trustframework", "saml"),
				},
			},
			"saml_metadata_url": schema.StringAttribute{
				MarkdownDescription: "The URL to the SAML metadata for the app.",
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.RegexMatches(regexp.MustCompile(`^.+\:\/\/.+$`), "format <protocol>://<uri>"),
				},
			},
		},
	}
}

func (r *ApplicationPatch) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *ApplicationPatch) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data ApplicationPatchModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	err := patchApplication(data, r.client)
	if err != nil {
		resp.Diagnostics.AddError("patch application failed", err.Error())
	}

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Trace(ctx, fmt.Sprintf("created %s application %s patch", data.Type.ValueString(), data.Id.ValueString()))

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ApplicationPatch) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data ApplicationPatchModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	app, err := r.client.GetApplication(data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("read application failed", err.Error())
	}

	if data.Type.ValueString() == "saml" {
		data.SamlMetadataUrl = types.StringPointerValue(app.GetSamlMetadataUrl())
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ApplicationPatch) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data ApplicationPatchModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	err := patchApplication(data, r.client)
	if err != nil {
		resp.Diagnostics.AddError("patch application failed", err.Error())
	}

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Trace(ctx, fmt.Sprintf("updated %s application %s patch", data.Type.ValueString(), data.Id.ValueString()))

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ApplicationPatch) Delete(_ context.Context, _ resource.DeleteRequest, _ *resource.DeleteResponse) {
}

func patchApplication(data ApplicationPatchModel, c *msgraph.ServiceClient) error {
	var patch map[string]interface{}
	switch data.Type.ValueString() {
	case "saml":
		if data.SamlMetadataUrl.ValueString() == "" {
			return fmt.Errorf("saml applications must patch saml_metadata_url")
		}
		_ = json.Unmarshal(samlApplicationPatch, &patch)
		patch["samlMetadataUrl"] = data.SamlMetadataUrl.ValueString()
	case "trustframework":
		_ = json.Unmarshal(trustframeworkApplicationPatch, &patch)
		patch["samlMetadataUrl"] = nil
	}
	return c.PatchApplication(data.Id.ValueString(), patch)
}
