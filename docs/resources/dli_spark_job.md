---
subcategory: "Data Lake Insight (DLI)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dli_spark_job"
description: ""
---

# huaweicloud_dli_spark_job

Manages spark job resource of DLI within HuaweiCloud

## Example Usage

### Submit a new spark job with jar packages

```hcl
variables "queue_name" {}
variables "job_name" {}

resource "huaweicloud_dli_spark_job" "default" {
  queue_name    = var.queue_name
  name          = var.job_name
  app_name      = "driver_package/driver_behavior.jar"
  main_class    = "driver_behavior"
  specification = "B"
  max_retries   = 20
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to submit a spark job.
  If omitted, the provider-level region will be used.
  Changing this parameter will submit a new spark job.

* `queue_name` - (Required, String, ForceNew) Specifies the DLI queue name.
  Changing this parameter will submit a new spark job.

* `name` - (Required, String, ForceNew) Specifies the spark job name.
  The value contains a maximum of 128 characters.
  Changing this parameter will submit a new spark job.

* `app_name` - (Required, String, ForceNew) Specifies the name of the package that is of the JAR or python file type and
  has been uploaded to the DLI resource management system.
  The OBS paths are allowed, for example, `obs://<bucket name>/<package name>`.
  Changing this parameter will submit a new spark job.

* `app_parameters` - (Optional, String, ForceNew) Specifies the input parameters of the main class.
  Changing this parameter will submit a new spark job.

* `main_class` - (Optional, String, ForceNew) Specifies the main class of the spark job.
  Required if the `app_name` is the JAR type.
  Changing this parameter will submit a new spark job.

* `jars` - (Optional, List, ForceNew) Specifies a list of the jar package name which has been uploaded to the DLI
  resource management system. The OBS paths are allowed, for example, `obs://<bucket name>/<package name>`.
  Changing this parameter will submit a new spark job.

* `python_files` - (Optional, List, ForceNew) Specifies a list of the python file name which has been uploaded to the
  DLI resource management system. The OBS paths are allowed, for example, `obs://<bucket name>/<python file name>`.
  Changing this parameter will submit a new spark job.

* `files` - (Optional, List, ForceNew) Specifies a list of the other dependencies name which has been uploaded to the
  DLI resource management system. The OBS paths are allowed, for example, `obs://<bucket name>/<dependent files>`.
  Changing this parameter will submit a new spark job.

* `dependent_packages` - (Optional, List, ForceNew) Specifies a list of package resource objects.
  The object structure is documented below.
  Changing this parameter will submit a new spark job.

* `configurations` - (Optional, Map, ForceNew) Specifies the configuration items of the DLI spark.
  Please following the document of Spark [configurations](https://spark.apache.org/docs/latest/configuration.html) for
  this argument. If you want to enable the `access metadata` of DLI spark in HuaweiCloud, please set
  `spark.dli.metaAccess.enable` to `true`. Changing this parameter will submit a new spark job.

* `modules` - (Optional, List, ForceNew) Specifies a list of modules that depend on system resources.
  The dependent modules and corresponding services are as follows.
  Changing this parameter will submit a new spark job.
  + **sys.datasource.hbase**: CloudTable/MRS HBase
  + **sys.datasource.opentsdb**: CloudTable/MRS OpenTSDB
  + **sys.datasource.rds**: RDS MySQL
  + **sys.datasource.css**: CSS

* `specification` - (Optional, String, ForceNew) Specifies the compute resource type for spark application.
  The available types and related specifications are as follows, default to minimum configuration (type **A**).
  Changing this parameter will submit a new spark job.

  | type | resource | driver cores | excutor cores | driver memory | executor memory | num executor |
  | ---- | ---- | ---- | ---- | ---- | ---- | ---- |
  | A | 8 vCPUs, 32-GB memory | 2 | 1 | 7G | 4G | 6 |
  | B | 16 vCPUs, 64-GB memory | 2 | 2 | 7G | 8G | 7 |
  | C | 32 vCPUs, 128-GB memory | 4 | 2 | 12G | 8G | 14 |

* `executor_memory` - (Optional, String, ForceNew) Specifies the executor memory of the spark application.
  application. The default value of this value corresponds to the configuration of the selected `specification`.
  If you set this value instead of the default value, `specification` will be invalid.
  Changing this parameter will submit a new spark job.

  ->**NOTE:** The unit must be provided, such as **GB** or **MB**.

* `executor_cores` - (Optional, Int, ForceNew) Specifies the number of CPU cores of each executor in the Spark
  application. The default value of this value corresponds to the configuration of the selected `specification`.
  If you set this value instead of the default value, `specification` will be invalid.
  Changing this parameter will submit a new spark job.

* `executors` - (Optional, Int, ForceNew) Specifies the number of executors in a spark application.
  The default value of this value corresponds to the configuration of the selected `specification`.
  If you set this value instead of the default value, `specification` will be invalid.
  Changing this parameter will submit a new spark job.

* `driver_memory` - (Optional, String, ForceNew) Specifies the driver memory of the spark application.
  The default value of this value corresponds to the configuration of the selected `specification`.
  If you set this value instead of the default value, `specification` will be invalid.
  Changing this parameter will submit a new spark job.

* `driver_cores` - (Optional, Int, ForceNew) Specifies the number of CPU cores of the Spark application driver.
  The default value of this value corresponds to the configuration of the selected `specification`.
  If you set this value instead of the default value, `specification` will be invalid.
  Changing this parameter will submit a new spark job.

* `max_retries` - (Optional, Int, ForceNew) Specifies the maximum retry times.
  The default value of this value corresponds to the configuration of the selected `specification`.
  If you set this value instead of the default value, `specification` will be invalid.
  Changing this parameter will submit a new spark job.

The `dependent_packages` block supports:

* `group_name` - (Required, String, ForceNew) Specifies the user group name.  
  Only letters, digits, dots (.), hyphens (-) and underscores (_) are allowed.  
  Changing this parameter will submit a new spark job.

* `packages` - (Required, List, ForceNew) Specifies the user group resource for details.
  Changing this parameter will submit a new spark job.
  The [object](#dependent_packages_packages) structure is documented below.

<a name="dependent_packages_packages"></a>
The `packages` block supports:

* `type` - (Required, String, ForceNew) Specifies the resource type of the package.
  Changing this parameter will submit a new spark job.

* `package_name` - (Required, String, ForceNew) Specifies the resource name of the package.
  Changing this parameter will submit a new spark job.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the spark job.

* `created_at` - Time of the DLI spark job submit.

* `owner` - The owner of the spark job.
