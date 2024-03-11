package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/schumann-it/azure-b2c-sdk-for-go/msgraph"
	"github.com/schumann-it/terraform-provider-azureadb2c/internal/model"
)

var _ datasource.DataSource = &TrustframeworkKeySetDataSource{}

func NewTrustframeworkKeySetDataSource() datasource.DataSource {
	return &TrustframeworkKeySetDataSource{}
}

type TrustframeworkKeySetDataSource struct {
	client *msgraph.ServiceClient
}

func (d *TrustframeworkKeySetDataSource) Metadata(_ context.Context, request datasource.MetadataRequest, response *datasource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_trustframework_keyset"
}

func (d *TrustframeworkKeySetDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = model.KeySetDataSourceSchema
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

func (d *TrustframeworkKeySetDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var keySet model.KeySet

	resp.Diagnostics.Append(req.Config.Get(ctx, &keySet)...)
	if resp.Diagnostics.HasError() {
		return
	}

	set, err := d.client.GetKeySet(keySet.GetId())
	if err != nil {
		resp.Diagnostics.AddError("read keyset failed", err.Error())
		return
	}

	diags := keySet.Consume(set)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &keySet)...)
}
