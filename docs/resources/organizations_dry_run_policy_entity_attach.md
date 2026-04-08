---
subcategory: "Organizations"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_organizations_dry_run_policy_entity_attach"
description: |-
  Manages an Organizations dry-run policy attach to entity resource within HuaweiCloud.
---

# huaweicloud_organizations_dry_run_policy_entity_attach

Manages an Organizations dry-run policy attach to entity resource within HuaweiCloud.

-> Before using this resource, you must ensure that the policy dry-run feature is enabled.

## Example Usage

```hcl
variable "policy_id" {}
variable "entity_id" {}

resource "huaweicloud_organizations_dry_run_policy_entity_attach" "test" {
  policy_id = var.policy_id
  entity_id = var.entity_id
}
```

## Argument Reference

The following arguments are supported:

* `policy_id` - (Required, String, NonUpdatable) Specifies the ID of the dry-run policy.

* `entity_id` - (Required, String, NonUpdatable) Specifies the ID of the entity (root, OU, or account).

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID, which is formatted as `<policy_id>/<entity_id>`.

* `entity_name` - The name of the entity.

* `entity_type` - The type of the entity.

## Import

The dry-run policy entity attach can be imported using the `policy_id` and `entity_id`, separated by a slash, e.g.

```bash
$ terraform import huaweicloud_organizations_dry_run_policy_entity_attach.test <policy_id>/<entity_id>
```
