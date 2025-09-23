---
subcategory: "Cloud Data Migration (CDM)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cdm_job_execution_records"
description: |-
  Use this data source to a list of CDM job execution records.
---

# huaweicloud_cdm_job_execution_records

Use this data source to a list of CDM job execution records.

## Example Usage

```hcl
variable "cluster_id" {}
variable "job_name" {}

data "huaweicloud_cdm_job_execution_records" "test" {
  cluster_id = var.cluster_id
  job_name   = var.job_name
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `cluster_id` - (Required, String) Specifies the cluster ID.

* `job_name` - (Required, String) Specifies the job name.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `records` - Indicates the records.
  The [records](#attrblock--records) structure is documented below.

<a name="attrblock--records"></a>
The `records` block supports:

* `counters` - Indicates the job running result statistics. Only return when status is SUCCEEDED.
  The [counters](#attrblock--records--counters) structure is documented below.

* `creation_date` - Indicates the creation time.

* `creation_user` - Indicates the user who created the job.

* `delete_rows` - Indicates the number of deleted rows.

* `error_details` - Indicates the error details.

* `error_summary` - Indicates the error summary.

* `execute_date` - Indicates the execution time.

* `external_id` - Indicates the job ID.

* `is_delete_job` - Indicates whether the job to be deleted after it is executed.

* `is_execute_auto` - Indicates whether the job executed as scheduled.

* `is_incrementing` - Indicates whether the job migrates incremental data.

* `is_stoping_increment` - Indicates whether incremental data migration stopped.

* `last_udpate_user` - Indicates the user who last updated the job status.

* `last_update_date` - Indicates the time when the job was last updated.

* `progress` - Indicates the Job progress.
  + If a job fails, the value is **-1**.
  + Otherwise, the value ranges from **0** to **100**.

* `status` - Indicates the Job status.
  Value can be as follows:
  + **BOOTING**: The job is starting.
  + **FAILURE_ON_SUBMIT**: The job failed to be submitted.
  + **RUNNING**: The job is running.
  + **SUCCEEDED**: The job was successfully executed.
  + **FAILED**: The job execution failed.
  + **UNKNOWN**: The job status is unknown.
  + **NEVER_EXECUTED**: The job was not executed.

* `submission_id` - Indicates the job submission ID.

* `update_rows` - Indicates the number of updated rows.

* `write_rows` - Indicates the number of write rows.

<a name="attrblock--records--counters"></a>
The `counters` block supports:

* `bytes_read` - Indicates the number of bytes that are read.

* `bytes_written` - Indicates the number of bytes that are written.

* `file_skipped` - Indicates the number of files that are skipped.

* `files_read` - Indicates the number of files that are read.

* `files_written` - Indicates the number of files that are written.

* `rows_read` - Indicates the number of rows that are read.

* `rows_written` - Indicates the number of rows that are written.

* `rows_written_skipped` - Indicates the number of rows that are skipped.

* `total_files` - Indicates the total number of files.

* `total_size` - Indicates the total number of bytes.
