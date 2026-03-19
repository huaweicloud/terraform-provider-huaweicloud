---
subcategory: "Cloud Firewall (CFW)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cfw_schedules"
description: |-
  Use this data source to get the list of CFW schedules.
---

# huaweicloud_cfw_schedules

Use this data source to get the list of CFW schedules.

## Example Usage

```hcl
variable object_id {}

data "huaweicloud_cfw_schedules" "test" {
  object_id = var.object_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `object_id` - (Required, String) Specifies the protected object ID. This ID is used to distinguish
  between Internet boundary protection and VPC boundary protection after the cloud firewall is created.
  You can get this value from data source `huaweicloud_cfw_firewalls`.

* `name` - (Optional, String) Specifies the schedule name.

* `description` - (Optional, String) Specifies the schedule description.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `records` - The schedule records.

  The [records](#data_records_struct) structure is documented below.

<a name="data_records_struct"></a>
The `records` block supports:

* `schedule_id` - The schedule ID.

* `name` - The schedule name.

* `description` - The schedule description.

* `ref_count` - The number of citations.

* `periodic` - The periodic planning.

  The [periodic](#records_periodic_struct) structure is documented below.

* `absolute` - The absolute time planning.

  The [absolute](#records_absolute_struct) structure is documented below.

<a name="records_periodic_struct"></a>
The `periodic` block supports:

* `type` - The periodic type. Valid values are:
  + `0`: each day
  + `1`: some day of each week
  + `2`: a period time of each week

* `start_time` - The start time of the periodic plan.

* `end_time` - The end time of the periodic plan.

* `week_mask` - The days of each week.

* `start_week` - The start date of the weekly periodic plan.

* `end_week` - The end date of the weekly periodic plan.

<a name="records_absolute_struct"></a>
The `absolute` block supports:

* `start_time` - The absolute start time of the plan, in milliseconds (timestamps).

* `end_time` - The absolute ent time of the plan, in milliseconds (timestamps).
