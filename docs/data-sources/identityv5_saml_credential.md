---
subcategory: "IAM"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_identityv5_saml_credential"
description: |-
  Use this data source to obtain the temporary security credentials by assuming a trust agency with a SAML authentication response.
---

# huaweicloud_identityv5_saml_credential

Use this data source to obtain the temporary security credentials by assuming a trust agency with a SAML
authentication response (API: `POST /v5/agencies/assume-with-saml`).

Before using this data source, you need to:
1. Create a SAML identity provider within IAM v5.
2. Create a trust agency whose trust principal is the SAML identity provider.

-> The temporary security credentials are refreshed every time `terraform plan` or `terraform apply` runs, and they
will expire according to the `expiration` attribute. Please use them in the scenarios where the credentials are
consumed immediately.

## Example Usage

```hcl
variable "saml_response" {}

data "huaweicloud_identityv5_saml_credential" "test" {
  saml_response       = var.saml_response
  agency_urn          = "iam::${var.account_id}:agency:demo"
  agency_session_name = "terraform-session"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `saml_response` - (Required, String, Sensitive) Specifies the SAML authentication response (assertion) issued by the
  SAML identity provider, which is a base64 encoded XML string.

* `agency_urn` - (Required, String) Specifies the uniform resource name of the trust agency to be assumed, in the
  format of `iam::{account_id}:agency:{agency_name}`.

* `agency_session_name` - (Required, String) Specifies the name of the agency session. The name contains 2 to 64
  characters and can only contain letters, digits, underscores (_) and hyphens (-).

* `provider_urn` - (Optional, String) Specifies the uniform resource name of the SAML identity provider, in the
  format of `iam::{account_id}:samlProvider:{provider_name}`.

* `duration_seconds` - (Optional, Int) Specifies the validity period of the temporary security credentials, in seconds.
  The value ranges from 900 to 43200 and defaults to 3600.

* `policy` - (Optional, String, Sensitive) Specifies a custom policy in JSON format, which is used to further restrict
  the permissions of the temporary security credentials.

* `policy_ids` - (Optional, List) Specifies the list of preset policy IDs, which is used to restrict the permissions
  of the temporary security credentials. A maximum of 20 policy IDs can be configured.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `access_key_id` - The access key (AK) of the temporary security credentials.

* `secret_access_key` - The secret access key (SK) of the temporary security credentials.

* `security_token` - The security token of the temporary security credentials, which is used as the request header
  `x-security-token`.

* `expiration` - The expiration time (UTC) of the temporary security credentials.

* `assumed_agency_urn` - The uniform resource name of the assumed agency session.

* `assumed_agency_id` - The ID of the assumed agency session.

* `subject` - The value of the NameID element of the Subject element in the SAML assertion.

* `subject_type` - The format of the NameID, defined by the Format attribute of the NameID element in the SAML assertion.

* `issuer` - The value of the Issuer element in the SAML assertion.

* `audience` - The value of the Recipient attribute of the SubjectConfirmationData element in the SAML assertion.

* `name_qualifier` - The hash value of the issuer, the account ID and the name of the SAML identity provider.

* `idp_urn` - The uniform resource name of the SAML identity provider.

* `source_identity` - The identity declared by the SourceIdentity attribute in the SAML assertion.
