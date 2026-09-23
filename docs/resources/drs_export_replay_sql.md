---
subcategory: "DRS"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_drs_export_replay_sql"
description: |-
  Use this resource to export the SQL report file of a DRS replay task within HuaweiCloud.
---

# huaweicloud_drs_export_replay_sql

Use this resource to export the SQL report file of a DRS replay task within HuaweiCloud.

-> This resource is a one-time action resource for exporting the SQL report file of a DRS replay task. Deleting this
   resource will not clear the corresponding request record, but will only remove the resource information from the
   tfstate file.

## Example Usage

### Export Abnormal SQL List

```hcl
variable "job_id" {}

resource "huaweicloud_drs_export_replay_sql" "test" {
  job_id    = var.job_id
  file_type = "abnormal_sql"
}
```

### Export Abnormal SQL Detail with Fields

```hcl
variable "job_id" {}

resource "huaweicloud_drs_export_replay_sql" "test" {
  job_id    = var.job_id
  file_type = "abnormal_sql_detail"
  field_names = [
    "id",
    "gmtCreate",
    "schemaName",
    "sqlStatement",
    "errorInfo",
  ]
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to export the replay SQL report.  
  If omitted, the provider-level region will be used. This parameter is non-updatable.

* `job_id` - (Required, String, NonUpdatable) Specifies the ID of the DRS replay job to export SQL report.

* `file_type` - (Required, String, NonUpdatable) Specifies the type of the SQL file to export.  
  The valid values are as follows:
  + **abnormal_sql**
  + **abnormal_sql_detail**
  + **slow_sql**
  + **slow_sql_detail**
  + **sql_result_inconsistent**

* `field_names` - (Optional, List, NonUpdatable) Specifies the list of field names to export.  
  Required when `file_type` is **abnormal_sql_detail** or **slow_sql_detail**.

* `start_time` - (Optional, String, NonUpdatable) Specifies the start time of the export range.

* `end_time` - (Optional, String, NonUpdatable) Specifies the end time of the export range.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID, which is the same as the job ID.
