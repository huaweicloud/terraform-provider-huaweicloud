---
subcategory: "TaurusDB"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_taurusdb_instance_backups"
description: |-
  Use this data source to query the full backups of a specified TaurusDB instance within HuaweiCloud.
---

# huaweicloud_taurusdb_instance_backups

Use this data source to query the full backups of a specified TaurusDB instance within HuaweiCloud.

## Example Usage

### Query all backups of an instance

```hcl
variable "instance_id" {}

data "huaweicloud_taurusdb_instance_backups" "test" {
  instance_id = var.instance_id
}
```

### Query backups with filter

```hcl
variable "instance_id" {}
variable "backup_name" {}

data "huaweicloud_taurusdb_instance_backups" "test" {
  instance_id    = var.instance_id
  filter_field   = "name"
  filter_content = var.backup_name
}
```

### Query backups with order

```hcl
variable "instance_id" {}

data "huaweicloud_taurusdb_instance_backups" "test" {
  instance_id = var.instance_id
  order_field = "beginTime"
  order_rule  = "asc"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the ID of the TaurusDB instance.

* `order_field` - (Optional, String) Specifies the field based on which the query results are sorted.
  Value options:
  + **name**: Backup name.
  + **beginTime**: Backup start time.
  + **type**: Backup type.

  -> **NOTE:** This field must be used together with `order_rule`.

* `order_rule` - (Optional, String) Specifies the sorting rule.
  Value options:
  + **asc**: Ascending order.
  + **desc**: Descending order.

  -> **NOTE:** This field must be used together with `order_field`.

* `filter_field` - (Optional, String) Specifies the filter field type.
  Value options:
  + **name**: Filter backups by backup name.

  -> **NOTE:** This field must be used together with `filter_content`.

* `filter_content` - (Optional, String) Specifies the filter content.

  -> **NOTE:** This field must be used together with `filter_field`.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `backups` - The list of backups.

  The [backups](#instance_backups_struct) structure is documented below.

<a name="instance_backups_struct"></a>
The `backups` block supports:

* `id` - The ID of the backup.

* `name` - The name of the backup.

* `description` - The description of the backup.

* `instance_id` - The ID of the TaurusDB instance.

* `instance_name` - The name of the TaurusDB instance.

* `size` - The backup size.

* `size_unit` - The unit of the backup size (KB).

* `status` - The backup status.
  Value options:
  + **BUILDING**: Backup in progress.
  + **COMPLETED**: Backup completed.
  + **FAILED**: Backup failed.
  + **DELETING**: Backup being deleted.

* `created` - The creation time.

* `updated` - The update time.

* `backup_type` - The backup type.
  Value options:
  + **Db**: Physical backup.
  + **Snapshot**: Snapshot backup.

* `backup_level` - The backup level.
  Value options:
  + **1**: Level-1 backup.
  + **2**: Level-2 backup.

* `backup_method` - The backup method.
  Value options:
  + **Db**: Physical backup.
  + **Snapshot**: Snapshot backup.

* `use_detail` - The usage details.

* `time_zone` - The UTC time zone.
