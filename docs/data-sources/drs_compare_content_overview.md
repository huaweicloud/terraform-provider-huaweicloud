---
subcategory: "Data Replication Service (DRS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_drs_compare_content_overview"
description: |-
  Use this data source to get the content comparison overview information of a DRS job.
---

# huaweicloud_drs_compare_content_overview

Use this data source to get the content comparison overview information of a DRS job.

## Example Usage

```hcl
variable "job_id" {} 
variable "compare_job_id" {}

data "huaweicloud_drs_compare_content_overview" "test" {
  job_id         = var.job_id
  compare_job_id = var.compare_job_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `job_id` - (Required, String) Specifies the ID of the DRS job.

* `compare_job_id` - (Required, String) Specifies the ID of the compare job.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `content_compare_result_infos` - The list of content comparison result information.

  The [content_compare_result_infos](#content_compare_result_infos_struct) structure is documented below.

<a name="content_compare_result_infos_struct"></a>
The `content_compare_result_infos` block supports:

* `status` - The comparison status.
  The valid values are as follows:
  + **0**: Inconsistent
  + **2**: Consistent
  + **3**: Target database or table does not exist
  + **4**: Comparison failed
  + **5**: Comparing
  + **6**: Waiting for comparison
  + **7**: Task cancelled
  + **8**: Source database is empty
  + **9**: Target database is empty
  + **10**: Both source and target databases are empty
  + **11**: Source table does not exist
  + **12**: Target table does not exist
  + **13**: Both source and target tables do not exist
  + **14**: Source database connection failed
  + **15**: Target database connection failed
  + **16**: Source database SQL execution timeout
  + **17**: Target database SQL execution timeout
  + **18**: Source database SQL execution error
  + **19**: Target database SQL execution error
  + **20**: Both source and target databases do not exist
  + **21**: Source database does not exist
  + **22**: Target database does not exist
  + **23**: Hundreds of millions of rows, comparison not performed
  + **27**: Timeout

* `source_db` - The source database name.

* `target_db` - The target database name.

* `compare_num` - The total number of comparisons.

* `compare_end_num` - The number of completed comparisons.

* `data_inconsistent_num` - The number of inconsistent data records.

* `uncomparable_num` - The number of uncomparable records.
