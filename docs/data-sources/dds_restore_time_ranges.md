---
subcategory: "Document Database Service (DDS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dds_restore_time_ranges"
description: |-
  Use this data source to get the list of DDS instance restore time ranges.
---

# huaweicloud_dds_restore_time_ranges

Use this data source to get the list of DDS instance restore time ranges.

## Example Usage

```hcl
variable "instance_id" {}
variable "date" {}

data "huaweicloud_dds_restore_time_ranges" "test"{
  instance_id = var.instance_id
  date        = var.date
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the instance ID.

* `date` - (Required, String) Specifies the date to be queried.
  The value is in the yyyy-mm-dd format, and the time zone is UTC.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `restore_times` - Indicates the restoration time ranges.

  The [restore_times](#restore_times_struct) structure is documented below.

<a name="restore_times_struct"></a>
The `restore_times` block supports:

* `start_time` - Indicates the start time in the UNIX timestamp format.
  The unit is millisecond and the time zone is UTC.

* `end_time` - Indicates the end time in the UNIX timestamp format.
  The unit is millisecond and the time zone is UTC.
