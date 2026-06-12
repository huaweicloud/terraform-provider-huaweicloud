---
subcategory: "DataArts Studio"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dataarts_factory_alarm_records"
description: |-
  Use this data source to query the alarm notification records within HuaweiCloud.
---

# huaweicloud_dataarts_factory_alarm_records

Use this data source to query the alarm notification records within HuaweiCloud.

## Example Usage

### Query all alarm records without any filter

```hcl
data "huaweicloud_dataarts_factory_alarm_records" "test" {}
```

### Query alarm records by workspace and time range

```hcl
variable "workspace_id" {}
variable "start_time" {}
variable "end_time" {}

data "huaweicloud_dataarts_factory_alarm_records" "test" {
  workspace  = var.workspace_id
  start_time = var.start_time
  end_time   = var.end_time
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region where the alarm records are located.  
  If omitted, the provider-level region will be used.

* `workspace` - (Optional, String) Specifies the workspace ID.

* `start_time` - (Optional, String) Specifies the start time of the alarm records, in RFC3339 format.  
  If omitted, it defaults to one hour before the current time.

* `end_time` - (Optional, String) Specifies the end time of the alarm records, in RFC3339 format.  
  If omitted, it defaults to the current time. Only records within the last week can be queried.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `records` - The list of alarm records.  
  The [records](#dataarts_factory_records_object) structure is documented below.

<a name="dataarts_factory_records_object"></a>
The `records` block supports:

* `alarm_time` - The alarm notification time, in RFC3339 format.

* `job_name` - The name of the job.

* `schedule_type` - The job instance scheduling mode.
  + **0**: Normal scheduling
  + **2**: Manual scheduling
  + **5**: Data supplement
  + **6**: Sub-job scheduling
  + **7**: One-time scheduling

* `send_msg` - The send message.

* `plan_time` - The plan time, in RFC3339 format.

* `remind_type` - The alarm notification type.
  + **0**: Run successfully
  + **1**: Run abnormally/failed
  + **12**: Not completed
  + **13**: Run cancelled

* `send_status` - The send status.
  + **0**: Send successfully
  + **1**: Send failed

* `job_id` - The ID of the job.
