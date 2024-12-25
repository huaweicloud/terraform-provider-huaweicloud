---
subcategory: "GaussDB"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_gaussdb_opengauss_restore"
description: |-
  Use this resource to restore a GaussDB OpenGauss instance with a backup within HuaweiCloud.
---

# huaweicloud_gaussdb_opengauss_restore

Use this resource to restore a GaussDB OpenGauss instance with a backup within HuaweiCloud.

-> **NOTE:** Deleting restoration record is not supported. If you destroy a resource of restoration record,
the restoration record is only removed from the state, but it remains in the cloud. And the instance doesn't return to
the state before restoration.

## Example Usage

### restore by backup_id

```hcl
variable "target_instance_id" {}
variable "source_instance_id" {}
variable "backup_id" {}

resource "huaweicloud_gaussdb_opengauss_restore" "test" {
  target_instance_id = var.target_instance_id
  source_instance_id = var.source_instance_id
  type               = "backup"
  backup_id          = var.backup_id
}
```

### restore by timestamp

```hcl
variable "target_instance_id" {}
variable "source_instance_id" {}

resource "huaweicloud_gaussdb_opengauss_restore" "test" {
  target_instance_id = var.target_instance_id
  source_instance_id = var.source_instance_id
  type               = "timestamp"
  restore_time       = 1673852043000
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) The region in which to create the GaussDB OpenGauss restore resource. If omitted,
  the provider-level region will be used. Changing this creates a new resource.

* `target_instance_id` - (Required, String, ForceNew) Specifies the target instance ID.

  Changing this creates a new resource.

* `source_instance_id` - (Required, String, ForceNew) Specifies the source instance ID.

  Changing this creates a new resource.

* `type` - (Required, String, ForceNew) Specifies the restoration type. Value options:
  + **backup**: indicates restoration from backup files. In this mode, `backup_id` is mandatory.
  + **timestamp**: indicates point-in-time restoration. In this mode, `restore_time` is mandatory.

  Changing this creates a new resource.

* `backup_id` - (Optional, String, ForceNew) Specifies the ID of the backup to be restored. It indicates the ID of the
  full backup corresponding to schema_type. This parameter must be specified when the backup file is used for restoration.

  Changing this creates a new resource.

* `restore_time` - (Optional, String, ForceNew) Specifies the time point of data restoration in the UNIX timestamp format.
  The unit is millisecond and the time zone is UTC.

  Changing this creates a new resource.

## Attribute Reference

In addition to all arguments above, the following attribute is exported:

* `id` - The resource ID. The value is the restore job ID.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 60 minutes.
