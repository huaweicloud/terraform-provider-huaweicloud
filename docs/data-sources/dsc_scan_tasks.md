---
subcategory: "Data Security Center (DSC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dsc_scan_tasks"
description: |-
  Use this data source to get the list of scan tasks for a specified scan job within HuaweiCloud.
---

# huaweicloud_dsc_scan_tasks

Use this data source to get the list of scan tasks for a specified scan job within HuaweiCloud.

## Example Usage

```hcl
variable "job_id" {}

data "huaweicloud_dsc_scan_tasks" "test" {
  job_id = var.job_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `job_id` - (Required, String) Specifies the scan job ID.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID in UUID format.

* `tasks` - The scan task list.
  The [tasks](#dsc_scan_tasks) structure is documented below.

<a name="dsc_scan_tasks"></a>
The `tasks` block supports:

* `id` - The task ID.

* `category` - The asset type.

* `status` - The task status.

* `progress` - The task progress.

* `asset_name` - The asset name.

* `asset_id` - The asset ID.

* `start_time` - The task start time.

* `end_time` - The task end time.

* `scanned_object_num` - The number of scanned objects.

* `to_be_scanned_object_num` - The number of objects to be scanned.

* `scan_speed` - The scan speed.

* `skip_object_num` - The number of skipped objects.

* `last_scan_risk` - The last scan risk result.

* `security_level_name` - The security level name.

* `security_level_color` - The security level color.
