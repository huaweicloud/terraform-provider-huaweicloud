---
subcategory: "GeminiDB"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_gaussdb_influx_instances"
description: |-
  Use this data source to get the list of GaussDB influx instances.
---

# huaweicloud_gaussdb_influx_instances

Use this data source to get the list of GaussDB influx instances.

## Example Usage

```hcl
data "huaweicloud_gaussdb_influx_instances" "this" {}
```

## Argument Reference

* `region` - (Optional, String) The region in which to obtain the instances. If omitted, the provider-level region will
  be used.

* `instance_id` - (Optional, String) Specifies the ID of the instance. If you enter an instance ID starting with
  an asterisk (*), fuzzy search results are returned. If you enter a valid instance ID, an exact result is returned.

* `name` - (Optional, String) Specifies the name of the instance. If you enter an instance name starting with an
  asterisk (*), fuzzy search results are returned. If you enter a valid instance name, an exact result is returned.

* `mode` - (Optional, String) Specifies the instance type. Value options:
  + **Cluster**: indicating that the instance is a GeminiDB Influx instance.
  + **InfluxdbSingle**: indicating that the instance is a single-node GeminiDB Influx instance.

* `vpc_id` - (Optional, String) Specifies the VPC ID.

* `subnet_id` - (Optional, String) Specifies the network ID of a subnet.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `instances` - Indicates the list of GaussDB influx instances.
  The [instances](#instances_struct) structure is documented below.

<a name="instances_struct"></a>
The `instances` block supports:

* `id` - Indicates the ID of the instance.

* `name` - Indicates the name of the instance.

* `status` - Indicates the DB instance status. The value can be:
  + **normal**: indicating that the instance is running normally.
  + **abnormal**: indicating that the instance is abnormal.
  + **creating**: indicating that the instance is being created.
  + **frozen**: indicating that the instance is frozen.
  + **data_disk_full**: indicating that the instance disk is full.
  + **createfail**: indicating that the instance failed to be created.
  + **enlargefail**: indicating that nodes failed to be added to the instance.

* `port` - Indicates the database port.

* `mode` - Indicates the instance type.

* `region` - Indicates the region where the instance is deployed.

* `datastore` - Indicates the database information.
  The [datastore](#datastore_struct) structure is documented below.

* `engine` - Indicates the storage engine. The value is **rocksDB**.

* `db_user_name` - Indicates the default username. The value is **rwuser**.

* `vpc_id` - Indicates the VPC ID.

* `subnet_id` - Indicates the network ID of a subnet.

* `security_group_id` - Indicates the security group ID.

* `backup_strategy` - Indicates the backup policy.
  The [backup_strategy](#backup_strategy_struct) structure is documented below.

* `pay_mode` - Indicates the billing mode. The value can be:
  + **0**: indicates the instance is billed on a pay-per-use basis.
  + **1**: indicates the instance is billed on a yearly/monthly basis.

* `maintain_begin` - Indicates the start time for a maintenance window.

* `maintain_end` - Indicates the end time for a maintenance window.

* `groups` - Indicates the group information.
  The [groups](#groups_struct) structure is documented below.

* `enterprise_project_id` - Indicates the enterprise project ID.

* `time_zone` - Indicates the time zone.

* `actions` - Indicates the operation that is executed on the instance.

* `dedicated_resource_id` - Indicates the dedicated resource ID. This parameter is returned only when the instance belongs
  to a dedicated resource pool.

* `lb_ip_address` - Indicates the IP address bound to the load balancer. This parameter is returned only when an IP
  address is specified for the load balancer.

* `lb_port` - Indicates the load balancing port number. This parameter is returned only when there is a load balancer address.

* `availability_zone` - Indicates the availability zone.

* `created_at` - Indicates the instance creation time.

* `updated_at` - Indicates the time when an instance is updated.

<a name="datastore_struct"></a>
The `datastore` block supports:

* `engine` - Indicates the database engine.

* `version` - Indicates the database version.

* `patch_available` - Indicates the whether there is an available patch for upgrade. If **true** is returned, you can
  install a patch to upgrade the instance.

<a name="backup_strategy_struct"></a>
The `backup_strategy` block supports:

* `start_time` - Indicates the backup time window.

* `keep_days` - Indicates the number of days to retain the generated.

<a name="groups_struct"></a>
The `groups` block contains:

* `id` - Indicates the group ID.

* `status` - Indicates the group status. The value can be:
  + **normal**: indicating that the group is normal.
  + **abnormal**: indicating that the group is abnormal.
  + **creating**: indicating that the group is being created.
  + **createfail**: indicating that the group failed to be created.
  + **deleted**: indicating that the group has been deleted.
  + **resizefailed**: indicating that the group specifications failed to be changed.
  + **enlargefail**: indicating the group failed to be scaled out.

* `volume` - Indicates the volume information.
  The [volume](#volume_struct) structure is documented below.

* `nodes` - Indicates the node information.
  The [nodes](#nodes_struct) structure is documented below.

<a name="volume_struct"></a>
The `volume` block contains:

* `size` - Indicates the storage space, in GB.

* `used` - Indicates the used storage space, in GB.

<a name="nodes_struct"></a>
The `nodes` block contains:

* `id` - Indicates the node ID.

* `name` - Indicates the node name.

* `status` - Indicates the node status. The value can be:
  + **normal**: indicating that the node is normal.
  + **abnormal**: indicating that the node is abnormal.
  + **creating**: indicating that the node is being created.
  + **createfail**: indicating that the node failed to be created.
  + **deleted**: indicating that the node has been deleted.
  + **resizefailed**: indicating that the node specifications failed to be changed.
  + **enlargefail**:  indicating nodes failed to be added.

* `subnet_id` - Indicates the ID of the subnet where the instance node is deployed.

* `private_ip` - Indicates the Private IP address of the node.

* `public_ip` - Indicates the bound EIP. This parameter is valid only for nodes bound with EIPs.

* `spec_code` - Indicates the resource specification code.

* `availability_zone` - Indicates the availability zone.

* `support_reduce` - Indicates whether instance nodes can be deleted. The value can be:
  + **true**: indicating that instance nodes can be deleted.
  + **false**: indicating that instance nodes cannot be deleted.
