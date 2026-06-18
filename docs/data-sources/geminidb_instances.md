---
subcategory: "GeminiDB"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_geminidb_instances"
description: |-
  Use this data source to get the list of GeminiDB instances.
---

# huaweicloud_geminidb_instances

Use this data source to get the list of GeminiDB instances.

## Example Usage

### Query all instances

```hcl
data "huaweicloud_geminidb_instances" "test" {}
```

### Query instance by instance ID

```hcl
variable "instance_id" {}

data "huaweicloud_geminidb_instances" "test" {
  instance_id = var.instance_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `instance_id` - (Optional, String) Specifies the instance ID.

* `name` - (Optional, String) Specifies the instance name.
  If the name starts with `*`, fuzzy search results are returned. Otherwise, an exact result is returned.

* `datastore_type` - (Optional, String) Specifies the database type.
  The valid values are as follows:
  + **cassandra**: Indicates GeminiDB Cassandra instances are queried.
  + **mongodb**: Indicates GeminiDB Mongo instances are queried.
  + **influxdb**: Indicates GeminiDB Influx instances are queried.
  + **redis**: Indicates GeminiDB Redis instances are queried.

* `mode` - (Optional, String) Specifies the instance type.
  The valid values are as follows:
  + **Cluster**:Indicates GeminiDB Cassandra, GeminiDB Influx, or proxy cluster GeminiDB Redis instance
  with classic storage.
  + **CloudNativeCluster**: Indicates GeminiDB Cassandra, Influx, or Redis cluster instance with cloud native storage.
  + **RedisCluster**: Indicates Redis Cluster GeminiDB Redis instance with classic storage.
  + **Replication**: Indicates primary/standby GeminiDB Redis instance with classic storage.
  + **InfluxdbSingle**: Indicates single-node GeminiDB Influx instance with classic storage.
  + **EnhancedCluster**: Indicates GeminiDB Influx cluster (performance-enhanced) instance with classic storage.
  + **ReplicaSet**: Indicates GeminiDB Mongo instance in a replica set.

  The parameter is invalid if `datastore_type` is not specified.

* `vpc_id` - (Optional, String) Specifies the VPC ID.

* `subnet_id` - (Optional, String) Specifies the subnet ID.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `instances` - The list of instances.
  The [instances](#instances_struct) structure is documented below.

<a name="instances_struct"></a>
The `instances` block supports:

* `id` - The instance ID.

* `name` - The instance name.

* `status` - The instance status.
  + **normal**
  + **abnormal**
  + **creating**
  + **frozen**
  + **data_disk_full**
  + **createfail**
  + **enlargefail**

* `port` - The database port.

* `mode` - The instance type.

* `product_type` - The product type used for GeminiDB Redis instances with cloud native storage.
  + **Capacity**
  + **Performance**

* `region` - The region where the instance is deployed.

* `datastore` - The database information.
  The [datastore](#datastore_struct) structure is documented below.

* `engine` - The storage engine.
  + **rocksDB**

* `created` - The creation time.

* `updated` - The update time.

* `db_user_name` - The default user name.

* `vpc_id` - The VPC ID.

* `subnet_id` - The subnet ID.

* `security_group_id` - The security group ID.

* `backup_strategy` - The backup policy information.
  The [backup_strategy](#backup_strategy_struct) structure is documented below.

* `pay_mode` - The billing mode.
  + **0**: The instance is billed on a pay-per-use basis.
  + **1**: The instance is billed on a yearly/monthly basis.

* `maintenance_window` - The maintenance time window.

* `groups` - The instance group information.
  The [groups](#groups_struct) structure is documented below.

* `enterprise_project_id` - The enterprise project ID.

* `time_zone` - The time zone.

* `actions` - The operation that is executed on the instance.

* `dedicated_resource_id` - The dedicated resource ID.

* `disk_encryption_id` - The key ID used for disk encryption.

* `lb_ip_address` - The load balancer IP address.

* `lb_port` - The load balancer port.

* `availability_zone` - The AZ information.

* `dr_instance_id` - The DR instance ID.

* `dual_active_info` - The active-active instance information.
  The [dual_active_info](#dual_active_info_struct) structure is documented below.

* `ccm_cert_info` - The CCM certificate information.
  The [ccm_cert_info](#ccm_cert_info_struct) structure is documented below.

<a name="datastore_struct"></a>
The `datastore` block supports:

* `type` - The database type.

* `version` - The database version.

* `patch_available` - Whether the current instance can be patche.

* `whole_version` - The whole version of a GeminiDB Cassandra or Redis instance.

<a name="backup_strategy_struct"></a>
The `backup_strategy` block supports:

* `start_time` - The backup time window.

* `keep_days` - The backup retention days.

<a name="groups_struct"></a>
The `groups` block supports:

* `id` - The group ID.

* `status` - The group status.

* `volume` - The volume information.
  The [volume](#volume_struct) structure is documented below.

* `nodes` - The nodes information.
  The [nodes](#nodes_struct) structure is documented below.

<a name="volume_struct"></a>
The `volume` block supports:

* `size` - The storage size, in GB.

* `used` - The used storage in GB.

<a name="nodes_struct"></a>
The `nodes` block supports:

* `id` - The node ID.

* `name` - The node name.

* `status` - The node status.
  + **normal**
  + **abnormal**
  + **creating**
  + **createfail**
  + **deleted**
  + **resizefailed**
  + **enlargefail**

* `role` - The node role.

* `subnet_id` - The subnet ID.

* `private_ip` - The private IP address.

* `public_ip` - The EIP.

* `spec_code` - The resource specification code.

* `availability_zone` - The AZ information.

* `support_reduce` - Whether instance nodes can be deleted.

<a name="dual_active_info_struct"></a>
The `dual_active_info` block supports:

* `role` - The active-active role.

* `status` - The active-active status.
  + **normal**
  + **abnormal**

* `destination_instance_id` - The ID of the peer instance in the active-active pair.

* `destination_region` - The peer region of the active-active pair.

* `destination_instance_name` - The name of the peer instance in the active-active pair.

* `destination_instance_node_num` - The number of nodes in the peer instance in the active-active pair.

* `destination_instance_spec_code` - The specifications of the peer instance in the active-active pair.

<a name="ccm_cert_info_struct"></a>
The `ccm_cert_info` block supports:

* `cert_id` - The certificate ID.

* `cert_type` - The certificate type.
  + **PCA**: CCM PCA certificate.
  + **SSL**: CCM SSL certificate.
