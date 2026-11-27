---
subcategory: "IAM"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_identityv5_oidc_credential"
description: |-
  Use this data source to obtain the temporary security credentials by assuming a trust agency with an OIDC ID token.
---

# huaweicloud_identityv5_oidc_credential

Use this data source to obtain the temporary security credentials by assuming a trust agency with an OIDC ID token
(API: `POST /v5/agencies/assume-with-oidc`).

-> Before using this data source, you need to:
  <br/>1. Create an OIDC identity provider within IAM v5.
  <br/>2. Create a trust agency whose trust principal is the OIDC identity provider.

-> The temporary security credentials are refreshed every time `terraform plan` or `terraform apply` runs, and they
will expire according to the `credentials.0.expiration` attribute. Please use them in the scenarios where the
credentials are consumed immediately.

## Example Usage

```hcl
variable "id_token" {}

data "huaweicloud_identityv5_oidc_credential" "test" {
  id_token            = var.id_token
  agency_urn          = "iam::${var.account_id}:agency:demo"
  agency_session_name = "terraform-session"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `id_token` - (Required, String, Sensitive) Specifies the OpenID Connect ID token issued by the OIDC identity provider.

* `agency_urn` - (Required, String) Specifies the uniform resource name of the trust agency to be assumed, in the
  format of `iam::{account_id}:agency:{agency_name}`.

* `agency_session_name` - (Required, String) Specifies the name of the agency session. The name contains 2 to 64
  characters and can only contain letters, digits, underscores (_) and hyphens (-).

* `provider_urn` - (Optional, String) Specifies the uniform resource name of the OIDC identity provider, in the
  format of `iam::{account_id}:oidcProvider:{provider_name}`.

* `duration_seconds` - (Optional, Int) Specifies the validity period of the temporary security credentials, in seconds.
  The value ranges from 900 to 43200 and defaults to 3600.

* `policy` - (Optional, String, Sensitive) Specifies a custom policy in JSON format, which is used to further restrict
  the permissions of the temporary security credentials.

* `policy_ids` - (Optional, List) Specifies the list of preset policy IDs, which is used to restrict the permissions
  of the temporary security credentials. A maximum of 20 policy IDs can be configured.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `assumed_agency` - The information about the assumed agency session.
  The [assumed_agency](#assumed_agency) structure is documented below.

* `credentials` - The temporary security credentials.
  The [credentials](#credentials) structure is documented below.

* `subject_from_id_token` - The unique user identifier returned from the `sub` field of the ID token.

* `audience` - The client ID returned from the `aud` or `azp` field of the ID token.

* `idp_urn` - The uniform resource name of the OIDC identity provider.

* `source_identity` - The source identity declared by the source identity field of the ID token.

<a name="assumed_agency"></a>
The `assumed_agency` block supports:

* `urn` - The uniform resource name of the assumed agency session.

* `id` - The ID of the assumed agency session.

<a name="credentials"></a>
The `credentials` block supports:

* `access_key_id` - The access key (AK) of the temporary security credentials.

* `secret_access_key` - The secret access key (SK) of the temporary security credentials.

* `security_token` - The security token of the temporary security credentials, which is used as the request header
  `x-security-token`.

* `expiration` - The expiration time (UTC) of the temporary security credentials.
