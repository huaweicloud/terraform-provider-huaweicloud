---
subcategory: "Data Lake Insight (DLI)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dli_flinkjar_job"
description: ""
---

# huaweicloud_dli_flinkjar_job

Manages a flink job resource which type is `Flink Jar` within HuaweiCloud DLI.

## Example Usage

### Create a flink job

```hcl
variable "name" {}
variable "queue_name" {}
variable "jar_obs_path" {}
variable "entrypoint_args" {}

resource "huaweicloud_dli_package" "test" {
  group_name  = "jarPackage"
  type        = "jar"
  object_path = var.jar_obs_path
}

resource "huaweicloud_dli_flinkjar_job" "test" {
  name            = var.name
  queue_name      = var.queue_name
  entrypoint      = "${huaweicloud_dli_package.test.group_name}/${huaweicloud_dli_package.test.object_name}"
  entrypoint_args = var.entrypoint_args

  tags = {
    foo = "bar"
    key = "value"
  }
}

```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) The region in which to create the DLI flink job resource. If omitted, the
  provider-level region will be used. Changing this parameter will create a new resource.

* `name` - (Required, String) Specifies the name of the job. Length range: `1` to `57` characters.
 Which may consist of letters, digits, underscores (_) and hyphens (-).

* `description` - (Optional, String) Specifies job description. Length range: `1` to `512` characters.

* `queue_name` - (Optional, String) Specifies the name of DLI queue which this job run in. The type of queue
 must be `general`.

* `main_class` - (Optional, String) Specifies job entry class. Default main class is specified by the Manifest file
 of the application.

* `entrypoint` - (Optional, String) Specifies the JAR file where the job main class is located. It is the name of the
 package that has been uploaded to the DLI.

* `entrypoint_args` - (Optional, String) Specifies job entry arguments. Multiple arguments are separated by spaces.
  The arguments are keys followed by values. Keys have to start with '-' or '--'.

  Example: `--key1 value1 --key2 value2 -key3 value3`

* `dependency_jars` - (Optional, List) Specifies other dependency jars. It is the name of the package that
 has been uploaded to the DLI.

 Example: `["myGroup/test.jar" , "myGroup/test1.jar"]`.

* `dependency_files` - (Optional, List) Specifies dependency files. It is the name of the package that has been
 uploaded to the DLI.
  
  Example: `["myGroup/test.cvs" , "myGroup/test1.csv"]`.
  You can add the following content to the application to access the corresponding dependency file: In the command,
  fileName indicates the name of the file to be accessed, and ClassName indicates the name of the class that needs to
  access the file: ClassName.class.getClassLoader().getResource("userData/fileName").

* `feature` - (Optional, String) Specifies job feature. Type of the Flink image used by a job.
  + **basic**: indicates that the basic Flink image provided by DLI is used.
  + **custom**: indicates that the user-defined Flink image is used.

  The default value is **basic**.
  
* `flink_version` - (Optional, String) Specifies flink version. This parameter is valid only when feature is set
 to basic. You can use this parameter with the feature parameter to specify the version of the DLI basic Flink image
 used for job running. The options are as follows: `1.10` and `1.7`.

* `image` - (Optional, String) Specifies custom image. The format is Organization name/Image name:Image version.
  This parameter is valid only when feature is set to `custom`. You can use this parameter with the feature parameter
  to specify a user-defined Flink image for job running. For details about how to use custom images, see the
  Data Lake Insight User Guide <https://support.huaweicloud.com/en-us/usermanual-dli/dli_01_0494.html>.
  
* `cu_num` - (Optional, Int) Specifies number of CUs selected for a job. The default value is `2`.

* `parallel_num` - (Optional, Int) Specifies number of parallel for a job. The default value is `1`.

* `manager_cu_num` - (Optional, Int) Specifies number of CUs in the JobManager selected for a job.
 The default value is `1`.
  
* `tm_cu_num` - (Optional, Int) Specifies number of CUs for each TaskManager. The default value is `1`.
  
* `tm_slot_num` - (Optional, Int) Specifies number of slots in each TaskManager.
 The default value is `(parallel_num * tm_cu_num) / (cu_num - manager_cu_num)`.

* `smn_topic` - (Optional, String) Specifies SMN topic. If a job fails, the system will send a message to users
  subscribed to the SMN topic.

* `log_enabled` - (Optional, Bool) Specifies whether to enable the function of uploading job logs to users' OBS buckets.
 The default value is `false`.

* `obs_bucket` - (Optional, String) Specifies OBS path. OBS path where users are authorized to save the log.
  This parameter is valid only when `log_enabled` is set to `true`.

* `restart_when_exception` - (Optional, Bool) Specifies whether to enable the function of restart upon exceptions.
 The default value is `false`.
  
* `resume_checkpoint` - (Optional, Bool) Specifies whether the abnormal restart is recovered from the checkpoint.
  
* `resume_max_num` - (Optional, Int) Specifies maximum number of retry times upon exceptions. The unit is
 `times/hour`. Value range: `-1` or greater than `0`. The default value is `-1`, indicating that the number of times is
 unlimited.

* `checkpoint_path` - (Optional, String) Specifies storage address of the checkpoint in the JAR file of the user.
 The path must be unique.

* `runtime_config` - (Optional, Map) Specifies customizes optimization parameters when a Flink job is running.

* `tags` - (Optional, Map) Specifies the key/value pairs to associate with the resource.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The Job ID in Int format.

* `status` - The Job status.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 20 minutes.
* `update` - Default is 30 minutes.
* `delete` - Default is 40 minutes.

## Import

The job can be imported by `id`. For example,

```bash
terraform import huaweicloud_dli_flinkjar_job.test 12345
```
