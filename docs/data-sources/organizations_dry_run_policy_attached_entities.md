---
subcategory: "Organizations"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_organizations_dry_run_policy_attached_entities"
description: |-
  Use this data source to get entitiy list attached to the specified dry run policy within HuaweiCloud.
---

# huaweicloud_organizations_dry_run_policy_attached_entities

Use this data source to get the entity list attached to the specified dry run policy within HuaweiCloud.

## Example Usage

```hcl
variable "policy_id" {}

data "huaweicloud_organizations_dry_run_policy_attached_entities" "test" {
  policy_id = var.policy_id
}
```

## Argument Reference

The following arguments are supported:

* `policy_id` - (Required, String) Specifies the ID of the dry run policy.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `entities` - The entities that are attached to the specified dry run policy.  
  The [entities](#dry_run_policy_attached_entities) structure is documented below.

<a name="dry_run_policy_attached_entities"></a>
The `entities` block supports:

* `id` - The unique ID of the entity.

* `type` - The type of the entity.

* `name` - The name of the entity.
