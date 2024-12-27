---
subcategory: "IAM Identity Center"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_identitycenter_permission_set_provisioning_accounts"
description: |-
  Use this data source to get the Identity Center permission set provisioning accounts.
---

# huaweicloud_identitycenter_permission_set_provisioning_accounts

Use this data source to get the Identity Center permission set provisioning accounts.

## Example Usage

```hcl
variable "instance_id" {}
variable "permission_set_id" {}

data "huaweicloud_identitycenter_permission_set_provisioning_accounts" "test" {
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

* `provisioning_status` - (Optional, String) Specifies the provisioning status of a permission set.
  The valid values are as follows:
  + **LATEST_PERMISSION_SET_PROVISIONED**
  + **LATEST_PERMISSION_SET_NOT_PROVISIONED**

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `account_ids` - The account ID list.
