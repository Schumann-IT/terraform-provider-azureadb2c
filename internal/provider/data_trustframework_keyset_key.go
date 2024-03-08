package provider

import (
	"context"
	"fmt"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/schumann-it/azure-b2c-sdk-for-go/msgraph"
	"github.com/schumann-it/terraform-provider-azureadb2c/internal/model"
)

var _ datasource.DataSource = &TrustframeworkKeySetKeyDataSource{}

func NewTrustframeworkKeySetKeyDataSource() datasource.DataSource {
	return &TrustframeworkKeySetKeyDataSource{}
}

type TrustframeworkKeySetKeyDataSource struct {
	client *msgraph.ServiceClient
}

func (d *TrustframeworkKeySetKeyDataSource) Metadata(_ context.Context, request datasource.MetadataRequest, response *datasource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_trustframework_keyset_key"
}

func (d *TrustframeworkKeySetKeyDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Trustframework keyset key",

		Attributes: map[string]schema.Attribute{
			"keyset_id": schema.StringAttribute{
				MarkdownDescription: "The id of the keyset",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(`^B2C_1A_[a-zA-Z]+$`),
						"must be prefixed with B2C_1A_ must only only contain only alphanumeric characters",
					),
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

func (d *TrustframeworkKeySetKeyDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *TrustframeworkKeySetKeyDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data model.KeySetKey

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	set, err := d.client.GetKeySet(data.KeySetId.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("get keyset failed", err.Error())
		return
	}

	diags := data.Consume(set.GetKeys()[0])
	if diags != nil {
		resp.Diagnostics.Append(diags...)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
