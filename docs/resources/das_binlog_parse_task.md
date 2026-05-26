---
subcategory: "Data Admin Service (DAS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_das_binlog_parse_task"
description: |-
  Manages a DAS binlog parse task resource within HuaweiCloud.
---

# huaweicloud_das_binlog_parse_task

Manages a DAS binlog parse task resource within HuaweiCloud.

## Example Usage

### Create a binlog parse task with latest type

```hcl
variable "user_id" {}
variable "file_name" {}

resource "huaweicloud_das_binlog_parse_task" "test" {
  user_id     = var.user_id
  file_name   = var.file_name
  binlog_type = "latest"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region where the binlog parse task is located.  
  If omitted, the provider-level region will be used.  
  Changing this parameter will create a new resource.

* `user_id` - (Required, String, NonUpdatable) Specifies the user ID of the database connection.

  -> You can use `huaweicloud_das_database_users` data source to get user IDs.

* `binlog_type` - (Required, String, NonUpdatable) Specifies the binlog type.  
  The valid values are as follows:
  + **latest**: Recent logs.
  + **backup**: Archived logs.
  + **fragment**: Fragment backup logs.

* `file_name` - (Required, String, NonUpdatable) Specifies the binlog file name.

* `backup_id` - (Optional, String, NonUpdatable) Specifies the backup ID.  
  This parameter is **Required** when the `binlog_type` is set to **backup**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The task ID.

* `status` - The task status.
  + **0**: Initialized
  + **1**: Running
  + **2**: Partially successful
  + **3**: Successful
  + **4**: Failed
  + **-1**: Deleted

* `position` - The binlog file parse position.

* `error_message` - The error message.

* `created_at` - The task creation time, in RFC3339 format.

* `updated_at` - The task modification time, in RFC3339 format.

## Import

The DAS binlog parse task can be imported using the `user_id` and `id` separated by a slash, e.g.

```bash
$ terraform import huaweicloud_das_binlog_parse_task.test <user_id>/<id>
```
