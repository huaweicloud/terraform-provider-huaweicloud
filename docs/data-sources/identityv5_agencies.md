---
subcategory: "Identity and Access Management (IAM)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_identityv5_agencies"
description: |-
  Use this data source to get a list of IAM V5 agencies or trust agencies.
---

# huaweicloud_identityv5_agencies

Use this data source to get a list of IAM V5 agencies or trust agencies.

## Example Usage

```hcl
variable "agency_id" {}

data "huaweicloud_identityv5_agencies" "all" {}

data "huaweicloud_identityv5_agencies" "test" {
  agency_id = var.agency_id
}
```

## Argument Reference

The following arguments are supported:

* `path_prefix` - (Optional, String) Specifies the resource path prefix, composed of segments of strings, each segment
  contains one or more letters, digits, `.`, `,`, `+`, `@`, `=`, `_`, or `-`, ending with `/`, for example `foo/bar/`.

* `agency_id` - (Optional, String) Specifies the id of the agency. This parameter conflicts with `path_prefix`.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `agencies` - The list of IAM agencies.
  The [agencies](#IdentityAgencies_List) structure is documented below.

<a name="IdentityAgencies_List"></a>
The `agencies` block supports:

* `trust_policy` - Indicates the trust policy of the agency.

* `agency_id` - Indicates the ID of the agency.

* `agency_name` - Indicates the name of the agency.

* `path` - Indicates the path of the agency.

* `trust_domain_id` - Indicates the account ID of the trusted domain, only exists in agencies, not in trust agencies.

* `trust_domain_name` - Indicates the account name of the trusted domain, only exists in agencies, not in trust agencies.

* `urn` - Indicates the URN of the agency.

* `created_at` - Indicates the creation time of the agency or trust agency.

* `description` - Indicates the description of the agency or trust agency.

* `max_session_duration` - Indicates the maximum session duration of the agency.
