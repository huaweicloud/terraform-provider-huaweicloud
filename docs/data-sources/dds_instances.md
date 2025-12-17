---
subcategory: "Document Database Service (DDS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dds_instances"
description: ""
---

# huaweicloud_dds_instances

Use this data source to get the list of DDS instances.

## Example Usage

```hcl
variable "vpc_id" {}
variable "subnet_id" {}

data "huaweicloud_dds_instances" "test" {
  name      = "test_name"
  mode      = "Sharding"
  vpc_id    = var.vpc_id
  subnet_id = var.subnet_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `name` - (Optional, String) Specifies the DB instance name.

* `mode` - (Optional, String) Specifies the mode of the database instance.

* `vpc_id` - (Optional, String) Specifies the VPC ID.

* `subnet_id` - (Optional, String) Specifies the subnet Network ID.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `instances` - Indicates the list of DDS instances.
  The [Instance](#DdsInstance_Instance) structure is documented below.

<a name="DdsInstance_Instance"></a>
The `Instance` block supports:

* `id` - Indicates the ID of the instance.

* `name` - Indicates the DB instance name.

* `ssl` - Indicates whether to enable or disable SSL.

* `port` - Indicates the database port number. The port range is 2100 to 9500.

* `datastore` - Indicates database information.
  The [Datastore](#DdsInstance_InstanceDatastore) structure is documented below.

* `backup_strategy` - Indicates the database information.
  The [BackupStrategy](#DdsInstance_InstanceBackupStrategy) structure is documented below.

* `vpc_id` - Indicates the VPC ID

* `subnet_id` - Indicates the subnet Network ID.

* `security_group_id` - Indicates the security group ID of the DDS instance.

* `disk_encryption_id` - Indicates the disk encryption ID of the instance.

* `mode` - Specifies the mode of the database instance.

* `db_username` - Indicates the DB Administrator name.

* `status` - Indicates the the DB instance status.

* `enterprise_project_id` - Indicates the enterprise project id of the dds instance.

* `groups` - Indicates the instance groups information.
  The [group](#DdsInstance_InstanceGroup) structure is documented below.

* `tags` - Indicates the key/value pairs to associate with the DDS instance.

<a name="DdsInstance_InstanceDatastore"></a>
The `InstanceDatastore` block supports:

* `type` - Indicates the DB engine.

* `version` - Indicates the DB instance version.

* `storage_engine` - Indicates the storage engine of the DB instance.

<a name="DdsInstance_InstanceBackupStrategy"></a>
The `InstanceBackupStrategy` block supports:

* `start_time` - Indicates the backup time window.

* `keep_days` - Indicates the number of days to retain the generated backup files.

<a name="DdsInstance_InstanceGroup"></a>
The `group` block supports:

* `type` - Indicates the node type.

* `id` - Indicates the group ID.

* `name` - Indicates the group name.

* `status` - Indicates the group status.

* `size` - Indicates the disk size.

* `used` - Indicates the disk usage.

* `nodes` - Indicates the nodes info.
  The [node](#DdsInstance_InstanceGroupNode) structure is documented below.

<a name="DdsInstance_InstanceGroupNode"></a>
The `node` block supports:

* `id` - Indicates the node ID.

* `name` - Indicates the node name.

* `role` - Indicates the node role.

* `private_ip` - Indicates the private IP address of a node.

* `public_ip` - Indicates the EIP that has been bound on a node.

* `status` - Indicates the node status.

* `spec_code` - Indicates the node spec code.

* `availability_zone` - Indicates the availability zone.
