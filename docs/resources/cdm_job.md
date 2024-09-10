---
subcategory: "Cloud Data Migration (CDM)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cdm_job"
description: ""
---

# huaweicloud_cdm_job

Manages CDM job resource within HuaweiCloud.

## Example Usage

### Create a cdm job

```hcl
variable "name" {}
variable "obs_link_name" {}
variable "obs_input_bucket" {}
variable "obs_output_bucket" {}
variable "obs_link_name" {}

resource "huaweicloud_obs_bucket" "input" {
  bucket        = "job-input"
  acl           = "private"
  force_destroy = true
}

resource "huaweicloud_obs_bucket" "output" {
  bucket        = "job-output"
  acl           = "private"
  force_destroy = true
}

resource "huaweicloud_cdm_job" "test" {
  name       = var.name
  job_type   = "NORMAL_JOB"
  cluster_id = huaweicloud_cdm_cluster.test.id

  source_connector = "obs-connector"
  source_link_name = var.obs_link_name
  source_job_config = {
    "bucketName"               = var.obs_input_bucket
    "inputDirectory"           = "/"
    "listTextFile"             = "false"
    "inputFormat"              = "BINARY_FILE"
    "fromCompression"          = "NONE"
    "fromFileOpType"           = "DO_NOTHING"
    "useMarkerFile"            = "false"
    "useTimeFilter"            = "false"
    "fileSeparator"            = "|"
    "filterType"               = "NONE"
    "useWildCard"              = "false"
    "decryption"               = "NONE"
    "nonexistentPathDisregard" = "false"
  }

  destination_connector = "obs-connector"
  destination_link_name = var.obs_link_name
  destination_job_config = {
    "bucketName"          = var.obs_output_bucket
    "outputDirectory"     = "/"
    "outputFormat"        = "BINARY_FILE"
    "validateMD5"         = "true"
    "recordMD5Result"     = "false"
    "duplicateFileOpType" = "REPLACE"
    "useCustomDirectory"  = "false"
    "encryption"          = "NONE"
    "copyContentType"     = "false"
    "shouldClearTable"    = "false"
  }

  config {
    retry_type                   = "NONE"
    scheduler_enabled            = false
    throttling_extractors_number = 4
    throttling_record_dirty_data = false
    throttling_max_error_records = 10
    throttling_loader_number     = 1
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) The region in which to create the job resource. If omitted, the
  provider-level region will be used. Changing this parameter will create a new resource.

* `name` - (Required, String) Specifies job name, which can contains of `1` to `240` characters, starting with a letter.
  Only letters, digits, hyphens (-), and underscores (_) are allowed.

* `cluster_id` - (Required, String, ForceNew) Specifies the ID of CDM cluster which this job run in.
 Changing this parameter will create a new resource.

* `job_type` - (Required, String, ForceNew) Specifies type of job. Changing this parameter will create a new resource.
 The options are as follows:

  + **NORMAL_JOB**: table/file migration.
  + **BATCH_JOB**: entire DB migration.
  + **SCENARIO_JOB**: scenario migration.

* `source_connector` - (Required, String, ForceNew) Specifies the connector name of source link.
 Changing this parameter will create a new resource.

* `source_link_name` - (Required, String, ForceNew) Specifies the source link name.
 Changing this parameter will create a new resource.

* `source_job_config` - (Required, Map) Specifies the source job configuration parameters. Each type of the data source
 to be connected has different configuration parameters, please refer to the document link below.

  + **From a Relational Database**: [configuration detail](https://support.huaweicloud.com/intl/en-us/api-cdm/cdm_02_0044.html)
  + **From OBS**: [configuration detail](https://support.huaweicloud.com/intl/en-us/api-cdm/cdm_02_0045.html)
  + **From HDFS**: [configuration detail](https://support.huaweicloud.com/intl/en-us/api-cdm/cdm_02_0046.html)
  + **From Hive**: [configuration detail](https://support.huaweicloud.com/intl/en-us/api-cdm/cdm_02_0047.html)
  + **From HBase/CloudTable**: [configuration detail](https://support.huaweicloud.com/intl/en-us/api-cdm/cdm_02_0048.html)
  + **From FTP/SFTP**: [configuration detail](https://support.huaweicloud.com/intl/en-us/api-cdm/cdm_02_0049.html)
  + **From HTTP/HTTPS**: [configuration detail](https://support.huaweicloud.com/intl/en-us/api-cdm/cdm_02_0086.html)
  + **From MongoDB/DDS**: [configuration detail](https://support.huaweicloud.com/intl/en-us/api-cdm/cdm_02_0071.html)
  + **From Redis**: [configuration detail](https://support.huaweicloud.com/intl/en-us/api-cdm/cdm_02_0050.html)
  + **From DIS**: [configuration detail](https://support.huaweicloud.com/intl/en-us/api-cdm/cdm_02_0078.html)
  + **From Kafka**: [configuration detail](https://support.huaweicloud.com/intl/en-us/api-cdm/cdm_02_0079.html)
  + **From Elasticsearch/Cloud Search Service**: [configuration detail](https://support.huaweicloud.com/intl/en-us/api-cdm/cdm_02_0067.html)

-> Please remove the `fromJobConfig.` in the parameter key listed in the document.

* `destination_connector` - (Required, String, ForceNew) Specifies the connector name of destination link.
 Changing this parameter will create a new resource.

* `destination_link_name` - (Required, String, ForceNew) Specifies the destination link name.
 Changing this parameter will create a new resource.

* `destination_job_config` - (Required, Map) Specifies the destination job configuration parameters. Each type of the
 data source to be connected has different configuration parameters, please refer to the document link below.

  + **To a Relational Database**: [configuration detail](https://support.huaweicloud.com/intl/en-us/api-cdm/cdm_02_0052.html)
  + **To OBS**: [configuration detail](https://support.huaweicloud.com/intl/en-us/api-cdm/cdm_02_0053.html)
  + **To HDFS**: [configuration detail](https://support.huaweicloud.com/intl/en-us/api-cdm/cdm_02_0054.html)
  + **To Hive**: [configuration detail](https://support.huaweicloud.com/intl/en-us/api-cdm/cdm_02_0055.html)
  + **To HBase/CloudTable**: [configuration detail](https://support.huaweicloud.com/intl/en-us/api-cdm/cdm_02_0056.html)
  + **To DDS**: [configuration detail](https://support.huaweicloud.com/intl/en-us/api-cdm/cdm_02_0084.html)
  + **To DLI**: [configuration detail](https://support.huaweicloud.com/intl/en-us/api-cdm/cdm_02_0080.html)
  + **To DIS**: [configuration detail](https://support.huaweicloud.com/intl/en-us/api-cdm/cdm_02_0088.html)
  + **To Elasticsearch/Cloud Search Service**: [configuration detail](https://support.huaweicloud.com/intl/en-us/api-cdm/cdm_02_0077.html)

-> Please remove the `toJobConfig.` in the parameter key listed in the document.

* `config` - (Optional, List) Specifies the job configuration. Structure is documented below.

The `config` block supports:

* `group_name` - (Optional, String) Specifies group to which a job belongs. The default group is `DEFAULT`.

* `retry_type` - (Optional, String) Specifies whether to automatically retry if a job fails to be executed.
 The options are as follows:
  + **NONE**: Do not retry.
  + **RETRY_TRIPLE**: Retry three times.

  Default value is `NONE`.

* `throttling_extractors_number` - (Optional, Int) Specifies maximum number of concurrent extraction jobs.

* `throttling_loader_number` - (Optional, Int) Specifies maximum number of loading jobs. This parameter is available
 only when HBase or Hive serves as the destination data source.

* `throttling_record_dirty_data` - (Optional, Bool) Specifies whether to write dirty data.

* `throttling_dirty_write_to_link` - (Optional, String) Specifies the link name to which dirty data is written to.
 The Dirty data can be written only to `OBS` or `HDFS`.

* `throttling_dirty_write_to_bucket` - (Optional, String) Specifies name of the OBS bucket to which dirty data is
 written. This parameter is valid only when dirty data is written to `OBS`.

* `throttling_dirty_write_to_directory` - (Optional, String) Specifies the directory in the OBS bucket or HDFS which
 dirty data is written to. For example, `/data/dirtydata/`.

* `throttling_max_error_records` - (Optional, Int) Specifies maximum number of error records in a single
 shard. When the number of error records of a map exceeds the upper limit, the task automatically ends.

* `scheduler_enabled` - (Optional, Bool) Specifies whether to enable a scheduled task.  Default value is `false`.

* `scheduler_cycle_type` - (Optional, String) Specifies cycle type of a scheduled task. The options are as follows:
 `minute`, `hour`, `day`, `week`, `month`.

* `scheduler_cycle` - (Optional, Int) Specifies cycle of a scheduled task. If `scheduler_cycle_type` is set to minute
 and `scheduler_cycle` is set to 10, the scheduled task is executed every 10 minutes.

* `scheduler_run_at` - (Optional, String) Specifies time when a scheduled task is triggered in a cycle. This parameter
 is valid only when `scheduler_cycle_type` is set to `hour`, `week`, or `month`.
  + If `scheduler_cycle_type` is set to month, cycle is set to 1, and runAt is set to 15, the scheduled task is executed
    on the 15th day of each month. You can set runAt to multiple values and separate the values with commas (,).
    For example, if runAt is set to 1,2,3,4,5, the scheduled task is executed on the first day, second day, third day,
    fourth day, and fifth day of each month.
  + If `scheduler_cycle_type` is set to week and runAt is set to mon,tue,wed,thu,fri, the scheduled task is executed on
    Monday to Friday.
  + If `scheduler_cycle_type` is set to hour and runAt is set to 27,57, the scheduled task is executed at the 27th and
    57th minute in the cycle.

* `scheduler_start_date` - (Optional, String) Specifies start time of a scheduled task.
 For example, `2018-01-24 19:56:19`

* `scheduler_stop_date` - (Optional, String) Specifies End time of a scheduled task. For example, `2018-01-27 23:59:00`.
 If you do not set the end time, the scheduled task is always executed and will never stop.

* `scheduler_disposable_type` - (Optional, String) Specifies whether to delete a job after the job is executed.
 The options are as follows:
  + **NONE**: The job will not be deleted after it is executed.
  + **DELETE_AFTER_SUCCEED**: The job will be deleted only after it is successfully executed. It is applicable to
    massive one-time jobs.
  + **DELETE**: The job will be deleted after it is executed, regardless of the execution result.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID in format of **cluster_id/job_name**. It is composed of the ID of CDM cluster which this
 job run in and the name of job, separated by a slash.

* `status` - Job status. The options are as follows:

  + **BOOTING**: The job is starting.
  + **FAILURE_ON_SUBMIT**: The job fails to be submitted.
  + **RUNNING**: The job is running.
  + **SUCCEEDED**: The job is executed successfully.
  + **FAILED**: The job failed.
  + **UNKNOWN**: The job status is unknown.
  + **NEVER_EXECUTED**: The job has not been executed.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 20 minutes.

* `update` - Default is 20 minutes.

* `delete` - Default is 20 minutes.

## Import

Jobs can be imported by `id`. It is composed of the ID of CDM cluster which this job run in and the name of job,
 separated by a slash. For example,

```bash
terraform import huaweicloud_cdm_job.test b11b407c-e604-4e8d-8bc4-92398320b847/jobName
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response, security or some other reason. The missing attributes include: `source_job_config` and `destination_job_config`.
 It is generally recommended running `terraform plan` after importing a cluster.
 You can then decide if changes should be applied to the cluster, or the resource definition
should be updated to align with the cluster. Also you can ignore changes as below.

```hcl
resource "huaweicloud_cdm_cluster" "test" {
    ...

  lifecycle {
    ignore_changes = [
      source_job_config, destination_job_config,
    ]
  }
}
```
