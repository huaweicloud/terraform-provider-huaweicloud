---
subcategory: "Data Admin Service (DAS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_das_slow_log_export_task"
description: |-
  Manages a DAS slow log export task resource within HuaweiCloud.
---

# huaweicloud_das_slow_log_export_task

Manages a DAS slow log export task resource within HuaweiCloud.

-> This resource is a one-time action resource for exporting slow log tasks. Deleting this resource will not clear the
   corresponding request record, but will only remove the resource information from the tfstate file.

## Example Usage

### Basic Usage

```hcl
variable "instance_id" {}
variable "obs_bucket_name" {}

resource "huaweicloud_das_slow_log_export_task" "test" {
  instance_id = var.instance_id
  bucket_name = var.obs_bucket_name
  start_time  = "2026-06-01T00:00:00+08:00"
  end_time    = "2026-06-02T00:00:00+08:00"
}
```

### Exporting task with execute time

```hcl
variable "instance_id" {}
variable "obs_bucket_name" {}

resource "huaweicloud_das_slow_log_export_task" "test" {
  instance_id          = var.instance_id
  bucket_name          = var.obs_bucket_name
  start_time           = "2026-06-01T00:00:00+08:00"
  end_time             = "2026-06-02T00:00:00+08:00"
  export_type          = "slowsqldetails"
  execute_time_min     = 100
  execute_time_max     = 3600
  min_avg_execute_time = 1
  max_avg_execute_time = 100
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region where the slow log export task is located.  
  If omitted, the provider-level region will be used.  
  Changing this parameter will create a new resource.

* `instance_id` - (Required, String, NonUpdatable) Specifies the instance ID.

* `bucket_name` - (Required, String, NonUpdatable) Specifies the OBS bucket name.

* `start_time` - (Required, String, NonUpdatable) Specifies the start time, in RFC3339 format.

* `end_time` - (Required, String, NonUpdatable) Specifies the end time, in RFC3339 format.

-> 1.The earliest `start_time` is at most `2` days earlier than the current time.  
2.The latest `end_time` is at most `1` day later than the current time.  
3.The `end_time` must be greater than the `start_time`.

* `file_path` - (Optional, String, NonUpdatable) Specifies the OBS file directory.  
  The maximum length is `1024` characters.

* `export_type` - (Optional, String, NonUpdatable) Specifies the export type. Defaults to **slowsql**.  
  The valid values are as follows:
  + **slowsql**: Export slow SQL summary.
  + **slowsqldetails**: Export slow SQL details.

* `sort_field` - (Optional, String, NonUpdatable) The sort field.  
  The valid values are as follows:
  + **occurrenceTime**. Execution start time.
  + **executeTime**. Execution loss time.
  + **lockWaitTime**. Lock wait time.
  + **rowsExamined**. The scanned row.
  + **rowsSent**. The returned row.

* `sort_asc` - (Optional, Bool, NonUpdatable) Whether to sort in ascending order.  
  The valid values are as follows:
  + **true**. Ascending.
  + **false**. Descending.

* `user_name` - (Optional, String, NonUpdatable) Specifies the user name.

* `client_ip_address` - (Optional, String, NonUpdatable) Specifies the client IP address.

* `killed` - (Optional, String, NonUpdatable) Specifies the execution status.

* `execute_time_min` - (Optional, Int, NonUpdatable) Specifies the minimum execution time, in milliseconds.

* `execute_time_max` - (Optional, Int, NonUpdatable) Specifies the maximum execution time, in milliseconds.

* `min_avg_execute_time` - (Optional, Float, NonUpdatable) Specifies the minimum average execution time.

* `max_avg_execute_time` - (Optional, Float, NonUpdatable) Specifies the maximum average execution time.

* `rows_max_examined` - (Optional, Int, NonUpdatable) Specifies the maximum number of scanned rows.

* `rows_min_examined` - (Optional, Int, NonUpdatable) Specifies the minimum number of scanned rows.

* `fuzzy_sql` - (Optional, String, NonUpdatable) Specifies the fuzzy SQL.

* `operation` - (Optional, String, NonUpdatable) Specifies the operation type (can be combined with commas).

* `time_zone` - (Optional, String, NonUpdatable) Specifies the time zone.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The export task ID.
