---
subcategory: "Document Database Service (DDS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dds_instance_restore"
description: |-
  Manages a DDS instance restore resource within HuaweiCloud.
---

# huaweicloud_dds_instance_restore

Manages a DDS instance restore resource within HuaweiCloud.

## Example Usage

### Restore insatnce by backup ID

```hcl
variable "source_id" {}
variable "backup_id" {}
variable "target_id" {}

resource "huaweicloud_dds_instance_restore" "test" {
  source_id = var.source_id
  backup_id = var.backup_id
  target_id = var.target_id
}
```

### Restore insatnce by time stamp

```hcl
variable "source_id" {}
variable "restore_time" {}
variable "target_id" {}

resource "huaweicloud_dds_instance_restore" "test" {
  source_id    = var.source_id
  restore_time = var.restore_time
  target_id    = var.target_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this creates a new resource.

* `source_id` - (Required, String, ForceNew) Specifies the ID of the source DDS instance.
  Changing this creates a new resource.

* `target_id` - (Required, String, ForceNew) Specifies the ID of the target DDS instance.

* `backup_id` - (Optional, String, ForceNew) Specifies the backup ID of the source DDS instance.
  Changing this creates a new resource.

* `restore_time` - (Optional, String, ForceNew) Specifies the restore time of the source DDS instance. The unit is
  millisecond and the time zone is UTC. This parameter takes effect only for replica set instances.
  Changing this creates a new resource.

-> Only one of `backup_id` and `restore_time`can be set and it will return a validation error if none are specified.

## Attribute Reference

* `id` - Indicates the resource ID. It's same as the target instance ID.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 30 minutes.
