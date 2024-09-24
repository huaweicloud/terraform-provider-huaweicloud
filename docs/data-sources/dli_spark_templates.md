---
subcategory: "Data Lake Insight (DLI)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dli_spark_templates"
description: |-
  Use this data source to get the list of the DLI spark templates.
---

# huaweicloud_dli_spark_templates

Use this data source to get the list of the DLI spark templates.

## Example Usage

```hcl
variable "template_id" {}

data "huaweicloud_dli_spark_templates" "test" {
  template_id = var.template_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `template_id` - (Optional, String) Specifies the ID of the spark template to be queried.

* `name` - (Optional, String) Specifies the name of the spark template to be queried.

* `group` - (Optional, String) Specifies the group name to which the spark templates belong.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `templates` - All templates that match the filter parameters.

  The [templates](#templates_struct) structure is documented below.

<a name="templates_struct"></a>
The `templates` block supports:

* `id` - The ID of template.

* `name` - The name of template.

* `group` - The group name to which the spark template belongs.

* `description` - The description of template.

* `body` - The body of template.

  The [body](#templates_body_struct) structure is documented below.

<a name="templates_body_struct"></a>
The `body` block supports:

* `name` - The spark job name.

* `dependent_packages` - The list of package resource objects.

  The [dependent_packages](#body_dependent_packages_struct) structure is documented below.

* `app_parameters` - The input parameters of the main class, that is application parameters.

* `obs_bucket` - The OBS bucket for storing the spark jobs.

* `auto_recovery` - Indicates whether to enable the retry function.

* `max_retry_times` - The maximum number of retries.

* `jars` - The name of the resource package of type jar upload to the DLI resource management system.

* `configurations` - The configuration items of the DLI spark.

* `executor_memory` - The executor memory of the spark application.

* `queue_name` - The DLI queue name.

* `main_class` - The spark main class of the template.

* `python_files` - Name of the package that is of the PyFile type.
  And has been uploaded to the DLI resource management system.

* `files` - The name of the resource package of type file upload to the DLI resource management system.

* `modules` - The name of the dependent system resource module.
  DLI provides dependencies for executing datasource jobs.
  The dependent modules and corresponding services are as follows.
  + **sys.datasource.hbase**: CloudTable/MRS HBase
  + **sys.datasource.opentsdb**: CloudTable/MRS OpenTSDB
  + **sys.datasource.rds**: RDS MySQL
  + **sys.datasource.css**: CSS

* `driver_memory` - The driver memory of the spark application.

* `executor_cores` - The number of CPU cores of each executor in the spark application.

* `app_name` - The name of the uploaded JAR or pyFile type package.

* `specification` - The compute resource type. Currently, resource types A, B, and C are available.  
  The available types and related specifications are as follows, default to minimum configuration (type **A**).

  | type | resource | driver cores | executor cores | driver memory | executor memory | num executor |
  | ---- | ---- | ---- | ---- | ---- | ---- | ---- |
  | A | 8 vCPUs, 32-GB memory | 2 | 1 | 7G | 4G | 6 |
  | B | 16 vCPUs, 64-GB memory | 2 | 2 | 7G | 8G | 7 |
  | C | 32 vCPUs, 128-GB memory | 4 | 2 | 12G | 8G | 14 |

* `resources` - The list of resource objects.

  The [resources](#body_resources_struct) structure is documented below.

* `driver_cores` - The number of CPU cores of the spark application driver.

* `num_executors` - The number of executors in a spark application.

<a name="body_dependent_packages_struct"></a>
The `dependent_packages` block supports:

* `name` - The name of a user group.

* `resources` - The resources of a user group.

  The [resources](#body_resources_struct) structure is documented below.

<a name="body_resources_struct"></a>
The `resources` block supports:

* `name` - The name of resource.

* `type` - The type of resource.
