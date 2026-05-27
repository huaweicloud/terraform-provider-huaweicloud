---
subcategory: "Data Replication Service (DRS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_drs_compare_line_detail"
description: |-
  Use this data source to get the line comparison detail information of a DRS job.
---

# huaweicloud_drs_compare_line_detail

Use this data source to get the line comparison detail information of a DRS job.

## Example Usage

```hcl
variable "job_id" {} 
variable "compare_job_id" {}

data "huaweicloud_drs_compare_line_detail" "test" {
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

* `status` - (Optional, Int) Specifies the comparison status for filtering.
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

* `type` - (Optional, String) Specifies the comparison type for filtering.
  The valid values are as follows:
  + **compare**: Query normal comparison items
  + **unCompare**: Query incomparable items

* `db_name` - (Optional, String) Specifies the source database name for filtering.

* `target_db_name` - (Optional, String) Specifies the target database name for filtering.

* `query_tb_name` - (Optional, String) Specifies the table name keyword for filtering.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `table_line_compare_result_infos` - The list of table line comparison result information.

  The [table_line_compare_result_infos](#table_line_compare_result_infos_struct) structure is documented below.

<a name="table_line_compare_result_infos_struct"></a>
The `table_line_compare_result_infos` block supports:

* `source_table_name` - The source table name.

* `source_row_num` - The number of rows in the source table.

* `target_table_name` - The target table name.

* `target_row_num` - The number of rows in the target table.

* `difference_row_num` - The difference in row count between source and target tables.

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

* `compare_line_config_filter` - The filter configuration for line comparison.
