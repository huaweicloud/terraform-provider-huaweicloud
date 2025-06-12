---
subcategory: "Document Database Service (DDS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dds_instance_node_num_update"
description: |-
  Manages a DDS instance node num update resource within HuaweiCloud.
---

# huaweicloud_dds_instance_node_num_update

Manages a DDS instance node num update resource within HuaweiCloud.

## Example Usage

### Update instance's mongos node num

```hcl
variable "instance_id" {}

resource "huaweicloud_dds_instance_node_num_update" "test" {
  instance_id = var.instance_id
  type        = "shard"
  spec_code   = "dds.mongodb.s6.medium.4.mongos"
  num         = "2"
}
```

### Update instance's shard node num

```hcl
variable "instance_id" {}

resource "huaweicloud_dds_instance_node_num_update" "test" {
  instance_id = var.instance_id
  type        = "shard"
  spec_code   = "dds.mongodb.s6.medium.4.shard"
  num         = "2"
  
  volume {
    size = "20"
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this creates a new resource.

* `instance_id` - (Required, String, NonUpdatable) Specifies the ID of a DDS instance.

* `type` - (Required, String, NonUpdatable) Specifies the type of the object to be scaled. Value options:
  + **mongos**: mongos nodes are to be added.
  + **shard**: shard nodes are to be added.

* `spec_code` - (Required, String, NonUpdatable) Specifies the resource specification code.

* `num` - (Required, String, NonUpdatable) Specifies the number of mongos or shard nodes to be added. A cluster instance
  supports up to 32 mongos nodes and 32 shard nodes.

* `volume` - (Optional, List, NonUpdatable) Specifies the disk capacity of all new shards.
  + This parameter is not transferred when the mongos nodes are to be added.
  + This parameter is mandatory when the shard nodes are to be added.

  The [volume](#volume_struct) structure is documented below.

<a name="volume_struct"></a>
The `volume` block supports:

* `size` - (Required, String, NonUpdatable) Specifies the disk capacity of all new shards. Value range:
  + 10 GB to 5,000 GB when the shard node has fewer than 8 vCPUs.
  + 10 GB to 10,000 GB when the shard node has 8 or more vCPUs.

## Attribute Reference

* `id` - Indicates the resource ID. It's same as the instance ID.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 60 minutes.
