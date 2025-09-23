---
subcategory: "IAM Identity Center"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_identitycenter_system_identity_policy_attachments"
description: |-
  Use this data source to get the Identity Center system identity policy attachments.
---

# huaweicloud_identitycenter_system_identity_policy_attachments

Use this data source to get the Identity Center system identity policy attachments.

## Example Usage

```hcl
variable "instance_id" {}
variable "permission_set_id" {}

data "huaweicloud_identitycenter_system_identity_policy_attachments" "test" {
  instance_id       = var.instance_id
  permission_set_id = var.permission_set_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the ID of an IAM Identity Center instance.

* `permission_set_id` - (Required, String) Specifies the ID of a permission set.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `policies` - The list of IAM system-defined identity policies.

  The [policies](#policies_struct) structure is documented below.

<a name="policies_struct"></a>
The `policies` block supports:

* `id` - The ID of the IAM system-defined identity policy.

* `name` - The name of the IAM system-defined identity policy.
