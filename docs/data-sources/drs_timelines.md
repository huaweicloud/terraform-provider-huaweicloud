---
subcategory: "Data Replication Service (DRS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_drs_timelines"
description: |-
  Use this data source to get the list of operation timelines for a specified DRS job within HuaweiCloud.
---

# huaweicloud_drs_timelines

Use this data source to get the list of operation timelines for a specified DRS job within HuaweiCloud.

## Example Usage

```hcl
variable "job_id" {}

data "huaweicloud_drs_timelines" "test" {
  job_id = var.job_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `job_id` - (Required, String) Specifies the ID of the DRS job.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `timelines` - The list of DRS job operation timelines.

  The [timelines](#timelines_struct) structure is documented below.

<a name="timelines_struct"></a>
The `timelines` block supports:

* `name` - The name of the timeline.

* `status` - The status of the timeline. Valid values are **success** and **failed**.

* `operation_time` - The operation time of the timeline, for example, **2026-05-18T02:28:54Z**.

* `user_name` - The user name of the timeline.
