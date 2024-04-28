---
subcategory: "GaussDB(for MySQL)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_gaussdb_mysql_database"
description: ""
---

# huaweicloud_gaussdb_mysql_database

Manages a GaussDB MySQL database resource within HuaweiCloud.

## Example Usage

```hcl
variable "instance_id" {}
variable "character_set" {}

resource "huaweicloud_gaussdb_mysql_database" "test" {
  instance_id   = var.instance_id
  name          = "test_db_name"
  character_set = var.character_set
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `instance_id` - (Required, String, ForceNew) Specifies the ID of the GaussDB MySQL instance.

  Changing this parameter will create a new resource.

* `name` - (Required, String, ForceNew) Specifies the database name. The name can contain 1 to 64 characters.
  Only letters, digits, hyphens (-), and underscores (_) are allowed. The total number of hyphens (-) cannot exceed 10.

  Changing this parameter will create a new resource.

* `character_set` - (Required, String, ForceNew) Specifies the database character set.
  Value options: **utf8mb4**, **utf8**, **latin1**, **gbk**.

  Changing this parameter will create a new resource.

* `description` - (Optional, String) Specifies the database remarks. The value can consist of up to 512 characters,
  and cannot contain the carriage return characters or special characters (!<"='>&).

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID which is formatted `<instance_id>/<name>`.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 10 minutes.
* `update` - Default is 10 minutes.
* `delete` - Default is 10 minutes.

## Import

The GaussDB MySQL database can be imported using the `instance_id` and `name` separated by a slash, e.g.

```bash
$ terraform import huaweicloud_gaussdb_mysql_database.test <instance_id>/<name>
```
