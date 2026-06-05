---
subcategory: "Data Admin Service (DAS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_das_instances"
description: |-
  Use this data source to query the DAS instances under a specified region within HuaweiCloud.
---

# huaweicloud_das_instances

Use this data source to query the DAS instances under a specified region within HuaweiCloud.

## Example Usage

### Query MySQL instances

```hcl
data "huaweicloud_das_instances" "test" {
  datastore_type = "MySQL"
}
```

### Query SQLServer instances

```hcl
data "huaweicloud_das_instances" "test" {
  datastore_type = "SQLServer"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region where the instances are located.  
  If omitted, the provider-level region will be used.

* `datastore_type` - (Required, String) Specifies the type of the database.  
  The valid values are as follows:
  + **MySQL**: Cloud Database RDS for MySQL.
  + **SQLServer**: Cloud Database RDS for SQL Server.
  + **PostgreSQL**: Cloud Database RDS for PostgreSQL.
  + **TaurusDB**: Cloud Database TaurusDB.
  + **gaussdbv5**: Cloud Database GaussDB.
  + **mongodb**: Document Database Service DDS.
  + **DDM**: Distributed Database Middleware DDM.
  + **MariaDB**: Cloud Database RDS for MariaDB.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `instances` - The list of instances that matched filter parameters.  
  The [instances](#das_instances_attr) structure is documented below.

<a name="das_instances_attr"></a>
The `instances` block supports:

* `id` - The ID of the instance.

* `name` - The name of the instance.

* `status` - The status of the instance.

* `version` - The version of the instance.

* `engine_type` - The engine type of the instance.

* `ip` - The IP address of the instance.

* `port` - The port of the instance.

* `cpu` - The CPU cores of the instance.

* `mem` - The memory size of the instance, in GB.

* `login_flag` - Whether login is enabled.

* `slow_sql_flag` - Whether slow SQL analysis is enabled.

* `deadlock_flag` - Whether deadlock analysis is enabled.

* `lock_blocking_flag` - Whether lock blocking analysis is enabled.

* `charge_flag` - Whether the instance is charged.

* `full_sql_flag` - Whether full SQL is enabled.
