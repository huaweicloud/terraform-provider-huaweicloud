---
subcategory: "GaussDB(for MySQL)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_gaussdb_mysql_incremental_backups"
description: |-
  Use this data source to get the list of incremental backups
---

# huaweicloud_gaussdb_mysql_incremental_backups

Use this data source to get the list of incremental backups

## Example Usage

```hcl
variable "instance_id" {}

data "huaweicloud_gaussdb_mysql_incremental_backups" "test" {
  instance_id = var.instance_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the instance ID.

* `begin_time` - (Optional, String) Specifies the query start time.
  The format is **yyyy-mm-ddThh:mm:ssZ**. It is mandatory when `end_time` is set.

* `end_time` - (Optional, String) Specifies the query end time.
  The format is **yyyy-mm-ddThh:mm:ssZ** and the end time must be later than the start time.
  It is mandatory when `begin_time` is set.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `backups` - Indicates the list of incremental backups.

  The [backups](#backups_struct) structure is documented below.

<a name="backups_struct"></a>
The `backups` block supports:

* `id` - Indicates the backup ID.

* `name` - Indicates the backup name.

* `instance_id` - Indicates the instance ID.

* `begin_time` - Indicates the backup start time in the **yyyy-mm-ddThh:mm:ssZ** format.

* `end_time` - Indicates the backup end time in the **yyyy-mm-ddThh:mm:ssZ** format.

* `size` - Indicates the backup size, in KB.
