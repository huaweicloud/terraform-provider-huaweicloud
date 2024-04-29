---
subcategory: "Distributed Cache Service (DCS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dcs_backup"
description: ""
---

# huaweicloud_dcs_backup

Manages a DCS backup resource within HuaweiCloud.

## Example Usage

```hcl
variable "dcs_instance_id" {}

resource "huaweicloud_dcs_backup" "test"{
  instance_id = var.dcs_instance_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `instance_id` - (Required, String, ForceNew) Specifies the ID of the DCS instance.

  Changing this parameter will create a new resource.

* `description` - (Optional, String, ForceNew) Specifies the description of DCS instance backup.

  Changing this parameter will create a new resource.

* `backup_format` - (Optional, String, ForceNew) Specifies the format of the DCS instance backup.
  Value options: **aof**, **rdb**. Default to rdb.

  Changing this parameter will create a new resource.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID in format of `<instance_id>/<backup_id>`.

* `backup_id` - Indicates the ID of the DCS instance backup.

* `name` - Indicates the backup name.

* `size` - Indicates the size of the backup file (byte).

* `type` - Indicates the backup type. Valid value:
  + **manual**: indicates manual backup.
  + **auto**: indicates automatic backup.

* `begin_time` - Indicates the time when the backup task is created. The format is yyyy-mm-dd hh:mm:ss.
  The value is in UTC format.

* `end_time` - Indicates the time at which DCS instance backup is completed. The format is yyyy-mm-dd hh:mm:ss.
  The value is in UTC format.

* `status` - Indicates the backup status. Valid value:
  + **waiting**: The task is waiting to begin.
  + **backuping**: DCS instance backup is in progress.
  + **succeed**: DCS instance backup succeeded.
  + **failed**: DCS instance backup failed.
  + **expired**: The backup file has expired.
  + **deleted**: The backup file has been deleted manually.

* `is_support_restore` - Indicates whether restoration is supported. Value Options: **TRUE**, **FALSE**.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 30 minutes.
* `delete` - Default is 10 minutes.

## Import

The DCS backup can be imported using the `instance_id` and `backup_id` separated by a slash, e.g.:

```bash
$ terraform import huaweicloud_dcs_backup.test <instance_id>/<backup_id>
```
