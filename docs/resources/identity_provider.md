---
subcategory: "Identity and Access Management (IAM)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_identity_provider"
description: ""
---

# huaweicloud_identity_provider

Manages the identity providers within HuaweiCloud IAM service.

-> **NOTE:** 1. You *must* have admin privileges to use this resource.
  <br/>2. You can create up to 10 identity providers.

## Example Usage

### Create a SAML protocol provider

```hcl
resource "huaweicloud_identity_provider" "provider_1" {
  name     = "example_com_provider_saml"
  protocol = "saml"
  metadata = file("/usr/local/data/files/metadata.txt")
}
```

### Create a OpenID Connect protocol provider

```hcl
resource "huaweicloud_identity_provider" "provider_2" {
  name     = "example_com_provider_oidc"
  protocol = "oidc"
  
  access_config {
    access_type            = "program_console"
    provider_url           = "https://accounts.example.com"
    client_id              = "your_client_id"
    authorization_endpoint = "https://accounts.example.com/o/oauth2/v2/auth"
    scopes                 = ["openid"]
    signing_key            = jsonencode(
    {
      keys = [
        {
          alg = "RS256"
          e   = "AQAB"
          kid = "..."
          kty = "RSA"
          n   = "..."
          use = "sig"
        },
      ]
    }
    )
  }
}
```

<!--markdownlint-disable MD033-->

## Argument Reference

The following arguments are supported:

* `name` - (Required, String, ForceNew) Specifies the name of the identity provider to be registered.
  The maximum length is `64` characters. Only letters, digits, underscores (_), and hyphens (-) are allowed.
  The name is unique, it is recommended to include domain name information.
  Changing this creates a new resource.

* `protocol` - (Optional, String, ForceNew) Specifies the protocol of the identity provider.
  Valid values are **saml** and **oidc**. Changing this creates a new resource.
  This parameter can be omitted during creation and using `huaweicloud_identity_provider_conversion` and
  `huaweicloud_identity_provider_protocol` resources to manage it. At the same time, please use
  `lifecycle.ignore_changes` to ignore changes to `conversion_rules` under this management way.

* `status` - (Optional, Bool) Enabled status for the identity provider. Defaults to true.

* `description` - (Optional, String) Specifies the description of the identity provider.

* `sso_type` - (Optional, String, ForceNew) Specifies the single sign-on type of the identity provider.
  Valid values are as follows:
  + **virtual_user_sso**: After a federated user logs in to HuaweiCloud, the system automatically creates a virtual user
    and assigns permissions to the user based on identity conversion rules.
  + **iam_user_sso**: After a federated user logs in to HuaweiCloud, the system automatically maps the external identity
    ID to an IAM user so that the federated user has the permissions of the mapped IAM user.

  The default value is **virtual_user_sso**. For details about how to choose an SSO type,
  see [Application Scenarios of Virtual User SSO and IAM User SSO](https://support.huaweicloud.com/intl/en-us/usermanual-iam/iam_08_0251.html).
  Changing this creates a new resource.

  -> **NOTE:** 1. Only one SSO type of identity provider can be created under an account.
    <br/>2. The identity provider with OIDC protocol only supports virtual user SSO type.
    <br/>3. An account can create multiple identity providers with virtual user SSO type.
    <br/>4. An account can create only one identity provider with IAM user SSO type.
    <br/>5. When you use IAM user SSO type, make sure that you have set **IAM_SAML_Attributes_xUserId** in the IDP
      and External Identity ID in the SP to the same value.

* `metadata` - (Optional, String) Specifies the metadata of the IDP(Identity Provider) server.
  To obtain the metadata file of your enterprise IDP, contact the enterprise administrator.
  This field is used to import a metadata file to IAM to implement federated identity authentication.
  This field is required only if the protocol is set to *saml*.
  The maximum length is 30,000 characters and it stores in the state with SHA1 algorithm.

  -> **NOTE:**
  The metadata file specifies API addresses and certificate information in compliance with the SAML 2.0 standard.
  It is usually stored in a file. In the TF script, you can import the metafile through the `file` function,
  for example:
  <br/>`metadata = file("/usr/local/data/files/metadata.txt")`

* `access_config` - (Optional, List) Specifies the description of the identity provider.
  This field is required only if the protocol is set to *oidc*.

The `access_config` block supports:

* `access_type` - (Required, String) Specifies the access type of the identity provider.
  Available options are:
  + **program**: programmatic access only.
  + **program_console**: programmatic access and management console access.

* `provider_url` - (Required, String) Specifies the URL of the identity provider.
  This field corresponds to the iss field in the ID token.

* `client_id` - (Required, String) Specifies the ID of a client registered with the OpenID Connect identity provider.

* `signing_key` - (Required, String) Public key used to sign the ID token of the OpenID Connect identity provider.
  This field is required only if the protocol is set to *oidc*.

* `authorization_endpoint` - (Optional, String) Specifies the authorization endpoint of the OpenID Connect identity
  provider. This field is required only if the access type is set to `program_console`.

* `scopes` - (Optional, List) Specifies the scopes of authorization requests. It is an array of one or more scopes.
  Valid values are *openid*, *email*, *profile* and other values defined by you.
  This field is required only if the access type is set to `program_console`.

  -> **NOTE:** 1. *openid* must be specified for this field.
  <br/>2. A maximum of 10 values can be specified, and they must be separated with spaces.
  <br/>Example: openid email host.

* `response_type` - (Optional, String) Response type. Valid values is *id_token*, default value is *id_token*.
  This field is required only if the access type is set to `program_console`.

* `response_mode` - (Optional, String) Response mode.
  Valid values is *form_post* and *fragment*, default value is *form_post*.
  This field is required only if the access type is set to `program_console`.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - A resource ID which equals the identity provider name.

* `login_link` - The login link of the identity provider.

* `conversion_rules` - The identity conversion rules of the identity provider.
  The [object](#conversion_rules) structure is documented below

<a name="conversion_rules"></a>
The `conversion_rules` block supports:

* `local` - The federated user information on the cloud platform.

* `remote` - The description of the identity provider.

The `local` block supports:

* `username` - The name of a federated user on the cloud platform.

* `group` - The user group to which the federated user belongs on the cloud platform.

The `remote` block supports:

* `attribute` - The attribute in the IDP assertion.

* `condition` - The condition of conversion rule.

* `value` - The rule is matched only if the specified strings appear in the attribute type.

## Import

Identity provider can be imported using the `name`, e.g.

```bash
$ terraform import huaweicloud_identity_provider.provider_1 example_com_provider_saml
```
