---
subcategory: "Organizations"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_organizations_policy_attach"
description: ""
---

# huaweicloud_organizations_policy_attach

Manages an Organizations policy attach resource within HuaweiCloud.

## Example Usage

```hcl
variable policy_id {}
variable entity_id {}

resource "huaweicloud_organizations_policy_attach" "test"{
  policy_id = var.policy_id
  entity_id = var.entity_id
}
```

## Argument Reference

The following arguments are supported:

* `policy_id` - (Required, String, ForceNew) Specifies the ID of the policy.

  Changing this parameter will create a new resource.

* `entity_id` - (Required, String, ForceNew) Specifies the unique ID of the root, OU, or account.

  Changing this parameter will create a new resource.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID which is formatted `<policy_id>/<entity_id>`.

* `entity_name` - Indicates the name of the entity.

* `entity_type` - Indicates the type of the entity.

## Import

The organizations policy attach can be imported using the `policy_id` and `entity_id` separated by a slash, e.g.:

```bash
$ terraform import huaweicloud_organizations_policy_attach.test <policy_id>/<entity_id>
```
