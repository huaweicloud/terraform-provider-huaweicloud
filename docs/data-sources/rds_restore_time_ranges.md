---
subcategory: "Relational Database Service (RDS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_rds_restore_time_ranges"
description: |-
  Use this data source to get the list of RDS restore time ranges.
---

# huaweicloud_rds_restore_time_ranges

Use this data source to get the list of RDS restore time ranges.

## Example Usage

```hcl
variable "instance_id" {}

data "huaweicloud_rds_restore_time_ranges" "test" {
  instance_id = var.instance_id
  date        = "2024-05-15"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the ID of RDS instance.

* `date` - (Optional, String) Specifies the date to be queried.
  The value is in the **yyyy-mm-dd** format, and the time zone is UTC.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `restore_time` - Indicates the list of restoration time ranges.

  The [restore_time](#restore_time_struct) structure is documented below.

<a name="restore_time_struct"></a>
The `restore_time` block supports:

* `start_time` - Indicates the start time of the restoration time range in the UNIX timestamp format.
  The unit is millisecond and the time zone is UTC.

* `end_time` - Indicates the end time of the restoration time range in the UNIX timestamp format.
  The unit is millisecond and the time zone is UTC.
