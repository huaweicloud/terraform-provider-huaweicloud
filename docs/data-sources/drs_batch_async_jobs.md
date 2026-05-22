---
subcategory: "Data Replication Service (DRS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_drs_batch_async_jobs"
description: |-
  Use this data source to get the list of DRS tasks created asynchronously in batches within HuaweiCloud.
---

# huaweicloud_drs_batch_async_jobs

Use this data source to get the list of DRS tasks created asynchronously in batches within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_drs_batch_async_jobs" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `async_job_id` - (Optional, String) Specifies the ID of the tasks created asynchronously in batches.

* `status` - (Optional, String) Specifies the status of the tasks created asynchronously in batches.  
  The valid values are as follows:
  + **ASYNC_JOB_VALIDATING**: The parameters of the tasks created asynchronously in batches are being verified.
  + **ASYNC_JOB_VALIDATE_FAILED**: The parameters of the tasks created asynchronously in batches fail to be verified.
  + **AUTO_PARAM_VALIDATE_SUCCESS**: The parameters of the tasks created asynchronously in batches are successfully
    verified.
  + **COMMIT_SUCCESS**: The tasks created asynchronously in batches are successfully submitted.

* `domain_name` - (Optional, String) Specifies the tenant name of the tasks created asynchronously in batches.

* `user_name` - (Optional, String) Specifies the username of the tasks created asynchronously in batches.

* `sort_key` - (Optional, String) Specifies the keyword based on which the returned results are sorted.
  The default value is **create_time**.

* `sort_dir` - (Optional, String) Specifies the result sorting order.  
  The valid values are as follows:
  + **desc**: Descending order.
  + **asc**: Ascending order.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `jobs` - The list of tasks created asynchronously in batches.

  The [jobs](#jobs_struct) structure is documented below.

<a name="jobs_struct"></a>
The `jobs` block supports:

* `async_job_id` - The ID of the tasks created asynchronously in batches.

* `status` - The status of the tasks created asynchronously in batches.

* `domain_name` - The tenant name of the tasks created asynchronously in batches.

* `user_name` - The username of the tasks created asynchronously in batches.

* `create_time` - The time when the tasks are asynchronously created in batches.
