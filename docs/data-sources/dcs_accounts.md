---
subcategory: "Distributed Cache Service (DCS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dcs_accounts"
description: |-
  Use this data source to get the list of DCS accounts.
---

# huaweicloud_dcs_accounts

Use this data source to get the list of DCS accounts.

## Example Usage

```hcl
variable "instance_id" {}

data "huaweicloud_dcs_accounts" "test"{
  instance_id = var.instance_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the instance ID.

* `account_name` - (Optional, String) Specifies the account name.

* `account_type` - (Optional, String) Specifies the account type. The value can be **normal** or **default**.

* `account_role` - (Optional, String) Specifies the account role. The value can be **read** or **write**.

* `description` - (Optional, String) Specifies the account description.

* `status` - (Optional, String) Specifies the account status.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `accounts` - ACL account list.

  The [accounts](#accounts_struct) structure is documented below.

<a name="accounts_struct"></a>
The `accounts` block supports:

* `id` - Account ID.

* `account_name` - Account name.

* `account_type` - Account type.

* `account_role` - Account permissions.

* `description` - Account description.

* `status` - Account status.
  The value can be:
  + **CREATING**: The account is being created.
  + **AVAILABLE**: The account is available.
  + **CREATEFAILED**: The account fails to be created.
  + **DELETED**: The account has been deleted.
  + **DELETEFAILED**: The account fails to be deleted.
  + **DELETING**: The account is being deleted.
  + **UPDATING**: The account is being updated.
  + **ERROR**: The account is abnormal.
