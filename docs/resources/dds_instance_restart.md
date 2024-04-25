---
subcategory: "Document Database Service (DDS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dds_instance_restart"
description: |-
  Manages a DDS instance restart resource within HuaweiCloud.
---

# huaweicloud_dds_instance_restart

Manages a DDS instance restart resource within HuaweiCloud.

## Example Usage

### Restart insatnce

```hcl
variable instance_id {}

resource "huaweicloud_dds_instance_restart" "test" {
  instance_id = var.instance_id
}
```

### Restart insatnce's mongos node

```hcl
variable "instance_id" {}
variable "node_id" {}

resource "huaweicloud_dds_instance_restart" "test" {
  instance_id = var.instance_id
  target_type = "mongos"
  target_id   = var.node_id
}
```

### Restart insatnce's shard group

```hcl
variable "instance_id" {}
variable "group_id" {}

resource "huaweicloud_dds_instance_restart" "test" {
  instance_id = var.instance_id
  target_type = "shard"
  target_id   = var.group_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this creates a new resource.

* `instance_id` - (Required, String, ForceNew) Specifies the ID of a DDS instance.
  Changing this creates a new resource.

* `target_type` - (Optional, String, ForceNew) Specifies the type of the object to restart. Valid values are **mongos**,
  **shard**, **config**. It's required with `target_id`. Changing this creates a new resource.

* `target_id` - (Optional, String, ForceNew) Specifies the ID of the object to be restarted. When you restart a node in
  a cluster instance, the value is the mongos node ID for a mongos node, and shard or config group ID for a shard or
  config group. It's required with `target_type`.
  Changing this creates a new resource.

-> If you want to restart instance, set both `target_type` and `target_id` empty.

## Attribute Reference

* `id` - Indicates the resource ID. It's same as the instance ID.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 30 minutes.
