---
subcategory: "Workspace"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_workspace_app_schedule_task_future_executions"
description: |-
  Use this data source to get the future execution time list of Workspace APP schedule task within HuaweiCloud.
---

# huaweicloud_workspace_app_schedule_task_future_executions

Use this data source to get the future execution time list of Workspace APP schedule task within HuaweiCloud.

## Example Usage

```hcl
variable "scheduled_time" {}

data "huaweicloud_workspace_app_schedule_task_future_executions" "test" {
  scheduled_type = "DAY"
  scheduled_time = var.scheduled_time
  day_interval   = 1
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region where the schedule task is located.  
  If omitted, the provider-level region will be used.

* `scheduled_type` - (Required, String) Specifies the type of execution cycle.  
  The valid values are as follows:
  + **FIXED_TIME**
  + **DAY**
  + **WEEK**
  + **MONTH**

* `scheduled_time` - (Required, String) Specifies the execution time of the schedule task.  
  The format is `HH:mm:ss`.

* `day_interval` - (Optional, Int) Specifies the interval in days for the scheduled task is to be executed.  
  The valid value ranges from `1` to `31`.  
  This parameter is **required** when `scheduled_type` is set to **DAY**.

* `week_list` - (Optional, String) Specifies the days of the weeks when the scheduled task is to be executed.  
  The valid value ranges from `1` to `7`, separated by commas, e.g. `1,2,7`.  
  `1` means Sunday, `2` means Monday, and so on.  
  This parameter is **required** when `scheduled_type` is set to **WEEK**.

* `month_list` - (Optional, String) Specifies the month when the scheduled task is to be executed.  
  The valid value ranges from `1` to `12`, separated by commas, e.g. `1,3,12`.  
  This parameter is **required** when `scheduled_type` is set to **MONTH**.

* `date_list` - (Optional, String) Specifies the days of month when the scheduled task is to be executed.  
  The valid value ranges from `1` to `31` and `L` (means the last day), separated by commas, e.g. `1,2,28` or `L`.  
  `L` can only be used alone, and cannot be used together with other values.  
  This parameter is **required** when `scheduled_type` is set to **MONTH**.

* `scheduled_date` - (Optional, String) Specifies the fixed date when the scheduled task is to be executed.  
  The format is `YYYY-MM-dd`.  
  This parameter is **required** when `scheduled_type` is set to **FIXED_TIME**.

* `time_zone` - (Optional, String) Specifies the time zone of the schedule task.  
  Defaults to **Asia/Shanghai**.

* `expire_time` - (Optional, String) Specifies the expiration time of the schedule task.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `future_executions` - The list of future execution times that match the filter parameters.
  
* `time_zone` - The time zone corresponding to the execution times.
