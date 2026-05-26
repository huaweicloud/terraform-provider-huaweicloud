---
subcategory: "GaussDB"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_gaussdb_wdr_snapshot_collection_results"
description: |-
  Use this data source to query the WDR snapshot collection results of a GaussDB instance within HuaweiCloud.
---

# huaweicloud_gaussdb_wdr_snapshot_collection_results

Use this data source to query the WDR snapshot collection results of a GaussDB instance within HuaweiCloud.

## Example Usage

```hcl
variable "instance_id" {}

data "huaweicloud_gaussdb_wdr_snapshot_collection_results" "test" {
  instance_id = var.instance_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the ID of the GaussDB instance.

* `start_time` - (Optional, String) Specifies the start time for querying snapshots. The format is **yyyy-mm-ddThh:mm:ssZ**,
  where T indicates the start of a certain time, and Z indicates the time zone offset.

* `end_time` - (Optional, String) Specifies the end time for querying snapshots. The format is **yyyy-mm-ddThh:mm:ssZ**,
  where T indicates the start of a certain time, and Z indicates the time zone offset.

* `job_id` - (Optional, String) Specifies the task ID to query collection results for a specific task.

* `status` - (Optional, String) Specifies the task collection status.
  The valid values are as follows:
  + **EXPORTING**
  + **SUCCESS**
  + **FAILED**

* `wdr_type` - (Optional, String) Specifies the collection type.
  The valid values are as follows:
  + **cluster**
  + **component**
  + **pdb**

* `job_start_time` - (Optional, String) Specifies the start time for querying task creation time.
  The format is `yyyy-mm-ddThh:mm:ssZ`. For example, Beijing time offset is `+0800`.

* `job_end_time` - (Optional, String) Specifies the end time for querying task creation time.
  The format is `yyyy-mm-ddThh:mm:ssZ`. For example, Beijing time offset is `+0800`.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `wdr_snapshots` - The list of WDR snapshot collection results.
  The [wdr_snapshots](#gaussdb_wdr_snapshot_collection_results_wdr_snapshots_attr) structure is documented below.

<a name="gaussdb_wdr_snapshot_collection_results_wdr_snapshots_attr"></a>
The `wdr_snapshots` block supports:

* `job_id` - The task ID.

* `file_size` - The file size in KB.

* `wdr_type` - The collection type.

* `start_time` - The start snapshot time when the collection was initiated.

* `end_time` - The end snapshot time when the collection was initiated.

* `job_create_time` - The creation time of the WDR report generation task.

* `start_snapshot_id` - The first comparison snapshot ID used to generate the WDR report.

* `end_snapshot_id` - The second comparison snapshot ID used to generate the WDR report.

* `download_url` - The report download link, valid for 30 minutes.

* `status` - The collection status.

* `notes` - The remarks. When the collection type is component level, the content includes the collected component IDs.

* `error_msg` - The error message for operations analysis.

* `file_name` - The temporary file name of the WDR report.

* `file_path` - The temporary file storage path of the WDR report.

* `obs_bucket` - The OBS bucket information for storing the WDR report temporary file.
  The [obs_bucket](#gaussdb_wdr_snapshot_collection_results_wdr_snapshots_obs_bucket_attr) structure is documented below.

<a name="gaussdb_wdr_snapshot_collection_results_wdr_snapshots_obs_bucket_attr"></a>
The `obs_bucket` block supports:

* `name` - The OBS bucket name.

* `type` - The OBS bucket type.

* `url` - The OBS service access address.

* `port` - The OBS service port number.

* `domain_id` - The final tenant ID.
