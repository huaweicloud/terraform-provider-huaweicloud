---
subcategory: "Data Admin Service (DAS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_das_binlogs"
description: |-
  Use this data source to get the list of DAS binlog files.
---

# huaweicloud_das_binlogs

Use this data source to get the list of DAS binlog files.

## Basic Usage

```hcl
variable "user_id" {}

data "huaweicloud_das_binlogs" "test" {
  user_id     = var.user_id
  binlog_type = "latest"
}
```

## Filter with time

```hcl
variable "user_id" {}

data "huaweicloud_das_binlogs" "test" {
  user_id     = var.user_id
  binlog_type = "latest"
  start_time  = "2025-06-01T00:00:00+08:00"
  end_time    = "2025-06-02T00:00:00+08:00"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region where the binlogs are located.  
  If omitted, the provider-level region will be used.

* `user_id` - (Required, String) Specifies the database user ID.

* `binlog_type` - (Required, String) Specifies the binlog file type.  
  The valid values are as follows:
  + **latest**. Most recent logs.
  + **backup**. Archived logs.

* `start_time` - (Optional, String) Specifies the start time of the query range, in RFC3339 format.

* `end_time` - (Optional, String) Specifies the end time of the query range, in RFC3339 format.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `binlogs` - The list of binlog files that matched the filter parameters.  
  The [binlogs](#binlogs_attr) structure is documented below.

<a name="binlogs_attr"></a>
The `binlogs` block supports:

* `file_name` - The file name.

* `backup_id` - The ID that has already been backed up.

* `file_size` - The file size.

* `task_info` - The archive log parse information.  
  The [task_info](#task_info_attr) structure is documented below.

<a name="task_info_attr"></a>
The `task_info` block supports:

* `id` - The task ID.

* `created_at` - The task creation time, in RFC3339 format.

* `updated_at` - The task modification time, in RFC3339 format.

* `project_id` - The tenant ID of the task.

* `project_name` - The tenant name of the task.

* `user_id` - The user ID of the task.

* `user_name` - The user name of the task.

* `connection_id` - The connection ID of the task.

* `binlog_type` - The binlog type of the task.

* `file_name` - The binlog file name of the task.

* `backup_id` - The backup file ID of the task.

* `status` - The status of the task.

* `err_msg` - The error message.
