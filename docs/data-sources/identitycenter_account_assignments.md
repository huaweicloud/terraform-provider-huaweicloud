---
subcategory: "IAM Identity Center"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_identitycenter_account_assignments"
description: |-
  Use this data source to get the Identity Center account assignments.
---

# huaweicloud_identitycenter_account_assignments

Use this data source to get the Identity Center account assignments.

## Example Usage

```hcl
variable "instance_id" {}
variable "principal_id" {}
variable "principal_type" {}

data "huaweicloud_identitycenter_account_assignments" "test" {
  instance_id    = var.instance_id
  principal_id   = var.principal_id
  principal_type = var.principal_type
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the ID of an IAM Identity Center instance.

* `principal_id` - (Required, String) Specifies the ID of the principal.

* `principal_type` - (Required, String) Specifies the type of the principal.
  The valid values are as follows:
  + **USER**
  + **GROUP**

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `account_assignments` - The list of the account assignments for principal.
  The [account_assignments](#account_assignments_struct) structure is documented below.

<a name="account_assignments_struct"></a>
The `account_assignments` block supports:

* `account_id` - The ID of the account.

* `permission_set_id` - The ID of the permission set.

* `principal_id` - The ID of the principal.

* `principal_type` - The type of the principal.
