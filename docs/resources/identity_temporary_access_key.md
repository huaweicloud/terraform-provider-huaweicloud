---
subcategory: "Identity and Access Management (IAM)"
layout: "huaweicloud"
page_title: "HuaweiCloud: identity_temporary_access_key"
description: |-
    Use this data source to obtain temporary access keys and security tokens within HuaweiCloud.
---

# huaweicloud_identity_temporary_access_key

Use this data source to obtain temporary access keys and security tokens within HuaweiCloud.

## Example Usage

### Obtain temporary access keys and security tokens by agency

```hcl
variable "token" {}
variable "agency_name" {}
variable "domain_name" {}

resource "huaweicloud_identity_temporary_access_key" "test" {
  token       = var.token
  methods     = "assume_role"
  agency_name = var.agency_name
  domain_name = var.domain_name
}
```

### Obtain temporary access keys and security tokens by token

```hcl
variable "token" {}

resource "huaweicloud_identity_temporary_access_key" "test" {
  token   = var.token
  methods = "token"
}
```

## Argument Reference

* `token` - (Required, String, ForceNew) Specifies the token.

* `methods` - (Required, String, ForceNew) Specifies the authentication method,
  the content of this field is either `token` or `assume_role`.

* `policy` - (Optional, String, ForceNew) Specifies the user-defined policy information. It is used to restrict the
  permissions of the obtained temporary access key and security token (currently only supported by the OBS service).

* `agency_name` - (Optional, String, ForceNew) Specifies the agency name. When the `methods` is `assume_role`, it is
  mandatory.

* `domain_id` - (Optional, String, ForceNew) Specifies the account(domain) id of the agency.
  When the `methods` is `assume_role`, `domain_id` or `domain_name` are mandatory.

* `domain_name` - (Optional, String, ForceNew) Specifies the account(domain) name of the agency.
  When the `methods` is `assume_role`, `domain_id` or `domain_name` are mandatory.

* `duration_seconds` - (Optional, Int, ForceNew) Specifies the validity period of AK/SK and security token.
  It's measured in seconds. The range is from 15 minutes to 24 hours, with a default value of 15 minutes.

* `session_user_name` - (Optional, String, ForceNew) Specifies the enterprise user information corresponding to the
  agency.
  When the `methods` is `assume_role`, it will be effective.

## Attribute Reference

* `expires_at` - Indicates the expiration time of AK/SK and security token. It is in UTC time format.

* `access` - Indicates the access key (AK).

* `secret` - Indicates the secret key (SK).

* `securitytoken` - Indicates an encrypted string containing the obtained AK, SK and other security information.
