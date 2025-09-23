---
subcategory: "Organizations"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_organizations_effective_policies"
description: |-
  Use this data source to get the effective policies of a specific type for the specified account.
---

# huaweicloud_organizations_effective_policies

Use this data source to get the effective policies of a specific type for the specified account.

## Example Usage

```hcl
variable "entity_id" {}

data "huaweicloud_organizations_effective_policies" "test"{
  entity_id   = var.entity_id
  policy_type = "tag_policy"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `entity_id` - (Required, String) Specifies the unique ID of an account.
  Currently, the effective policy of the root and organizational units cannot be queried.

* `policy_type` - (Required, String) Specifies the name of a policy type.
  Currently, the value **tag_policy** is available.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `last_updated_at` - Indicates the time when the effective policy is mostly updated.

* `policy_content` - Indicates the content of the effective policy.
