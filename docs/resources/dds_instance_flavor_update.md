---
subcategory: "Document Database Service (DDS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dds_instance_flavor_update"
description: |-
  Manages a DDS instance flavor update resource within HuaweiCloud.
---

# huaweicloud_dds_instance_flavor_update

Manages a DDS instance flavor update resource within HuaweiCloud.

## Example Usage

### Update instance's mongos node flavor

```hcl
variable "instance_id" {}
variable "node_id" {}

resource "huaweicloud_dds_instance_flavor_update" "test" {
  instance_id      = var.instance_id
  target_spec_code = "dds.mongodb.c6.2xlarge.4.mongos"
  target_type      = "mongos"
  target_id        = var.node_id
}
```

### Update instance's shard group flavor

```hcl
variable "instance_id" {}
variable "group_id" {}

resource "huaweicloud_dds_instance_flavor_update" "test" {
  instance_id      = var.instance_id
  target_spec_code = "dds.mongodb.s6.medium.4.shard"
  target_type      = "shard"
  target_id        = var.group_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this creates a new resource.

* `instance_id` - (Required, String, NonUpdatable) Specifies the ID of a DDS instance.

* `target_spec_code` - (Required, String, NonUpdatable) Specifies the resource specification code of the new specification.

* `target_type` - (Optional, String, NonUpdatable) Specifies the type of the instance. Value options: **mongos**,
  **shard**, **config**, **readonly**.
  + This parameter is mandatory for a cluster instance.
      - To change the specifications of a mongos node, set this parameter to **mongos**.
      - To change the specifications of a single shard or multiple shards in batches, set this parameter to **shard**.
      - To change the specifications of a config node, set this parameter to **config**.
  + This parameter is not specified for replica set instances. If you modify the specifications of a read replica, the
    value is **readonly**.
  + This parameter is not specified for single node instances.

* `target_id` - (Optional, String, NonUpdatable) Specifies the ID of the node or instance whose specifications are to be
  modified.
  + For a cluster instance:
      - When you change the specifications of a mongos node, the value is the mongos node ID.
      - When you change the specifications of a single shard group, the value is the shard group ID.
      - When you change the specifications of multiple shard groups in batches, this parameter is not specified.
      - When you change the specifications of a config group, the value is the config group ID.
  + For a replica set instance, the value is the DB instance ID. If you modify the specifications of a read replica, the
    value is the read replica ID.
  + For a single node instance, the value is the DB instance ID.

* `target_ids` - (Optional, List, NonUpdatable) Specifies the IDs of the node groups whose specifications are to be modified.
  + For a cluster instance:
      - This parameter is not transferred when the specifications of a mongos node are to be changed when the specifications
        of a single shard group are to be changed, or when the specifications of a config group are to be changed.
      - When you change the specifications of multiple shard groups in batches, the value is the IDs of the shard groups.
        A maximum of 16 shard groups can be selected in batches.
  + This parameter is not specified for replica set instances.
  + This parameter is not specified for single node instances.

## Attribute Reference

* `id` - Indicates the resource ID. It's same as the instance ID.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 60 minutes.
