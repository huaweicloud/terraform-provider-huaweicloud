---
subcategory: "Document Database Service (DDS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dds_instance_storage_space_update"
description: |-
  Manages a DDS instance storage space update resource within HuaweiCloud.
---

# huaweicloud_dds_instance_storage_space_update

Manages a DDS instance storage space update resource within HuaweiCloud.

## Example Usage

```hcl
variable "instance_id" {}
variable "group_id" {}

resource "huaweicloud_dds_instance_storage_space_update" "test" {
  instance_id = var.instance_id
  size        = 40
  group_id    = var.group_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this creates a new resource.

* `instance_id` - (Required, String, NonUpdatable) Specifies the ID of a DDS instance.

* `size` - (Required, Int, NonUpdatable) Specifies the disk capacity. The value must be an integer multiple of 10 and
  greater than the current storage space.
  + In a cluster instance, this parameter indicates the storage space of shard nodes. Value range: 10 GB to 5,000 GB when
    the shard node has fewer than 8 vCPUs. 10 GB to 10,000 GB when the shard node has 8 or more vCPUs.
  + In a replica set instance, this parameter indicates the disk capacity of the DB instance to be expanded. The value
    range is from 10 GB to 5,000 GB when the instance has fewer than 8 vCPUs. The value range is from 10 GB to 10,000 GB
    when the instance has 8 or more vCPUs.
  + In a single node instance, this parameter indicates the disk capacity of the DB instance to be expanded. The value
    range is from 10 GB to 1,000 GB.

* `group_id` - (Optional, String, NonUpdatable) Specifies the role ID.
  + This parameter is not specified for replica set instances.
  + For a cluster instance, this parameter is set to the ID of the shard group.

* `node_ids` - (Optional, List, NonUpdatable) Specifies the node IDs. This parameter is required when the disk capacity
  of the read replica of a replica set instance is expanded. Only one element can be transferred in the list.

## Attribute Reference

* `id` - Indicates the resource ID. It's same as the instance ID.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 30 minutes.
