---
subcategory: "IAM Identity Center"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_identitycenter_account_provisioning_permission_sets"
description: |-
  Use this data source to get the Identity Center account provisioning permission sets.
---

# huaweicloud_identitycenter_account_provisioning_permission_sets

Use this data source to get the Identity Center account provisioning permission sets.

## Example Usage

```hcl
variable "instance_id" {}
variable "account_id" {}

data "huaweicloud_identitycenter_account_provisioning_permission_sets" "test" {
  instance_id = var.instance_id
  account_id  = var.account_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the ID of an IAM Identity Center instance.

* `account_id` - (Required, String) Specifies the ID of a specified account.

* `provisioning_status` - (Optional, String) Specifies the authorization status of a permission set.
  The valid values are as follows:
  + **LATEST_PERMISSION_SET_PROVISIONED**
  + **LATEST_PERMISSION_SET_NOT_PROVISIONED**

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `permission_sets` - The permission set ID list.
