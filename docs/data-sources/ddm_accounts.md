---
subcategory: "Distributed Database Middleware (DDM)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_ddm_accounts"
description: ""
---

# huaweicloud_ddm_accounts

Use this data source to get the list of DDM accounts.

## Example Usage

```hcl
variable "ddm_instance_id" {}
variable "account_name" {}

data "huaweicloud_ddm_accounts" "test" {
  instance_id = var.ddm_instance_id
  name        = var.account_name
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the ID of DDM instance.

* `name` - (Optional, String) Specifies the name of the DDM account.

* `status` - (Optional, String) Specifies the status of the DDM account.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `accounts` - Indicates the list of DDM account.
  The [Account](#DdmAccounts_Account) structure is documented below.

<a name="DdmAccounts_Account"></a>
The `Account` block supports:

* `name` - Indicates the name of the DDM account.

* `status` - Indicates the status of the DDM account.

* `permissions` - Indicates the basic permissions of the DDM account.

* `description` - Indicates the description of the DDM account.

* `schemas` - Indicates the schemas that associated with the account.
  The [Schema](#DdmAccounts_AccountSchema) structure is documented below.

<a name="DdmAccounts_AccountSchema"></a>
The `AccountSchema` block supports:

* `name` - Indicates the name of the associated schema.

* `description` - Indicates the schema description.
