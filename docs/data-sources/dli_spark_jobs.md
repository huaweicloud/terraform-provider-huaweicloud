---
subcategory: "Data Lake Insight (DLI)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dli_spark_jobs"
description: |-
  Use this data source to query the list of the DLI spark jobs within HuaweiCloud.
---

# huaweicloud_dli_spark_jobs

Use this data source to query the list of the DLI spark jobs within HuaweiCloud.

## Example Usage

### Basic Usage

```hcl
data "huaweicloud_dli_spark_jobs" "test" {}
```

### Filter by job ID

```hcl
variable "job_id" {}

data "huaweicloud_dli_spark_jobs" "test" {
  job_id = var.job_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region where the spark jobs are located.  
  If omitted, the provider-level region will be used.

* `cluster_name` - (Optional, String) Specifies the DLI queue name of the spark job to be queried.

* `queue_name` - (Optional, String) Specifies the queue name of the spark job to be queried.

* `job_name` - (Optional, String) Specifies the name of the spark job to be queried.

* `job_id` - (Optional, String) Specifies the ID of the spark job to be queried.

* `state` - (Optional, String) Specifies the state of the spark job to be queried.  
  The valid values are as follows:
  + **starting**: The job is being started.
  + **running**: The job is running.
  + **dead**: The job is dead.
  + **success**: The job is successful.
  + **recovering**: The job is being recovered.

* `owner` - (Optional, String) Specifies the owner of the spark job to be queried.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `jobs` - The list of spark jobs that matched filter parameters.  
  The [jobs](#spark_jobs_attr) structure is documented below.

<a name="spark_jobs_attr"></a>
The `jobs` block supports:

* `id` - The ID of the spark job.

* `app_id` - The backend app ID of the spark job.

* `name` - The name of the spark job.

* `owner` - The owner of the spark job.

* `queue` - The queue of the spark job.

* `cluster_name` - The cluster name of the spark job.

* `state` - The state of the spark job.

* `kind` - The type of the spark job.

* `duration` - The running duration of the spark job, in milliseconds.

* `sc_type` - The compute resource type of the spark job.

* `image` - The custom image of the spark job.

* `log` - The last `10` log records of the spark job.

* `req_body` - The request body details of the spark job.

* `proxy_user` - The proxy user of the spark job.

* `created_at` - The creation time of the spark job, in RFC3339 format.

* `updated_at` - The update time of the spark job, in RFC3339 format.
