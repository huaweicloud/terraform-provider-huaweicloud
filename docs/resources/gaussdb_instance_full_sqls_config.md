---
subcategory: "GaussDB"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_gaussdb_instance_full_sqls_config"
description: |-
  Manages a GaussDB instance full SQL collection configuration resource within HuaweiCloud.
---

# huaweicloud_gaussdb_instance_full_sqls_config

Manages a GaussDB instance full SQL collection configuration resource within HuaweiCloud.

## Example Usage

```hcl
variable "instance_id" {}
variable "storage_mode" {}
variable "is_exclude_sys_user" {}
variable "save_days" {}
variable "lts_config" {}
variable "sql_type_range" {}

resource "huaweicloud_gaussdb_instance_full_sqls_config" "test" {
  instance_id         = var.instance_id
  storage_mode        = var.storage_mode
  save_days           = var.save_days
  is_exclude_sys_user = var.is_exclude_sys_user
  lts_config {
    log_group_name  = "GROUP_GAUSSDB_APS_1"
    log_stream_name = "STREAM_APS_FULL_SQL_1"
  }
  sql_type_range {
    category = "all"
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `instance_id` - (Required, String, NonUpdatable) Specifies the ID of the GaussDB instance.

* `storage_mode` - (Required, String, NonUpdatable) Specifies the storage mode for full SQL data.
  Valid value is **LTS**.

* `save_days` - (Required, Int) Specifies the number of days to retain full SQL data.
  The value ranges from `1` to `30`.

* `is_exclude_sys_user` - (Optional, Bool) Specifies whether to exclude system users from SQL collection.

* `lts_config` - (Required, List) Specifies the LTS (Log Tank Service) configuration.
  The [lts_config](#gaussdb_full_sqls_lts_config) structure is documented below.

<a name="gaussdb_full_sqls_lts_config"></a>
The `lts_config` block supports:

* `log_group_name` - (Required, String) Specifies the name of the LTS log group.

* `log_stream_name` - (Required, String) Specifies the name of the LTS log stream.

* `sql_type_range` - (Optional, List) Specifies the SQL type range configuration.
  The [sql_type_range](#gaussdb_full_sqls_sql_type_range) structure is documented below.

<a name="gaussdb_full_sqls_sql_type_range"></a>
The `sql_type_range` block supports:

* `category` - (Required, String) Specifies the SQL category.
  Valid values include **all**, **ddl**, **dml**, **dcl**, **tcl**, **dql**, **custom**.

* `prefixes` - (Optional, List) Specifies the SQL statement prefixes to match.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID which is the same as `instance_id`.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 20 minutes.
* `update` - Default is 20 minutes.
* `delete` - Default is 20 minutes.
