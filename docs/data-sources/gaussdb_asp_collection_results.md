---
subcategory: "GaussDB"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_gaussdb_asp_collection_results"
description: |-
  Use this data source to query ASP collection results of a specified GaussDB instance within HuaweiCloud.
---

# huaweicloud_gaussdb_asp_collection_results

Use this data source to query ASP collection results of a specified GaussDB instance within HuaweiCloud.

## Example Usage

```hcl
variable "instance_id" {}

data "huaweicloud_gaussdb_asp_collection_results" "test" {
  instance_id = var.instance_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the ASP collection results. If omitted,
  the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the ID of the GaussDB instance to query ASP collection results for.

* `start_time` - (Optional, String) Specifies the start time for querying ASP collection results. The format is
  **yyyy-mm-ddThh:mm:ssZ**, where T indicates the start of a certain time, and Z indicates the time zone offset.  

* `end_time` - (Optional, String) Specifies the end time for querying ASP collection results. The format is
  **yyyy-mm-ddThh:mm:ssZ**, where T indicates the start of a certain time, and Z indicates the time zone offset.  

* `job_id` - (Optional, String) Specifies the ASP collection task ID.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `asp` - The list of ASP collection results.  
  The [asp](#asp_attr) structure is documented below.

<a name="asp_attr"></a>
The `asp` block supports:

* `job_id` - The task ID of the ASP collection.

* `file_size` - The file size of the ASP report, in KB.

* `file_path` - The file path of the ASP report.

* `file_name` - The file name of the ASP report.

* `start_time` - The start time of the ASP collection. The format is **yyyy-mm-ddThh:mm:ssZ**.

* `end_time` - The end time of the ASP collection. The format is **yyyy-mm-ddThh:mm:ssZ**.

* `download_url` - The download link for the ASP report. The link is valid for 30 minutes.

* `status` - The status of the ASP collection. The valid values are as follows:
  + **SUCCESS**: The collection is successful.
  + **FAILED**: The collection is failed.
  + **EXPORTING**: The collection is in progress.

* `obs_bucket` - The obs bucket of ASP collection. The [obs_bucket](#obs_bucket_attr) structure is documented below.

<a name="obs_bucket_attr"></a>
The `obs_bucket` block supports:

* `name` - The name of the OBS bucket.

* `type` - The type of the OBS bucket.

* `url` - The url of the OBS bucket.

* `port` - The port of the OBS bucket.

* `domain_id` - The domain ID.
