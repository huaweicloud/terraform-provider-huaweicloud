---
subcategory: "Data Admin Service (DAS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_das_binlog_parse_task_export"
description: |-
  Manages a DAS binlog parse task export resource within HuaweiCloud.
---

# huaweicloud_das_binlog_parse_task_export

Manages a DAS binlog parse task export resource within HuaweiCloud.

## Example Usage

### Export with filter conditions

```hcl
variable "user_id" {}
variable "task_id" {}
variable "bucket_name" {}

resource "huaweicloud_das_binlog_parse_task_export" "test" {
  user_id     = var.user_id
  task_id     = var.task_id
  bucket_name = var.bucket_name

  filter_condition {
    db_names             = ["test_db"]
    tb_names             = ["test_table"]
    types                = ["insert", "update"]
    parse_double_insert  = true

    columns {
      name  = "name"
      value = "terraform"
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region where the binlog parse task export is located.  
  If omitted, the provider-level region will be used.  
  Changing this parameter will create a new resource.

* `user_id` - (Required, String, NonUpdatable) Specifies the user ID of the database connection.

  -> You can use `huaweicloud_das_database_users` data source to get user IDs.

* `task_id` - (Required, Int, NonUpdatable) Specifies the binlog parse task ID.

* `bucket_name` - (Required, String, NonUpdatable) Specifies the OBS bucket name.

* `filter_condition` - (Required, List, NonUpdatable) Specifies the filter conditions for the export task.
  The [filter_condition](#filter_condition) structure is documented below.

<a name="filter_condition"></a>
The `filter_condition` block supports:

* `db_names` - (Optional, List, NonUpdatable) Specifies the list of database names to filter.

* `tb_names` - (Optional, List, NonUpdatable) Specifies the list of table names to filter.

* `file_names` - (Optional, List, NonUpdatable) Specifies the list of file names to filter.

* `start_time` - (Optional, String, NonUpdatable) Specifies the start time of the export range, in RFC3339 format.

* `end_time` - (Optional, String, NonUpdatable) Specifies the end time of the export range, in RFC3339 format.

* `types` - (Optional, List, NonUpdatable) Specifies the list of SQL types to filter.  
  The valid values are as follows:
  + **insert**
  + **update**
  + **delete**
  + **ddl**

* `columns` - (Optional, List, NonUpdatable) Specifies the list of columns to filter.  
  The [columns](#binlog_parse_task_export_columns) structure is documented below.

* `parse_double_insert` - (Optional, Bool, NonUpdatable) Specifies whether to export **UPDATE** statements as
  two **INSERT** statements.

<a name="binlog_parse_task_export_columns"></a>
The `columns` block supports:

* `name` - (Optional, String, NonUpdatable) Specifies the column name.

* `value` - (Optional, String, NonUpdatable) Specifies the column value.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The export task ID.

* `status` - The task status.
  + **0**: Initialized
  + **1**: Running
  + **2**: Partially successful
  + **3**: Successful
  + **4**: Failed
  + **-1**: Deleted

* `instance_id` - The instance ID.

* `last_record_time` - The last record time, in RFC3339 format.

* `created_at` - The task creation time, in RFC3339 format.

* `export_line_num` - The number of exported lines.

* `download_url` - The download URL of the exported file.

## Import

The DAS binlog parse task export can be imported using `<user_id>/<bucket_name>/<id>`, e.g.

```bash
$ terraform import huaweicloud_das_binlog_parse_task_export.test <user_id>/<bucket_name>/<id>
```
