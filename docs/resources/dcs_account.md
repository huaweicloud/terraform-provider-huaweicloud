---
subcategory: "Distributed Cache Service (DCS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dcs_account"
description: ""
---

# huaweicloud_dcs_account

Manages a DCS account resource within HuaweiCloud.

## Example Usage

```hcl
variable "instance_id" {}

resource "huaweicloud_dcs_account" "test" {
  instance_id      = var.instance_id
  account_name     = "user"
  account_role     = "read"
  account_password = "Terraform@123"
  description      = "add account"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this creates a new resource.

* `instance_id` - (Required, String, ForceNew) Specifies the ID of the DCS instance.
  Changing this creates a new resource.

* `account_name` - (Required, String, ForceNew) Specifies the name of the account.
  Changing this creates a new resource.

* `account_password` - (Required, String) Specifies the password of the account.

* `account_role` - (Required, String) Specifies the role of the account.
  Value options:
  + **read**: The account has read-only privilege.
  + **write**: The account has read and write privilege.

* `description` - (Optional, String) Specifies the description of the account.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID. The value is the ID of the account.

* `account_type` - Indicates the type of the account. The value can be **normal** or **default**.

* `status` - Indicates the status of the account. The value can be **CREATING**, **AVAILABLE**, **CREATEFAILED**,
  **DELETED**, **DELETEFAILED**, **DELETING**, **UPDATING** or **ERROR**.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 30 minutes.
* `update` - Default is 30 minutes.
* `delete` - Default is 30 minutes.

## Import

The DCS account can be imported using the DCS instance ID and the DCS account ID separated by a slash, e.g.

```bash
$ terraform import huaweicloud_dcs_account.test <instance_id>/<id>
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response, security or some other reason. The missing attributes include: `account_password`.
It is generally recommended running `terraform plan` after importing the account.
You can then decide if changes should be applied to the account, or the resource definition should be updated to
align with the account. Also you can ignore changes as below.

```hcl
resource "huaweicloud_dcs_account" "test" {
    ...

  lifecycle {
    ignore_changes = [
      account_password,
    ]
  }
}
```
