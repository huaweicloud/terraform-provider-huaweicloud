---
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_mrs_job"
sidebar_current: "docs-huaweicloud-resource-mrs-job"
description: |-
  Manages resource job within HuaweiCloud MRS.
---

# huaweicloud\_mrs\_job

Manages resource job within HuaweiCloud MRS.
This is an alternative to `huaweicloud_mrs_job_v1`

## Example Usage

```hcl

resource "huaweicloud_mrs_job" "job1" {
  job_type   = 1
  job_name   = "test_mapreduce_job1"
  cluster_id = "ef43d2ff-1ecf-4f13-bd0c-0004c429a058"
  jar_path   = "s3a://wordcount/program/hadoop-mapreduce-examples-2.7.5.jar"
  input      = "s3a://wordcount/input/"
  output     = "s3a://wordcount/output/"
  job_log    = "s3a://wordcount/log/"
  arguments  = "wordcount"
}

```

## Argument Reference

The following arguments are supported:

* `job_type` - (Required) Job type 1: MapReduce 2: Spark 3: Hive Script 4: HiveQL
    (not supported currently) 5: DistCp, importing and exporting data.  6: Spark
    Script 7: Spark SQL, submitting Spark SQL statements. (not supported in this
    APIcurrently) NOTE: Spark and Hive jobs can be added to only clusters including
    Spark and Hive components.

* `job_name` - (Required) Job name Contains only 1 to 64 letters, digits, hyphens
    (-), and underscores (_). NOTE: Identical job names are allowed but not recommended.

* `cluster_id` - (Required) Cluster ID

* `jar_path` - (Required) Path of the .jar package or .sql file for program
    execution The parameter must meet the following requirements: Contains a maximum
    of 1023 characters, excluding special characters such as ;|&><'$. The address
    cannot be empty or full of spaces. Starts with / or s3a://. Spark Script must
    end with .sql; while MapReduce and Spark Jar must end with .jar. sql and jar
    are case-insensitive.

* `arguments` - (Optional) Key parameter for program execution. The parameter
    is specified by the function of the user's program. MRS is only responsible
    for loading the parameter. The parameter contains a maximum of 2047 characters,
    excluding special characters such as ;|&>'<$, and can be empty.

* `input` - (Optional) Path for inputting data, which must start with / or s3a://.
    A correct OBS path is required. The parameter contains a maximum of 1023 characters,
    excluding special characters such as ;|&>'<$, and can be empty.

* `output` - (Optional) Path for outputting data, which must start with / or
    s3a://. A correct OBS path is required. If the path does not exist, the system
    automatically creates it. The parameter contains a maximum of 1023 characters,
    excluding special characters such as ;|&>'<$, and can be empty.

* `job_log` - (Optional) Path for storing job logs that record job running status.
    This path must start with / or s3a://. A correct OBS path is required. The parameter
    contains a maximum of 1023 characters, excluding special characters such as
    ;|&>'<$, and can be empty.

* `hive_script_path` - (Optional) SQL program path This parameter is needed
    by Spark Script and Hive Script jobs only and must meet the following requirements:
    Contains a maximum of 1023 characters, excluding special characters such as
    ;|&><'$. The address cannot be empty or full of spaces. Starts with / or s3a://.
    Ends with .sql. sql is case-insensitive.

* `is_protected` - (Optional) Whether a job is protected true false The current
    version does not support this function.

* `is_public` - (Optional) Whether a job is public true false The current version
    does not support this function.

## Attributes Reference

The following attributes are exported:

* `job_type` - See Argument Reference above.
* `job_name` - See Argument Reference above.
* `cluster_id` - See Argument Reference above.
* `jar_path` - See Argument Reference above.
* `arguments` - See Argument Reference above.
* `input` - See Argument Reference above.
* `output` - See Argument Reference above.
* `job_log` - See Argument Reference above.
* `hive_script_path` - See Argument Reference above.
* `is_protected` - See Argument Reference above.
* `is_public` - See Argument Reference above.
