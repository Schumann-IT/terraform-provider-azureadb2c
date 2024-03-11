package model

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	KeySetResourceSchema = schema.SingleNestedAttribute{
		MarkdownDescription: "Key set data",
		Optional:            true,
		Computed:            true,
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The id of the keyset",
				Optional:            true,
				Computed:            true,
				Validators: []validator.String{
					stringvalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("name")),
					KeySetIdValidator,
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the keyset",
				Optional:            true,
				Computed:            true,
				Validators: []validator.String{
					stringvalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("id")),
					KeySetNameValidator,
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"keys": schema.ListNestedAttribute{
				MarkdownDescription: "Represents a list of JWK (JSON Web Key). TrustFrameworkKey is a JSON data structure that represents a cryptographic key. The structure of this resource follows the format defined in RFC 7517 Section 4.",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"kid": schema.StringAttribute{
							MarkdownDescription: "The unique identifier for the key.",
							Computed:            true,
						},
						"kty": schema.StringAttribute{
							MarkdownDescription: "The kty (key type) parameter identifies the cryptographic algorithm family used with the key, The valid values are rsa, oct.",
							Computed:            true,
						},
						"use": schema.StringAttribute{
							MarkdownDescription: "The use (public key use) parameter identifies the intended use of the public key. The use parameter is employed to indicate whether a public key is used for encrypting data or verifying the signature on data. Possible values are: sig (signature), enc (encryption)",
							Computed:            true,
						},
						"n": schema.StringAttribute{
							MarkdownDescription: "RSA Key - modulus",
							Computed:            true,
							Sensitive:           true,
						},
						"e": schema.StringAttribute{
							MarkdownDescription: "RSA Key - public exponent",
							Computed:            true,
							Sensitive:           true,
						},
						"exp": schema.Int64Attribute{
							MarkdownDescription: "This value is a NumericDate as defined in RFC 7519 (A JSON numeric value representing the number of seconds from 1970-01-01T00:00:00Z UTC until the specified UTC date/time, ignoring leap seconds.)",
							Computed:            true,
							Sensitive:           true,
						},
						"nbf": schema.Int64Attribute{
							MarkdownDescription: "This value is a NumericDate as defined in RFC 7519 (A JSON numeric value representing the number of seconds from 1970-01-01T00:00:00Z UTC until the specified UTC date/time, ignoring leap seconds.)",
							Computed:            true,
							Sensitive:           true,
						},
						"x5c": schema.ListAttribute{
							ElementType:         types.StringType,
							MarkdownDescription: "The x5c (X.509 certificate chain) parameter contains a chain of one or more PKIX certificates RFC 5280.",
							Computed:            true,
							Sensitive:           true,
						},
						"x5t": schema.StringAttribute{
							MarkdownDescription: "The x5t (X.509 certificate SHA-1 thumbprint) parameter is a base64url-encoded SHA-1 thumbprint (also known as digest) of the DER encoding of an X.509 certificate RFC 5280.",
							Computed:            true,
							Sensitive:           true,
						},
					},
				},
			},
		},
	}
)
