---
subcategory: "Identity and Access Management (IAM)"
---

# huaweicloud_identity_provider_oidc_config

Manages the configuration information of the identity provider that uses the open connection protocol within HuaweiCloud.

## Example Usage

```hcl
variable provider_id {}
resource "huaweicloud_identity_provider_oidc_config" "config" {
  provider_id            = var.provider_id
  access_type            = "program_console"
  provider_url           = "https://accounts.example.com"
  client_id              = "client_id_example"
  authorization_endpoint = "https://accounts.example.com/o/oauth2/v2/auth"
  scopes                 = ["openid"]
  signing_key            = jsonencode(
  {
    keys = [
      {
        alg = "RS256"
        e   = "AQAB"
        kid = "your_kid"
        kty = "RSA"
        n   = "the_value_of_n"
        use = "sig"
      }
    ]
  }
  )
}
```

<!--markdownlint-disable MD033-->

## Argument Reference

The following arguments are supported:

* `provider_id` - (Required, String, ForceNew) Specifies the identity provider ID.
  Changing this creates a new resource.

* `access_type` - (Required, String) Specifies the access type of the identity provider.
  Available options are:
  + `program`: programmatic access only.
  + `program_console`: programmatic access and management console access.

* `provider_url` - (Required, String) Specifies the URL of the identity provider.
  This field corresponds to the iss field in the ID token.

* `client_id` - (Required, String) Specifies the ID of a client registered with the identity provider.

* `signing_key` - (Required, String) The public key in json string format, used to sign the ID token of the identity
  provider.

* `authorization_endpoint` - (Optional, String) Specifies the authorization endpoint of the identity provider.
  This field is required only if the access type is set to `program_console`.

* `scopes` - (Optional, List) Specifies the scopes of authorization requests. It is an array of one or more scopes.
  Valid values are *openid*, *email*, *profile* and other values defined by you.
  This field is required only if the access type is set to `program_console`.

-> **NOTE:** 1. *openid* must be specified for scopes.
<br/>2. A maximum of 10 values can be specified.

* `response_type` - (Optional, String) Response type. Only support *id_token* currently.
  This field is required only if the access type is set to `program_console`.

* `response_mode` - (Optional, String) Response mode.
  Valid values are *form_post* and *fragment*, default value is *form_post*.
  This field is required only if the access type is set to `program_console`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - A resource ID in UUID format.

## Import

The configuration information of the identity provider can be imported using the `provider_id`, e.g.

```
$ terraform import huaweicloud_identity_provider_oidc_config.config example_com_provider_oidc
```
