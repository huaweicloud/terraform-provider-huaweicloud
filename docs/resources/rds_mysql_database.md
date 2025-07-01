---
subcategory: "Relational Database Service (RDS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_rds_mysql_database"
description: ""
---

# huaweicloud_rds_mysql_database

Manages RDS Mysql database resource within HuaweiCloud.

## Example Usage

```hcl
variable "instance_id" {}

resource "huaweicloud_rds_mysql_database" "test" {
  instance_id   = var.instance_id
  name          = "test"
  character_set = "utf8"
  description   = "test database"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) The region in which to create the RDS database resource. If omitted, the
  provider-level region will be used. Changing this creates a new resource.

* `instance_id` - (Required, String, NonUpdatable) Specifies the RDS instance ID.

* `name` - (Required, String, NonUpdatable) Specifies the database name. The database name contains `1` to `64`
  characters. The name can only consist of lowercase letters, digits, hyphens (-), underscores (_) and dollar signs
  ($). The total number of hyphens (-) and dollar signs ($) cannot exceed `10`. RDS for **MySQL 8.0** does not
  support dollar signs ($).

* `character_set` - (Required, String, NonUpdatable) Specifies the character set used by the database, For example **utf8**,
  **gbk**, **ascii**, etc.

* `description` - (Optional, String) Specifies the database description. The value can contain `0` to `512` characters.
  This parameter takes effect only for DB instances whose kernel versions are at least **5.6.51.3**, **5.7.33.1**,
  or **8.0.21.4**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID of database which is formatted `<instance_id>/<name>`.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 30 minutes.
* `update` - Default is 30 minutes.
* `delete` - Default is 30 minutes.

## Import

RDS database can be imported using the `instance id` and `name` separated by slash, e.g.

```bash
$ terraform import huaweicloud_rds_mysql_database.database_1 <instance_id>/<name>
```
