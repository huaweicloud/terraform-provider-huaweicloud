---
subcategory: "Document Database Service (DDS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dds_instances"
description: |-
  Use this data source to get the list of DDS instances.
---

# huaweicloud_dds_instances

Use this data source to get the list of DDS instances.

## Example Usage

### Query all instances

```hcl
data "huaweicloud_dds_instances" "test" {}
```

### Filter by name

```hcl
variable "instance_name" {}

data "huaweicloud_dds_instances" "test" {
  name = var.instance_name
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `name` - (Optional, String) Specifies the DB instance name.
  If you use asterisk (*) at the beginning of the name, fuzzy search results are returned. Otherwise,
  the exact results are returned.

* `mode` - (Optional, String) Specifies the instance type.
  The valid values are as follows:
  + **Sharding**: Indicates the cluster instance.
  + **ReplicaSet**: Indicates the replica set instance.
  + **Single**: Indicates the single node instance.

* `vpc_id` - (Optional, String) Specifies the VPC ID.

* `subnet_id` - (Optional, String) Specifies the subnet Network ID.

* `instance_id` - (Optional, String) Specifies the instance ID.

* `datastore_type` - (Optional, String) Specifies the DB type. The value is **DDS-Community**.

* `tags` - (Optional, String) Specifies query based on the instance tag key and value.
  A maximum of `20` key-value pairs are supported. The key cannot be empty or duplicate, but the value can be empty.
  If multiple tag key-value pairs are used for the query at the same time, they should be separated by commas(,),
  indicating that the query includes instances that have the specified tag key-value pairs.
  e.g. **key1=value1,key2=value2**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `instances` - Indicates the list of DDS instances.
  The [instances](#instance_struct) structure is documented below.

<a name="instance_struct"></a>
The `instances` block supports:

* `id` - Indicates the ID of the instance.

* `name` - Indicates the DB instance name.

* `ssl` - Indicates whether to enable or disable SSL.

* `port` - Indicates the database port number. The port range is `2,100` to `9,500`.

* `datastore` - Indicates database information.
  The [datastore](#datastore_struct) structure is documented below.

* `backup_strategy` - Indicates the backup policy information.
  The [backup_strategy](#backup_strategy_struct) structure is documented below.

* `vpc_id` - Indicates the VPC ID

* `subnet_id` - Indicates the subnet Network ID.

* `security_group_id` - Indicates the security group ID of the DDS instance.

* `disk_encryption_id` - Indicates the disk encryption ID of the instance.

* `mode` - Indicates the instance type.

* `db_username` - Indicates the DB default user name.

* `status` - Indicates the DB instance status.

* `enterprise_project_id` - Indicates the enterprise project ID.

* `groups` - Indicates the instance groups information.
  The [groups](#groups_struct) structure is documented below.

* `tags` - Indicates the key/value pairs to associate with the DDS instance.

* `remark` - Indicates the instance description.

* `engine` - Indicates the storage engine.

* `created` - Indicates the instance creation time.

* `updated` - Indicates the instance update time.

* `pay_mode` - Indicates the billing mode.

* `maintenance_window` - Indicates the maintenance time window.

* `time_zone` - Indicates the time zone.

* `dss_pool_id` - Indicates the DSS storage pool ID of the Dec user.

* `actions` - Indicates the action that is being executed on an instance.

<a name="datastore_struct"></a>
The `datastore` block supports:

* `type` - Indicates the DB engine.

* `version` - Indicates the DB instance version.

* `whole_version` - Indicates the DB instance complete version number.

* `patch_available` - Whether there is an available patch for upgrade.

<a name="backup_strategy_struct"></a>
The `backup_strategy` block supports:

* `start_time` - Indicates the backup time window.

* `keep_days` - Indicates the number of days to retain the generated backup files.

<a name="groups_struct"></a>
The `groups` block supports:

* `type` - Indicates the node type.

* `id` - Indicates the group ID.

* `name` - Indicates the group name.

* `status` - Indicates the group status.

* `size` - Indicates the disk size, in GB.

* `used` - Indicates the disk usage, in GB.

* `nodes` - Indicates the nodes info.
  The [nodes](#node_struct) structure is documented below.

<a name="node_struct"></a>
The `nodes` block supports:

* `id` - Indicates the node ID.

* `name` - Indicates the node name.

* `role` - Indicates the node role.

* `private_ip` - Indicates the private IP address of a node.

* `public_ip` - Indicates the EIP that has been bound on a node.

* `status` - Indicates the node status.

* `spec_code` - Indicates the node spec code.

* `availability_zone` - Indicates the availability zone.
