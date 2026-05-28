---
subcategory: "Data Admin Service (DAS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_das_binlog_parse_task_exports"
description: |-
  Use this data source to get the list of DAS binlog exported tasks.
---

# huaweicloud_das_binlog_parse_task_exports

Use this data source to get the list of DAS binlog exported tasks.

## Basic Usage

```hcl
variable "user_id" {}

data "huaweicloud_das_binlog_parse_task_exports" "test" {
  user_id = var.user_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region where the binlog parse task exports are located.  
  If omitted, the provider-level region will be used.

* `user_id` - (Required, String) Specifies the database user ID.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `tasks` - The list of binlog parse task exports.  
  The [tasks](#binlog_parse_task_exports_tasks) structure is documented below.

<a name="binlog_parse_task_exports_tasks"></a>
The `tasks` block supports:

* `exported_task_id` - The exported task ID.

* `parsed_task_id` - The parsed task ID.

* `instance_id` - The instance ID.

* `status` - The task status.
  + **0**. Initializing.
  + **1**. Running.
  + **2**. Partially successful.
  + **3**. Successful.
  + **4**. Failed.
  + **-1**. Deleted.

* `start_time` - The start time, in RFC3339 format.

* `end_time` - The end time, in RFC3339 format.

* `last_record_time` - The last record time, in RFC3339 format.

* `created_at` - The task creation time, in RFC3339 format.

* `export_line_num` - The number of exported lines.

* `download_url` - The download URL of the exported file.

* `source_file_name` - The binlog source file name.
