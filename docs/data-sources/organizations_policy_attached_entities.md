---
subcategory: "Organizations"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_organizations_policy_attached_entities"
description: |-
  Use this data source to get the list of the entities that the specified policy is attached to.
---

# huaweicloud_organizations_policy_attached_entities

Use this data source to get the list of the entities that the specified policy is attached to.

## Example Usage

```hcl
variable "policy_id "{}

data "huaweicloud_organizations_policy_attached_entities" "test" {
  policy_id = var.policy_id
}
```

## Argument Reference

The following arguments are supported:

* `policy_id` - (Required, String) Specifies the ID of the policy.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `attached_entities` - Indicates the entities that the specified policy is attached to.

  The [attached_entities](#attached_entities_struct) structure is documented below.

<a name="attached_entities_struct"></a>
The `attached_entities` block supports:

* `id` - Indicates the ID of the entity.

* `type` - Indicates the type of the entity.

* `name` - Indicates the name of the entity.
