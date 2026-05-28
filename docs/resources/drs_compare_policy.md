---
subcategory: "Data Replication Service (DRS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_drs_compare_policy"
description: |-
  Manages a DRS compare policy resource within HuaweiCloud.
---

# huaweicloud_drs_compare_policy

Manages a DRS compare policy resource within HuaweiCloud.

-> Deleting this resource will close the compare policy.

## Example Usage

```hcl
variable "job_id" {}

resource "huaweicloud_drs_compare_policy" "test" {
  job_id         = var.job_id
  period         = "* * 1,3,5"
  begin_time     = "00:00:00"
  end_time       = "04:00:00"
  compare_type   = ["lines"]
  compare_policy = "normal"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `job_id` - (Required, String, NonUpdatable) Specifies the job ID.

* `period` - (Optional, String) Specifies the comparison period.
  Weekly comparison: format example `* * 1,3,5`, where `1,3,5` corresponds to Monday, Wednesday, Friday.
  Daily comparison: fixed value `* * 1,2,3,4,5,6,7`.
  Hourly comparison: fixed value `* * 1,2,3,4,5,6,7`.

* `begin_time` - (Optional, String) Specifies the start time when the comparison policy takes effect, UTC time,
  `24`-hour format, time format HH:mm:ss, e.g. **00:00:00**.

* `end_time` - (Optional, String) Specifies the end time when the comparison policy takes effect, UTC time,
  `24`-hour format, time format HH:mm:ss, e.g. **04:00:00**.

* `compare_type` - (Optional, List of String) Specifies the list of comparison types.
  Valid values are: **object_comparison**, **lines**, and **account**.

* `compare_policy` - (Optional, String) Specifies the comparison policy.
  Valid values are **normal** and **manyToOne**.

* `interval_hour` - (Optional, Int) Specifies the interval time, filled in when comparing by hour,
  indicating how often to perform comparison, unit is hour.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID (also `job_id`).

* `status` - The status of the comparison policy.
  Valid values are: **OPEN**, **CLOSED**, and **NO_SUPPORT**.

* `next_compare_time` - The next comparison time, UTC time, e.g. **2023-06-12T08:00:00Z**.

## Import

The DRS compare policy can be imported by `id`. e.g.

```bash
$ terraform import huaweicloud_drs_compare_policy.test <id>
```
