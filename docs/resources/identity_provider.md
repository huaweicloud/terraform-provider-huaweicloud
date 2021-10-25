---
subcategory: "Identity and Access Management (IAM)"
---

# huaweicloud_identity_provider

Manages the identity providers within HuaweiCloud IAM service.

-> **NOTE:** You can create up to 10 identity providers.

## Example Usage

### Create a SAML protocol provider

```hcl
resource "huaweicloud_identity_provider" "provider_1" {
  name     = "example_com_provider_saml"
  protocol = "saml"

  conversion_rules {
    local {
      username = "FederationUser"
    }
    remote {
      attribute = "username"
      condition = "any_one_of"
      value     = ["Tom", "Jerry"]
    }
  }
}
```

### Create a OpenID Connect protocol provider

```hcl
resource "huaweicloud_identity_provider" "provider_2" {
  name     = "example_com_provider_oidc"
  protocol = "oidc"

  conversion_rules {
    local {
      username = "FederationUser"
    }
    remote {
      attribute = "username"
      condition = "any_one_of"
      value     = ["Tom", "Jerry"]
    }
  }

  access_type            = "program_console"
  provider_url           = "https://accounts.example.com"
  client_id              = "client_id_example"
  authorization_endpoint = "https://accounts.example.com/o/oauth2/v2/auth"
  scopes                 = ["openid", "email"]
  signing_key            = jsonencode(
  {
    keys = [
      {
        alg = "RS256"
        e   = "AQAB"
        kid = "d05ef20c4512645vv1..."
        kty = "RSA"
        n   = "cws_cnjiwsbvweolwn_-vnl..."
        use = "sig"
      },
    ]
  }
  )
}
```

<!--markdownlint-disable MD033-->

## Argument Reference

The following arguments are supported:

* `name` - (Required, String, ForceNew) Specifies the name of the identity provider to be registered.
  The maximum length is 64 characters. Only letters, digits, underscores (_), and hyphens (-) are allowed.
  The name is unique, it is recommended to include domain name information.
  Changing this creates a new resource.

* `protocol` - (Required, String, ForceNew) Specifies the protocol of the identity provider.
  Valid values are *saml* and *oidc*.
  Changing this creates a new resource.

* `conversion_rules` - (Required, String) Specifies the identity conversion rules of the identity provider.
  You can use identity conversion rules to map the identities of existing users to Huawei Cloud and control their access
  to cloud resources.
  The [object](#conversion_rules) structure is documented below.

* `status` - (Optional, String) Specifies the status of the identity provider.
  Valid values are *enabled* and *disabled*. Default value is *enabled*.

* `sso_type` - (Optional, String, ForceNew) Specifies the type of the identity provider.
  Valid values are *virtual_user_sso* and *iam_user_sso*. Default value is *virtual_user_sso*.

* `access_type` - (Optional, String) Specifies the access type of the identity provider.
  Available options are:
  + `program`: programmatic access only.
  + `program_console`: programmatic access and management console access.
  This field is required only if the protocol is set to *oidc*.

* `provider_url` - (Optional, String) Specifies the URL of the identity provider.
  This field corresponds to the iss field in the ID token.
  This field is required only if the protocol is set to *oidc*.

* `client_id` - (Optional, String) Specifies the ID of a client registered with the OpenID Connect identity provider.
  This field is required only if the protocol is set to *oidc*.

* `signing_key` - (Optional, String) Public key used to sign the ID token of the OpenID Connect identity provider.
  This field is required only if the protocol is set to *oidc*.

* `authorization_endpoint` - (Optional, String) Specifies the authorization endpoint of the OpenID Connect identity
  provider. This field is required only if the protocol is set to *oidc* and the access type is set to programmatic
  access and management console access.

* `scopes` - (Optional, List) Specifies the scopes of authorization requests. It is an array of one or more scopes.
  Valid values are *openid*, *email*, *profile* and other values defined by you.
  This field is required only if the protocol is set to *oidc* and the access type is set to programmatic access and
  management console access.

-> **NOTE:** 1. *openid* must be specified for this field.
<br/>2. A maximum of 10 values can be specified, and they must be separated with spaces.
<br/>Example: openid email host.

* `response_type` - (Optional, String) Response type. Only support *id_token* currently.
  This field is required only if the protocol is set to *oidc* and the access type is set to programmatic
  access and management console access.

* `response_mode` - (Optional, String) Response mode.
  Valid values are *form_post* and *fragment*, default value is *form_post*.
  This field is required only if the protocol is set to *oidc* and the access type is set to programmatic
  access and management console access.

* `metadata` - (Optional, String) Specifies the metadata of the IdP(Identity Provider) server.
  To obtain the metadata file of your enterprise IdP, contact the enterprise administrator.
  This field is used to import a metadata file to IAM to implement federated identity authentication.
  This field is required only if the protocol is set to *saml*.

-> **NOTE:**
The metadata file specifies API addresses and certificate information in compliance with the SAML 2.0 standard.
It is usually stored in a file. In the TF script, you can import the metafile through the `file` function,
for example:
<br/>`metadata = file("/usr/local/data/files/metadata.txt")`

* `description` - (Optional, String) Specifies the description of the identity provider.

<a name="conversion_rules"></a>
The `conversion_rules` block supports:

* `local` - (Required, List) Specifies the federated user information on the cloud platform.

* `remote` - (Required, List) Specifies the description of the identity provider.

The `local` block supports:

* `username` - (Required, String) Specifies the name of a federated user on the cloud platform.

* `group` - (Optional, String) Specifies the user group to which the federated user belongs on the cloud platform.

The `remote` block supports:

* `attribute` - (Required, String) Specifies the attribute in the IdP assertion.

* `condition` - (Optional, String) Specifies the condition of conversion rule.
  Available options are:
  + `any_one_of`: The rule is matched only if the specified strings appear in the attribute type.
  + `not_any_of`: The rule is matched only if the specified strings do not appear in the attribute type.

-> **NOTE:** 1. The condition result is Boolean rather than the argument that is passed as input.
  <br/>2. In a remote array, `any_one_of` and `not_any_of` are mutually exclusive and cannot be set at the same time.

* `value` - (Optional, List) Specifies the rule is matched only if the specified strings appear in the attribute type.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - A resource ID in UUID format.

* `login_link` - The login link of the identity provider.

## Import

Identity provider can be imported using the `name`, e.g.

```
$ terraform import huaweicloud_identity_provider.provider_1 example_com_provider_saml
```
