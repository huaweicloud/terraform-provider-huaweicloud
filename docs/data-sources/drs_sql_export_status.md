---
subcategory: "DRS"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_drs_sql_export_status"
description: |-
  Use this data source to query the SQL file export status of a DRS replay task within HuaweiCloud.
---

# huaweicloud_drs_sql_export_status

Use this data source to query the SQL file export status of a DRS replay task within HuaweiCloud.

## Example Usage

### Basic Usage

```hcl
variable "job_id" {}

data "huaweicloud_drs_sql_export_status" "test" {
  job_id    = var.job_id
  file_type = "abnormal_sql"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the SQL export status.  
  If omitted, the provider-level region will be used.

* `job_id` - (Required, String) Specifies the ID of the DRS replay job to query SQL export status.

* `file_type` - (Required, String) Specifies the type of the SQL file to query export status.  
  The valid values are as follows:
  + **abnormal_sql**
  + **abnormal_sql_detail**
  + **slow_sql**
  + **slow_sql_detail**

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `export_status` - The export status of the SQL file.

* `failed_reason` - The reason of the export failure.

* `total_count` - The total amount of exported data.

* `current_count` - The amount of currently processed data.

* `progress_percentage` - The progress percentage of the export task.

* `uploaded_file_names` - The names of the files uploaded to OBS.
