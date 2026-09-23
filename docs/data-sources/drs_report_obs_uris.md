---
subcategory: "Data Replication Service (DRS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_drs_report_obs_uris"
description: |-
  Use this data source to query the OBS temporary download links of the DRS replay reports within HuaweiCloud.
---

# huaweicloud_drs_report_obs_uris

Use this data source to query the OBS temporary download links of the DRS replay reports within HuaweiCloud.

## Example Usage

```hcl
variable "job_id" {}

data "huaweicloud_drs_report_obs_uris" "test" {
  job_id      = var.job_id
  report_type = "abnormal_sql"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the OBS temporary download links.
  If omitted, the provider-level region will be used.

* `job_id` - (Required, String) Specifies the ID of the DRS job.

* `report_type` - (Required, String) Specifies the type of the exported SQL file.  
  The valid values are as follows:
  + **abnormal_sql**: The list of abnormal SQL statements.
  + **abnormal_sql_detail**: The details of abnormal SQL statements.
  + **slow_sql**: The list of slow SQL statements.
  + **slow_sql_detail**: The details of slow SQL statements.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `count` - The total number of downloadable files.

* `export_report_obs_files` - The list of downloadable files.

  The [export_report_obs_files](#export_report_obs_files_struct) structure is documented below.

<a name="export_report_obs_files_struct"></a>
The `export_report_obs_files` block supports:

* `file_name` - The name of the file.

* `obs_temp_uri` - The OBS temporary URI of the file.
