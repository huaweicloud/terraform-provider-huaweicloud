---
subcategory: "GaussDB(for MySQL)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_gaussdb_mysql_account"
description: ""
---

# huaweicloud_gaussdb_mysql_account

Manages a GaussDB MySQL account resource within HuaweiCloud.

## Example Usage

```hcl
variable "instance_id" {}
variable "password" {}

resource "huaweicloud_gaussdb_mysql_account" "test" {
  instance_id = var.instance_id
  name        = "test_account_name"
  password    = var.password
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `instance_id` - (Required, String, ForceNew) Specifies the ID of the GaussDB MySQL instance.

  Changing this parameter will create a new resource.

* `name` - (Required, String, ForceNew) Specifies the database username. The value can contain 1 to 32 characters,
  including letters, digits, and underscores (_).

  Changing this parameter will create a new resource.

* `password` - (Required, String) Specifies the password of the database user. It cannot be the same as the username.
  The value cannot be empty and must consist of 8 to 32 characters and contain at least three types of the following:
  uppercase letters, lowercase letters, digits, and special characters (~!@#$%^*-_=+?,()&).

* `host` - (Optional, String, ForceNew) Specifies the host IP address. The default value is %, indicating that all IP
  addresses are allowed to access your GaussDB(for MySQL) instance. If its value is 10.10.10.%, all 10.10.10.X IP
  addresses can access your GaussDB(for MySQL) instance.

  Changing this parameter will create a new resource.

* `description` - (Optional, String) Specifies the database user remarks. The value can consist of up to 512 characters,
  and cannot contain the carriage return characters or special characters (!<"='>&).This field is only suitable for
  instances 2.0.13.0 or later.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID which is formatted `<instance_id>/<name>/<host>`.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 10 minutes.
* `update` - Default is 10 minutes.
* `delete` - Default is 10 minutes.

## Import

The GaussDB MySQL account can be imported using the `instance_id`, `name` and `host` separated by slashes, e.g.

```bash
$ terraform import huaweicloud_gaussdb_mysql_account.test <instance_id>/<name>/<host>
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response, security or some other reason. The missing attributes include: `password`. It is generally recommended
running `terraform plan` after importing a cluster. You can then decide if changes should be applied to the GaussDB
MySQL account, or the resource definition should be updated to align with the GaussDB MySQL account. Also you can
ignore changes as below.

```hcl
resource "huaweicloud_gaussdb_mysql_account" "test" {
    ...

  lifecycle {
    ignore_changes = [
      password,
    ]
  }
}
```
