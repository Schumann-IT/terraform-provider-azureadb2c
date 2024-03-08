package provider

import (
	"context"
	"fmt"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework-validators/datasourcevalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/schumann-it/azure-b2c-sdk-for-go/msgraph"
	"github.com/schumann-it/terraform-provider-azureadb2c/internal/model"
)

var _ datasource.DataSource = &TrustframeworkKeySetDataSource{}
var _ datasource.DataSourceWithConfigValidators = &TrustframeworkKeySetDataSource{}

func NewTrustframeworkKeySetDataSource() datasource.DataSource {
	return &TrustframeworkKeySetDataSource{}
}

type TrustframeworkKeySetDataSource struct {
	client *msgraph.ServiceClient
}

func (d *TrustframeworkKeySetDataSource) Metadata(_ context.Context, request datasource.MetadataRequest, response *datasource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_trustframework_keyset"
}

func (d *TrustframeworkKeySetDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Represents a trust framework keyset/policy key.",

		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the keyset",
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(`^[a-zA-Z]+$`),
						"must only contain only alphanumeric characters",
					),
				},
			},
			"id": schema.StringAttribute{
				MarkdownDescription: "Unique identifier of the trustframework keyset",
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(`^B2C_1A_[a-zA-Z]+$`),
						"must be prefixed with B2C_1A_ must only only contain only alphanumeric characters",
					),
				},
			},
			"metadata": schema.SingleNestedAttribute{
				MarkdownDescription: "key set metadata",
				Computed:            true,
				Attributes: map[string]schema.Attribute{
					"odata_context": schema.StringAttribute{
						MarkdownDescription: "The context",
						Computed:            true,
					},
				},
			},
			"keys": schema.ListAttribute{
				MarkdownDescription: "Represents a list of JWK (JSON Web Key). TrustFrameworkKey is a JSON data structure that represents a cryptographic key. The structure of this resource follows the format defined in RFC 7517 Section 4.",
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

func (d *TrustframeworkKeySetDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	d.client = client
}

func (d *TrustframeworkKeySetDataSource) ConfigValidators(ctx context.Context) []datasource.ConfigValidator {
	return []datasource.ConfigValidator{
		datasourcevalidator.Conflicting(path.MatchRoot("name"), path.MatchRoot("id")),
	}
}

func (d *TrustframeworkKeySetDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data model.KeySet

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	var id string
	if data.Id.IsNull() {
		id = fmt.Sprintf("B2C_1A_%s", data.Name.ValueString())
	} else {
		id = data.Id.ValueString()
	}

	set, err := d.client.GetKeySet(id)
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
