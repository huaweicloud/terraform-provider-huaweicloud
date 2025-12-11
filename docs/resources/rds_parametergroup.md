---
subcategory: "Relational Database Service (RDS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_rds_parametergroup"
description: |-
  Manages an RDS parameter group resource within HuaweiCloud.
---

# huaweicloud_rds_parametergroup

Manages an RDS parameter group resource within HuaweiCloud.

## Example Usage

```hcl
resource "huaweicloud_rds_parametergroup" "test" {
  name        = "test_name"
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

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource. If omitted, the
  provider-level region will be used. Changing this creates a new resource.

* `name` - (Required, String) Specifies the parameter group name. It contains a maximum of 64 characters.

* `description` - (Optional, String) Specifies the parameter group description. It contains a maximum of 256 characters
  and cannot contain the following special characters:>!<"&'= the value is left blank by default.

* `values` - (Optional, Map) Specifies the parameter group values key/value pairs defined by users based on the default
  parameter groups.

* `datastore` - (Required, List, NonUpdatable) Specifies the database object.
  The [datastore](#datastore_struct) structure is documented below.

<a name="datastore_struct"></a>
The `datastore` block supports:

* `type` - (Required, String, NonUpdatable) Specifies the DB engine. Currently, MySQL, PostgreSQL, Microsoft SQL Server
  and MariaDB are supported. The value is case-insensitive and can be **mysql**, **postgresql**, **sqlserver**,
  or **mariadb**.

* `version` - (Required, String, NonUpdatable) Specifies the database version.
  + MySQL databases support MySQL 5.6 and 5.7. Example value: 5.7.
  + PostgreSQL databases support PostgreSQL 9.5 and 9.6. Example value: 9.5.
  + Microsoft SQL Server databases support 2014 SE, 2016 SE, and 2016 EE. Example value: 2014_SE.
  + MariaDB databases support MariaDB 10.5. Example value: 10.5.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - Indicates the resource ID.

* `configuration_parameters` - Indicates the parameter configuration defined by users based on the default parameters
  groups. The [configuration_parameters](#configuration_parameters_struct) structure is documented below.

* `created_at` - Indicates the creation time, in UTC format.

* `updated_at` - Indicates the last update time, in UTC format.

<a name="configuration_parameters_struct"></a>
The `configuration_parameters` block supports:

* `name` - Indicates the parameter name.

* `value` - Indicates the parameter value.

* `restart_required` - Indicates whether a restart is required.

* `readonly` - Indicates whether the parameter is read-only.

* `value_range` - Indicates the parameter value range.

* `type` - Indicates the parameter type.

* `description` - Indicates the parameter description.

## Import

The RDS parameter group can be imported using the `id`, e.g.

```bash
$ terraform import huaweicloud_rds_parametergroup.test <id>
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response, security or some other reason. The missing attributes include: `values`. It is generally recommended
running `terraform plan` after importing the RDS parameter group. You can then decide if changes should be applied to
the RDS parameter group, or the resource definition should be updated to align with the RDS parameter group. Also you
can ignore changes as below.

```hcl
resource "huaweicloud_rds_parametergroup" "test" {
    ...

  lifecycle {
    ignore_changes = [
      values
    ]
  }
}
```
