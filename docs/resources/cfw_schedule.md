---
subcategory: "Cloud Firewall (CFW)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cfw_schedule"
description: |-
  Manages a CFW schedule resource within HuaweiCloud.
---

# huaweicloud_cfw_schedule

Manages a CFW schedule resource within HuaweiCloud.

## Example Usage

```hcl
variable "object_id" {}

resource "huaweicloud_cfw_schedule" "test" {
  object_id   = var.object_id
  name        = "test-name"
  description = "test description"

  absolute {
    end_time   = 1774076220000
    start_time = 1773730600349
  }

  periodic {
    type       = 0
    start_time = "00:00:00"
    end_time   = "23:59:59"
  }

  periodic {
    type       = 1
    start_time = "00:00:00"
    end_time   = "23:59:59"
    week_mask = [
      1,
    ]
  }

  periodic {
    type       = 2
    start_time = "00:00:00"
    end_time   = "23:59:59"
    start_week = 1
    end_week   = 1
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this will create new resource.

* `object_id` - (Required, String, NonUpdatable) Specifies the protected object ID. This ID is used to distinguish
  between Internet boundary protection and VPC boundary protection after the cloud firewall is created.

* `name` - (Required, String) Specifies the schedule name.

* `description` - (Optional, String) Specifies the schedule description.

* `periodic` - (Optional, List) Specifies the periodic planning.

  The [periodic](#periodic_struct) structure is documented below.

* `absolute` - (Optional, List) Specifies the absolute time planning. The max length of this field is `1`.

  The [absolute](#absolute_struct) structure is documented below.

-> The fields `periodic` and `absolute` cannot both be empty.

<a name="periodic_struct"></a>
The `periodic` block supports:

* `type` - (Required, Int) Specifies the periodic type. Valid values are:
  + `0`: each day
  + `1`: some day of each week
  + `2`: a period time of each week

* `start_time` - (Required, String) Specifies the start time of the periodic plan.

* `end_time` - (Required, String) Specifies the end time of the periodic plan.

* `week_mask` - (Optional, List) Specifies the days of each week.

* `start_week` - (Optional, Int) Specifies the start date of the weekly periodic plan.

* `end_week` - (Optional, Int) Specifies the end date of the weekly periodic plan.

<a name="absolute_struct"></a>
The `absolute` block supports:

* `start_time` - (Optional, Int) Specifies the absolute start time of the plan, in milliseconds (timestamps).

* `end_time` - (Optional, Int) Specifies the absolute ent time of the plan, in milliseconds (timestamps).

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID (also the schedule ID).

* `ref_count` - The number of citations.

## Import

The resource can be imported using `object_id`, `id` (schedule ID), separated by a slash, e.g.

```bash
$ terraform import huaweicloud_cfw_schedule.test <object_id>/<id>
```
