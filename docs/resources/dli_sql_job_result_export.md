---
subcategory: "Data Lake Insight (DLI)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dli_sql_job_result_export"
description: |-
  Use this resource to export query result of the SQL job to the specified OBS bucket within HuaweiCloud.
---

# huaweicloud_dli_sql_job_result_export

Use this resource to export query result of the SQL job to the specified OBS bucket within HuaweiCloud.

-> 1. Only the query result of the **QUERY** type SQL job can be exported.
   <br>2. This resource is a one-time action resource for exporting query result of the SQL job to OBS. Deleting this
   resource will not clear the corresponding request record, but will only remove the resource information from the
   tfstate file.

## Example Usage

```hcl
variable "job_id" {}
variable "data_path" {}

resource "huaweicloud_dli_sql_job_result_export" "test" {
  job_id    = var.job_id
  data_path = var.data_path
  data_type = "json"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region where the SQL job with result to be exported is
  located.  
  If omitted, the provider-level region will be used.  
  Changing this will create a new resource.

* `job_id` - (Required, String, NonUpdatable) Specifies the ID of the SQL job.

* `data_path` - (Required, String, NonUpdatable) Specifies the OBS path for storing the exported data.  
  The path must be specified at the folder level.  
  The folder name cannot contain special characters (`\/:*?"<>|`), and cannot start or end with dot (.).
  e.g. `obs://bucket_name/folder_name`.

* `data_type` - (Required, String, NonUpdatable) Specifies the storage format of the exported data.  
  The valid values are as follows:
  + **csv**
  + **json**

* `compress` - (Optional, String, NonUpdatable) Specifies the compression format of the exported data.  
  The valid values are as follows:
  + **none**
  + **gzip**
  + **bzip2**
  + **deflate**

  Defaults to **none**.

* `queue_name` - (Optional, String, NonUpdatable) Specifies the queue name used to execute the export task.  
  If omitted, the default queue is used.

* `export_mode` - (Optional, String, NonUpdatable) Specifies the mode of data export.  
  The valid values are as follows:
  + **ErrorIfExists**: The specified export directory must not exist. If the specified directory already exists,
    the export operation will fail.
  + **Overwrite**: If the specified export directory already exists, the directory will be overwritten.

  Defaults to **ErrorIfExists**.

* `with_column_header` - (Optional, Bool, NonUpdatable) Specifies whether to export column names when exporting data.  
  Defaults to **false**.

* `limit_num` - (Optional, Int, NonUpdatable) Specifies the number of data records to be exported.  
  The default value is **0**, meaning that all data records are exported.

* `encoding_type` - (Optional, String, NonUpdatable) Specifies the encoding format of the exported data.  
  The valid values are as follows:
  + **utf-8**
  + **gb2312**
  + **gbk**

  Defaults to **utf-8**.
  
* `quote_char` - (Optional, String, NonUpdatable) Specifies the custom quote character.  
  Only one character is supported.  
  This parameter is valid only when `data_type` is set to **csv**.  
  Defaults to `"`.

* `escape_char` - (Optional, String, NonUpdatable) Specifies the custom escape character.  
  Only one character is supported.  
  This parameter is valid only when `data_type` is set to **csv**.  
  Defaults to `\\`.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 10 minutes.
