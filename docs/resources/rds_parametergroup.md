---
subcategory: "Relational Database Service (RDS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_rds_parametergroup"
description: ""
---

# huaweicloud_rds_parametergroup

Manages a RDS ParameterGroup resource within HuaweiCloud.

## Example Usage

```hcl
resource "huaweicloud_rds_parametergroup" "pg_1" {
  name        = "pg_1"
  description = "description_1"
  values      = {
    max_connections = "10"
    autocommit      = "OFF"
  }
  datastore {
    type    = "mysql"
    version = "5.6"
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) The region in which to create the RDS parameter group. If omitted, the
  provider-level region will be used. Changing this creates a new parameter group.

* `name` - (Required, String) The parameter group name. It contains a maximum of 64 characters.

* `description` - (Optional, String) The parameter group description. It contains a maximum of 256 characters and cannot
  contain the following special characters:>!<"&'= the value is left blank by default.

* `values` - (Optional, Map) Parameter group values key/value pairs defined by users based on the default parameter
  groups.

* `datastore` - (Required, List, ForceNew) Database object. The database object structure is documented below. Changing
  this creates a new parameter group.

The `datastore` block supports:

* `type` - (Required, String) The DB engine. Currently, MySQL, PostgreSQL, Microsoft SQL Server and MariaDB are supported.
  The value is case-insensitive and can be **mysql**, **postgresql**, **sqlserver**, or **mariadb**.

* `version` - (Required, String) Specifies the database version.

  + MySQL databases support MySQL 5.6 and 5.7. Example value: 5.7.
  + PostgreSQL databases support PostgreSQL 9.5 and 9.6. Example value: 9.5.
  + Microsoft SQL Server databases support 2014 SE, 2016 SE, and 2016 EE. Example value: 2014_SE.
  + MariaDB databases support MariaDB 10.5. Example value: 10.5.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the parameter group.

* `configuration_parameters` - Indicates the parameter configuration defined by users based on the default parameters
  groups.

  + `name` - Indicates the parameter name.
  + `value` - Indicates the parameter value.
  + `restart_required` - Indicates whether a restart is required.
  + `readonly` - Indicates whether the parameter is read-only.
  + `value_range` - Indicates the parameter value range.
  + `type` - Indicates the parameter type.
  + `description` - Indicates the parameter description.

* `created_at` - The creation time, in UTC format.

* `updated_at` - The last update time, in UTC format.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 10 minutes.
* `delete` - Default is 10 minutes.

## Import

Parameter groups can be imported using the `id`, e.g.

```bash
$ terraform import huaweicloud_rds_parametergroup.pg_1 7117d38e-4c8f-4624-a505-bd96b97d024c
```
