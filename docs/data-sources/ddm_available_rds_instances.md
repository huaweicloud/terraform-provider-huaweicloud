---
subcategory: "Distributed Database Middleware (DDM)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_ddm_available_rds_instances"
description: |-
  Use this data source to obtain a list of RDS instances available for DDM schema creation.
---

# huaweicloud_ddm_available_rds_instances

Use this data source to obtain a list of RDS instances available for DDM schema creation.

## Example Usage

```hcl
variable "instance_id" {}

data "huaweicloud_ddm_available_rds_instances" "test" {
  instance_id = var.instance_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the ID of the DDM instance.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `instances` - Indicates the list of available RDS instances.
  The [Instance](#DdmInstances_Instance) structure is documented below.

<a name="DdmInstances_Instance"></a>
The `Instances` block supports:

* `id` - Indicates the DB instance ID.

* `project_id` - Indicates the project ID of the tenant whom the DB instance belongs to in a region.

* `status` - Indicates the DB instance status.

* `name` - Indicates the DB instance name.

* `engine_name` - Indicates the engine name of the DB instance.

* `engine_software_version` - Indicates the engine version of the DB instance.

* `private_ip` - Indicates the private IP address for connecting to the DB instance.

* `mode` - Indicates the DB instance type (primary/standby or single-node).

* `port` - Indicates the port for connecting to the DB instance.

* `az_code` - Indicates the AZ.

* `time_zone` - Indicates the time zone.
