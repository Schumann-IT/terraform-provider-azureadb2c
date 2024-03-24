package provider

import (
	"context"
	_ "embed"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/resourcevalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/schumann-it/azure-b2c-sdk-for-go/msgraph"
	"github.com/schumann-it/terraform-provider-azureadb2c/internal/model"
)

var _ resource.Resource = &OrganizationalBrandingLocalization{}
var _ resource.ResourceWithConfigValidators = &OrganizationalBrandingLocalization{}

func NewOrganizationalBrandingLocalizationResource() resource.Resource {
	return &OrganizationalBrandingLocalization{}
}

type OrganizationalBrandingLocalization struct {
	client *msgraph.ServiceClient
}

func (r *OrganizationalBrandingLocalization) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_organizational_branding_localization"
}

func (r *OrganizationalBrandingLocalization) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: `
Resource that supports managing language-specific branding. While you can't change your original configuration's language, this resource allows you to create a new configuration for a different language.
`,

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "An identifier that represents the locale specified using culture names. Culture names follow the RFC 1766 standard in the format \"languagecode2-country/regioncode2\", where \"languagecode2\" is a lowercase two-letter code derived from ISO 639-1 and 'country/regioncode2' is an uppercase two-letter code derived from ISO 3166. For example, U.S. English is en-US. Use \"0\" to manage the default branding.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"background_color": schema.StringAttribute{
				MarkdownDescription: "Color that appears in place of the background image in low-bandwidth connections. We recommend that you use the primary color of your banner logo or your organization color. Specify this in hexadecimal format, for example, white is #FFFFFF.",
				Optional:            true,
			},
			"sign_in_page_text": schema.StringAttribute{
				MarkdownDescription: "Text that appears at the bottom of the sign-in box. Use this to communicate additional information, such as the phone number to your help desk or a legal statement. This text must be in Unicode format and not exceed 1024 characters.",
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(1024),
				},
			},
			"banner_logo": schema.StringAttribute{
				MarkdownDescription: "A banner version of your company logo that appears on the sign-in page. The allowed types are PNG or JPEG not larger than 36 × 245 pixels. We recommend using a transparent image with no padding around the logo.",
				Optional:            true,
			},
			"banner_logo_url": schema.StringAttribute{
				MarkdownDescription: "The URL to the banner logo",
				Computed:            true,
			},
		},
	}
}

func (r *OrganizationalBrandingLocalization) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *OrganizationalBrandingLocalization) ConfigValidators(ctx context.Context) []resource.ConfigValidator {
	return []resource.ConfigValidator{
		resourcevalidator.AtLeastOneOf(path.MatchRoot("background_color"), path.MatchRoot("sign_in_page_text")),
	}
}

func (r *OrganizationalBrandingLocalization) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data model.OrganizationalBrandingLocalization

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	exists := true
	b, _ := r.client.OrganizationClient.GetBrandingLocalization(data.Id.ValueString())
	if b == nil {
		exists = false
		b = r.client.OrganizationClient.NewBrandingLocalization(data.Id.ValueString())
	}

	diags := data.Populate(b)
	if diags != nil {
		resp.Diagnostics.Append(diags...)
		return
	}

	var err error
	if exists {
		_, err = r.client.OrganizationClient.UpdateBrandingLocalization(b)
	} else {
		_, err = r.client.OrganizationClient.CreateBrandingLocalization(b)
	}
	if err != nil {
		resp.Diagnostics.AddError("create organizational branding failed", err.Error())
		return
	}

	if len(data.GetBannerLogoBytes()) > 0 {
		err = r.client.OrganizationClient.UploadBannerLogo(data.Id.ValueString(), data.GetBannerLogoBytes())
		if err != nil {
			resp.Diagnostics.AddError("failed to upload banner logo", err.Error())
			return
		}
	} else {
		resp.Diagnostics.AddWarning("banner logo is not set", "removing the banner logo is currently not supported. please remove it manually")
	}

	b, err = r.client.OrganizationClient.GetBrandingLocalization(data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("get created organizational branding", err.Error())
		return
	}
	diags = data.Consume(b)
	if diags != nil {
		resp.Diagnostics.Append(diags...)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *OrganizationalBrandingLocalization) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data model.OrganizationalBrandingLocalization

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	b, err := r.client.OrganizationClient.GetBrandingLocalization(data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("get organizational branding failed", err.Error())
		return
	}

	diags := data.Consume(b)
	if diags != nil {
		resp.Diagnostics.Append(diags...)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *OrganizationalBrandingLocalization) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data model.OrganizationalBrandingLocalization

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	b, err := r.client.OrganizationClient.GetBrandingLocalization(data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("get organizational branding failed", err.Error())
		return
	}
	diags := data.Populate(b)
	if diags != nil {
		resp.Diagnostics.Append(diags...)
		return
	}

	_, err = r.client.OrganizationClient.UpdateBrandingLocalization(b)
	if err != nil {
		resp.Diagnostics.AddError("update organizational branding failed", err.Error())
		return
	}

	if len(data.GetBannerLogoBytes()) > 0 {
		err = r.client.OrganizationClient.UploadBannerLogo(data.Id.ValueString(), data.GetBannerLogoBytes())
		if err != nil {
			resp.Diagnostics.AddError("failed to upload banner logo", err.Error())
			return
		}
	} else {
		resp.Diagnostics.AddWarning("banner logo is not set", "removing the banner logo is currently not supported. please remove it manually")
	}

	b, err = r.client.OrganizationClient.GetBrandingLocalization(data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("get created organizational branding", err.Error())
		return
	}
	diags = data.Consume(b)
	if diags != nil {
		resp.Diagnostics.Append(diags...)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *OrganizationalBrandingLocalization) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data model.OrganizationalBrandingLocalization

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	if !data.BannerLogo.IsNull() {
		resp.Diagnostics.AddWarning("removing the banner logo is not supported", "Please remove the banner logo manually, before destroying this resource")
		return
	}

	err := r.client.OrganizationClient.DeleteBrandingLocalization(data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddWarning("delete organizational branding failed", err.Error())
		return
	}
}
