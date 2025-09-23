---
subcategory: "Distributed Database Middleware (DDM)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_ddm_account"
description: ""
---

# huaweicloud_ddm_account

Manages a DDM account resource within HuaweiCloud.

## Example Usage

```hcl
variable "ddm_instance_id" {}
variable "name" {}
variable "password" {}
variable "schema_name" {}

resource "huaweicloud_ddm_account" "test"{
  instance_id = var.ddm_instance_id
  name        = var.name
  password    = var.password

  permissions = [
    "CREATE",
    "SELECT"
  ]

  schemas {
   name = var.schema_name
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `instance_id` - (Required, String, ForceNew) Specifies the ID of a DDM instance.
  Changing this parameter will create a new resource.

* `name` - (Required, String, ForceNew) Specifies the name of the DDM account.
  An account name starts with a letter, consists of 1 to 32 characters, and can contain only letters,
  digits, and underscores (_).
  Changing this parameter will create a new resource.

* `password` - (Required, String) Specifies the DDM account password.

* `permissions` - (Required, List) Specifies the basic permissions of the DDM account. Value options: **CREATE**
  **DROP**、**ALTER**、**INDEX**、**INSERT**、**DELETE**、**UPDATE**、**SELECT**.

* `description` - (Optional, String) Specifies the description of the DDM account.

* `schemas` - (Optional, List) Specifies the schemas that associated with the account.
  The [Schema](#DdmAccount_Schema) structure is documented below.

<a name="DdmAccount_Schema"></a>
The `Schema` block supports:

* `name` - (Optional, String) Specifies the name of the associated schema.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `status` - Indicates the status of the DDM account.

* `schemas` - (Optional, List) Specifies the schemas that associated with the account.
  The [Schema](#DdmAccount_Schema) structure is documented below.

<a name="DdmAccount_Schema"></a>
The `Schema` block supports:

* `description` - (Optional, String) Specifies the schema description.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 30 minutes.
* `update` - Default is 10 minutes.
* `delete` - Default is 10 minutes.

## Import

The DDM account can be imported using the instance ID and account name separated by a slash, e.g.:

```bash
$ terraform import huaweicloud_ddm_account.test 0a8f1c6baa124e99853719d9257324dfin09/account_name
```
