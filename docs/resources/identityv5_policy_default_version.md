---
subcategory: "Identity and Access Management (IAM)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_identityv5_policy_default_version"
description: |-
  Use this resource to set the default version of an identity policy within HuaweiCloud.
---

# huaweicloud_identityv5_policy_default_version

Use this resource to set the default version of an identity policy within HuaweiCloud.

-> This resource is a one-time action resource used to set the default version of an identity policy. Deleting this
   resource will not clear the corresponding request record, but will only remove the resource information from the
   tfstate file.

## Example Usage

```hcl
variable "policy_id" {}
variable "version_id" {}

resource "huaweicloud_identityv5_policy_default_version" "test" {
  policy_id  = var.policy_id
  version_id = var.version_id
}
```

## Argument Reference

The following arguments are supported:

* `policy_id` - (Required, String, NonUpdatable) Specifies the ID of the identity policy.

* `version_id` - (Required, String, NonUpdatable) Specifies the version ID of the identity policy to be set
  as default.  
  The value must be a string starting with `v` followed by a number, e.g. `v1`.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.
