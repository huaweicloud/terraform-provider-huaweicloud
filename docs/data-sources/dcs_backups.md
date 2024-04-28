---
subcategory: "Distributed Cache Service (DCS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dcs_backups"
description: ""
---

# huaweicloud_dcs_backups

Use this data source to get the list of DCS backups.

## Example Usage

```hcl
var "instance_id" {}

data "huaweicloud_dcs_backups" "test"{
  instance_id = var.instance_id
  name        = "test_name"
  type        = "auto"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the ID of the DCS instance.

* `begin_time` - (Optional, String) Specifies the start time (UTC) of DCS backups.
  The format is **yyyymmddhhmmss**, for example, 20231201000000.

* `end_time` - (Optional, String) Specifies the end time (UTC) of DCS backups.
  The format is **yyyymmddhhmmss**, for example, 20231201000000.

* `backup_id` - (Optional, String) Specifies the ID of the DCS instance backup.

* `name` - (Optional, String) Specifies the backup name.

* `type` - (Optional, String) Specifies the backup type.
  Value options: **manual**, **auto**.

* `status` - (Optional, String) Specifies the backup status.
  Value options: **waiting**, **backuping**, **succeed**, **failed**, **expired**, **deleted**.

* `is_support_restore` - (Optional, String) Specifies whether restoration is supported.
  Value Options: **TRUE**, **FALSE**.

* `backup_format` - (Optional, String) Specifies the format of the DCS instance backup.
  Value options: **aof**, **rdb**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `backups` - The list of backups.
  The [backups](#DCS_backups) structure is documented below.

<a name="DCS_backups"></a>
The `backups` block supports:

* `id` - Indicates the ID of the DCS instance backup.

* `name` - Indicates the backup name.

* `instance_id` - Indicates the ID of the DCS instance.

* `size` - Indicates the size of the backup file (byte).

* `type` - Indicates the backup type.

* `begin_time` - Indicates the time when the backup task is created. The format is yyyy-mm-dd hh:mm:ss.
  The value is in UTC format.

* `end_time` - Indicates the time at which DCS instance backup is completed. The format is yyyy-mm-dd hh:mm:ss.
  The value is in UTC format.

* `status` - Indicates the backup status.

* `description` - Indicates the description of the DCS instance backup.

* `is_support_restore` - Indicates whether restoration is supported.

* `backup_format` - Indicates the format of the DCS instance backup.

* `error_code` - Indicates the error code displayed for a backup failure.

* `progress` - Indicates the backup progress.
