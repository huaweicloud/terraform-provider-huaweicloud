---
subcategory: "Data Lake Insight (DLI)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dli_sql_job"
description: ""
---

# huaweicloud_dli_sql_job

Manages DLI SQL job resource within HuaweiCloud

## Example Usage

### Create a Sql job

```hcl
variable "database_name" {}
variable "queue_name" {}
variable "sql" {}

resource "huaweicloud_dli_sql_job" "test" {
  sql           = var.sql
  database_name = var.database_name
  queue_name    = var.queue_name
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the DLI table resource. If omitted,
  the provider-level region will be used. Changing this parameter will create a new resource.

* `sql` - (Required, String, ForceNew) Specifies SQL statement that you want to execute.
  Changing this parameter will create a new resource.

* `database_name` - (Optional, String, ForceNew) Specifies the database where the SQL is executed. This argument does
 not need to be configured during database creation. Changing this parameter will create a new resource.

* `queue_name` - (Optional, String, ForceNew) Specifies queue which this job to be submitted belongs.
 Changing this parameter will create a new resource.

* `tags` - (Optional, Map, ForceNew) Specifies label of a Job. Changing this parameter will create a new resource.

* `conf` - (Optional, List, ForceNew) Specifies the configuration parameters for the SQL job. Changing this parameter
 will create a new resource. Structure is documented below.

 The `conf` block supports:

   * `spark_sql_max_records_per_file` - (Optional, Int, ForceNew) Maximum number of records to be written
    into a single file. If the value is zero or negative, there is no limit. Default value is `0`.
     Changing this parameter will create a new resource.

   * `spark_sql_auto_broadcast_join_threshold` - (Optional, Int, ForceNew) Maximum size of the table that
    displays all working nodes when a connection is executed. You can set this parameter to -1 to disable the display.
    Default value is `209715200`. Changing this parameter will create a new resource.

   -> Currently, only the configuration unit metastore table that runs the ANALYZE TABLE COMPUTE statistics noscan
    command and the file-based data source table that directly calculates statistics based on data files are supported.
     Changing this parameter will create a new resource.

   * `spark_sql_shuffle_partitions` - (Optional, Int, ForceNew) Default number of partitions used to filter
    data for join or aggregation. Default value is `4096`. Changing this parameter will create a new resource.

   * `spark_sql_dynamic_partition_overwrite_enabled` - (Optional, Bool, ForceNew) In dynamic mode, Spark does not delete
    the previous partitions and only overwrites the partitions without data during execution. Default value is `false`.
    Changing this parameter will create a new resource.
   * `spark_sql_files_max_partition_bytes` - (Optional, Int, ForceNew) Maximum number of bytes to be packed into a
    single partition when a file is read. Default value is `134217728`. Changing this parameter will create a new
     resource.

   * `spark_sql_bad_records_path` - (Optional, String, ForceNew) Path of bad records. Changing this parameter will create
    a new resource.

   * `dli_sql_sqlasync_enabled` - (Optional, Bool, ForceNew) Specifies whether DDL and DCL statements are executed
    asynchronously. The value true indicates that asynchronous execution is enabled. Default value is `false`.
     Changing this parameter will create a new resource.

   * `dli_sql_job_timeout` - (Optional, Int, ForceNew) Sets the job running timeout interval. If the timeout interval
    expires, the job is canceled. Unit: `ms`. Changing this parameter will create a new resource.
  
## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - Indicates a resource ID in UUID format.

* `owner` - User who submits a job.

* `job_type` - Type of a job, Includes **DDL**, **DCL**, **IMPORT**, **EXPORT**, **QUERY**, **INSERT**,
 **DATA_MIGRATION**, **UPDATE**, **DELETE**, **RESTART_QUEUE** and **SCALE_QUEUE**.

* `status` - Status of a job, including **RUNNING**, **SCALING**, **LAUNCHING**, **FINISHED**, **FAILED**,
  and **CANCELLED.**

* `start_time` - Time when a job is started, in RFC-3339 format. e.g. `2019-10-12T07:20:50.52Z`

* `duration` - Job running duration (unit: millisecond).

* `schema` - When the statement type is DDL, the column name and type of DDL are displayed.

* `rows` - When the statement type is DDL, results of the DDL are displayed.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 20 minutes.

* `delete` - Default is 45 minutes.

## Import

DLI SQL job can be imported by `id`. For example,

```bash
terraform import huaweicloud_dli_sql_job.example 7f803d70-c533-469f-8431-e378f3e97123
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response, security or some other reason. The missing attributes include: `conf`, `rows` and `schema`.
It is generally recommended running `terraform plan` after importing a resource. You can then decide if changes should
be applied to the resource, or the resource definition should be updated to align with the resource. Also you can
ignore changes as below.

```hcl
resource "huaweicloud_dli_sql_job" "test" {
    ...

  lifecycle {
    ignore_changes = [
      conf, rows, schema
    ]
  }
}
```
