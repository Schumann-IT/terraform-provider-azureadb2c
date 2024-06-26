---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "azureadb2c_application_patch Resource - azureadb2c"
subcategory: ""
description: |-
  Manages special requirements for application registration within Azure AD B2C when
  using custom policies https://learn.microsoft.com/en-us/azure/active-directory-b2c/user-flow-overview#custom-policies.
  Please refer to the following examples:
  Identity Experience Framework applications https://learn.microsoft.com/en-us/azure/active-directory-b2c/tutorial-create-user-flows?pivots=b2c-custom-policy#register-identity-experience-framework-applicationsSAML applications https://learn.microsoft.com/en-us/azure/active-directory-b2c/saml-service-provider?tabs=windows&pivots=b2c-custom-policyDaemon applications https://learn.microsoft.com/en-us/azure/active-directory-b2c/client-credentials-grant-flow?pivots=b2c-custom-policy
  Other applications (like web and native apps) can still be configured via Azure Active Directory Provider https://registry.terraform.io/providers/hashicorp/azuread/latest/docs/resources/application
---

# azureadb2c_application_patch (Resource)

Manages special requirements for application registration within Azure AD B2C when 
using [custom policies](https://learn.microsoft.com/en-us/azure/active-directory-b2c/user-flow-overview#custom-policies).

Please refer to the following examples:
- [Identity Experience Framework applications](https://learn.microsoft.com/en-us/azure/active-directory-b2c/tutorial-create-user-flows?pivots=b2c-custom-policy#register-identity-experience-framework-applications)
- [SAML applications](https://learn.microsoft.com/en-us/azure/active-directory-b2c/saml-service-provider?tabs=windows&pivots=b2c-custom-policy) 
- [Daemon applications](https://learn.microsoft.com/en-us/azure/active-directory-b2c/client-credentials-grant-flow?pivots=b2c-custom-policy) 

Other applications (like web and native apps) can still be configured via [Azure Active Directory Provider](https://registry.terraform.io/providers/hashicorp/azuread/latest/docs/resources/application)

## Example Usage

```terraform
data "azuread_application" "example" {
  display_name = "My First AzureAD Application"
}

resource "azureadb2c_application_patch" "example" {
  object_id  = azuread_application.example.object_id
  patch_file = format("%s/path/to/patch.json", path.module)
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `object_id` (String) The object id of the application to be patched
- `patch_file` (String) The path to the patch file. Must be an absolute path to a JSON file

### Optional

- `saml_metadata_url` (String) The SAML metadata url

### Read-Only

- `data` (Attributes) identity experience framework app data (see [below for nested schema](#nestedatt--data))

<a id="nestedatt--data"></a>
### Nested Schema for `data`

Read-Only:

- `app_id` (String) The application id (client id)
- `display_name` (String) The display name
- `id` (String) The id of the application
- `identifier_uris` (List of String) The identifier uris
- `saml_metadata_url` (String) The saml metadata url
