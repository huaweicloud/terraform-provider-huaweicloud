---
subcategory: "Data Lake Insight (DLI)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dli_spark_template"
description: ""
---

# huaweicloud_dli_spark_template

Manages a DLI Spark template resource within HuaweiCloud.  

## Example Usage

```hcl
  resource "huaweicloud_dli_spark_template" "test" {
    name        = "demo"
    description = "This is a demo"
    group       = "demo"

    body {
      queue_name    = "queue_demo"
      name          = "demo"
      app_name      = "jar_package/demo.jar"
      main_class    = "com.demo.main"
      specification = "B"
    }
  }
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `name` - (Required, String) The name of the spark template.

* `body` - (Required, List) The content of the spark template.
  The [body](#SparkTemplate_body) structure is documented below.

* `group` - (Optional, String) The group of the spark template.

* `description` - (Optional, String) The description of the spark template.

<a name="SparkTemplate_body"></a>
The `body` block supports:

* `queue_name` - (Optional, String) The DLI queue name.

* `name` - (Optional, String) The spark job name.

* `app_name` - (Optional, String) Name of the package that is of the JAR or pyFile type.  
  You can also specify an OBS path, for example, obs://Bucket name/Package name.

* `main_class` - (Optional, String) Java/Spark main class of the template.

* `app_parameters` - (Optional, List) Input parameters of the main class, that is application parameters.

* `specification` - (Optional, String) Compute resource type. Currently, resource types A, B, and C are available.  
  The available types and related specifications are as follows, default to minimum configuration (type **A**).

  | type | resource | driver cores | executor cores | driver memory | executor memory | num executor |
  | ---- | ---- | ---- | ---- | ---- | ---- | ---- |
  | A | 8 vCPUs, 32-GB memory | 2 | 1 | 7G | 4G | 6 |
  | B | 16 vCPUs, 64-GB memory | 2 | 2 | 7G | 8G | 7 |
  | C | 32 vCPUs, 128-GB memory | 4 | 2 | 12G | 8G | 14 |

* `jars` - (Optional, List) Name of the package that is of the JAR type and has been uploaded to the DLI
  resource management system. You can also specify an OBS path, for example, obs://Bucket name/Package name.

* `python_files` - (Optional, List) Name of the package that is of the PyFile type and has been uploaded to the DLI
  resource management system. You can also specify an OBS path, for example, obs://Bucket name/Package name.

* `files` - (Optional, List) Name of the package that is of the file type and has been uploaded to the
  DLI resource management system. You can also specify an OBS path, for example, obs://Bucket name/Package name.

* `modules` - (Optional, List) Name of the dependent system resource module.
  DLI provides dependencies for executing datasource jobs.
  The dependent modules and corresponding services are as follows.
    + **sys.datasource.hbase**: CloudTable/MRS HBase
    + **sys.datasource.opentsdb**: CloudTable/MRS OpenTSDB
    + **sys.datasource.rds**: RDS MySQL
    + **sys.datasource.css**: CSS

* `resources` - (Optional, List) The list of resource objects.
  The [resources](#SparkTemplate_Resources) structure is documented below.

* `dependent_packages` - (Optional, List) The list of package resource objects.  
  The [dependent_packages](#SparkTemplate_Dependent_packages) structure is documented below.

* `configurations` - (Optional, Map) The configuration items of the DLI spark.  
  For details, see [Spark configuration](https://spark.apache.org/docs/latest/configuration.html)
  If you want to enable the **access metadata** of DLI spark in HuaweiCloud, please set
  **spark.dli.metaAccess.enable** to **true**.

* `driver_memory` - (Optional, String) Driver memory of the Spark application, for example, 2 GB and 2048 MB.  
  This configuration item replaces the default parameter in **specification**.
  The unit must be provided. Otherwise, the startup fails.

* `driver_cores` - (Optional, Int) Number of CPU cores of the Spark application driver.  
  This configuration item replaces the default parameter in **specification**.

* `executor_memory` - (Optional, String) Executor memory of the Spark application, for example, 2 GB and 2048 MB.  
  This configuration item replaces the default parameter in **specification**.
  The unit must be provided. Otherwise, the startup fails.

* `executor_cores` - (Optional, Int) Number of CPU cores of each Executor in the Spark application.  
  This configuration item replaces the default parameter in **specification**.

* `num_executors` - (Optional, Int) Number of Executors in a Spark application.  
  This configuration item replaces the default parameter in **specification**.

* `obs_bucket` - (Optional, String) OBS bucket for storing the Spark jobs.  
  Set this parameter when you need to save jobs.

* `auto_recovery` - (Optional, Bool) Whether to enable the retry function.  
  If enabled, Spark jobs will be automatically retried after an exception occurs.
  The default value is false.

* `max_retry_times` - (Optional, Int) Maximum retry times.  
  The maximum value is 100, and the default value is 20.

<a name="SparkTemplate_Resources"></a>
The `resources` block supports:

* `name` - (Optional, String) Resource name.  
 You can also specify an OBS path, for example, obs://Bucket name/Package name.

* `type` - (Optional, String) Resource type.

<a name="SparkTemplate_Dependent_packages"></a>
The `dependent_packages` block supports:

* `name` - (Optional, String) User group name.

* `resources` - (Optional, List) User group resource.
The [resources](#SparkTemplate_Resources) structure is documented above.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

## Import

The spark template can be imported using the `id`, e.g.

```bash
$ terraform import huaweicloud_dli_spark_template.test 9680ed93-fa3f-47e5-8471-ff6e7e1a6499
```
