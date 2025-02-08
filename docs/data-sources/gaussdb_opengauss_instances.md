---
subcategory: "GaussDB"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_gaussdb_opengauss_instances"
description: |-
  Use this data source to get available GaussDB OpenGauss instances.
---

# huaweicloud_gaussdb_opengauss_instances

Use this data source to get available GaussDB OpenGauss instances.

## Example Usage

```hcl
data "huaweicloud_gaussdb_opengauss_instances" "this" {
  name = "gaussdb-instance"
}
```

## Argument Reference

* `region` - (Optional, String) The region in which to obtain the instance. If omitted, the provider-level region will
  be used.

* `name` - (Optional, String) Specifies the name of the instance.

* `vpc_id` - (Optional, String) Specifies the VPC ID.

* `subnet_id` - (Optional, String) Specifies the network ID of a subnet.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - Indicates the ID of the data source.

* `instances` - An array of available instances.

The `instances` block supports:

* `region` - The region of the instance.

* `id` - Indicates the id of the instance.

* `name` - Indicates the name of the instance.

* `vpc_id` - Indicates the VPC ID.

* `subnet_id` - Indicates the network ID of a subnet.

* `status` - Indicates the DB instance status.

* `type` - Indicates the instance type.

* `flavor` - Indicates the instance specifications.

* `security_group_id` - Indicates the security group ID.

* `enterprise_project_id` - Indicates the enterprise project id.

* `db_user_name` - Indicates the default username.

* `time_zone` - Indicates the default username.

* `availability_zone` - Indicates the instance availability zone.

* `port` - Indicates the database port.

* `switch_strategy` - Indicates the switch strategy.

* `maintenance_window` - Indicates the maintenance window.

* `coordinator_num` - Indicates the count of coordinator node.

* `sharding_num` - Indicates the sharding num.

* `replica_num` - Indicates the replica num.

* `private_ips` - Indicates the list of private IP address of the nodes.

* `public_ips` - Indicates the public IP address of the DB instance.

* `volume` - Indicates the volume information. Structure is documented below.

* `datastore` - Indicates the database information. Structure is documented below.

* `backup_strategy` - Indicates the advanced backup policy. Structure is documented below.

* `nodes` - Indicates the instance nodes information. Structure is documented below.

* `ha` - Indicates the instance ha information. Structure is documented below.

* `mysql_compatibility_port` - Indicates the port for MySQL compatibility.

The `volume` block supports:

* `type` - Indicates the volume type.
* `size` - Indicates the volume size.

The `datastore` block supports:

* `engine` - Indicates the database engine.
* `version` - Indicates the database version.

The `backup_strategy` block supports:

* `start_time` - Indicates the backup time window.
* `keep_days` - Indicates the number of days to retain the generated

The `nodes` block contains:

* `id` - Indicates the node ID.
* `name` - Indicates the node name.
* `status` - Indicates the node status.
* `role` - Indicates whether the node support reduce.
* `availability_zone` - Indicates the availability zone where the node resides.

The `ha` block supports:

* `replication_mode` - Indicates the replication mode.
