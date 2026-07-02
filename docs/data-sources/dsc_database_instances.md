---
subcategory: "DSC"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dsc_database_instances"
description: |-
  Use this data source to get the list of database instances within HuaweiCloud.
---

# huaweicloud_dsc_database_instances

Use this data source to get the list of database instances within HuaweiCloud.

## Example Usage

```hcl
variable "instance_type" {}

data "huaweicloud_dsc_database_instances" "test" {
  instance_type = var.instance_type
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `instance_type` - (Required, String) Specifies the database instance type.
  Valid values: **RDS**, **DDS**, **GaussDB**, **GaussDB(for MySQL)**, **GaussDB(for Mongo)**,
  **GaussDB(for Redis)**, **GaussDB(for Cassandra)**, **GaussDB(for Influx)**, **GaussDB(for ClickHouse)**,
  **GaussDB(for OpenGauss)**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `instances` - The basic information list of the instances.

  The [instances](#instances_struct) structure is documented below.

<a name="instances_struct"></a>
The `instances` block supports:

* `id` - The unique identifier of the instance.

* `ins_id` - The instance ID.

* `ins_name` - The instance name.

* `ins_status` - The instance status.

* `ins_type` - The instance type.

* `db_type` - The database type.

* `version` - The instance version.

* `address` - The network address of the instance.

* `port` - The network port of the instance.

* `bind_database` - The number of bound databases.

* `is_external` - Whether the instance is external.

* `project_id` - The project ID.

* `vpc_id` - The VPC ID.

* `subnet_id` - The subnet ID.

* `security_group_id` - The security group ID.

* `create_time` - The creation time of the instance.
