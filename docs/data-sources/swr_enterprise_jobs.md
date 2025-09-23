---
subcategory: "Software Repository for Container (SWR)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_swr_enterprise_jobs"
description: |-
  Use this data source to get the list of SWR enterprise jobs.
---

# huaweicloud_swr_enterprise_jobs

Use this data source to get the list of SWR enterprise jobs.

## Example Usage

```hcl
data "huaweicloud_swr_enterprise_jobs" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `status` - (Optional, String) Specifies the job status.
  Value can be **Creating**, **Initializing**, **Running**, **Failed**, **Success**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `jobs` - Indicates the jobs.

  The [jobs](#jobs_struct) structure is documented below.

<a name="jobs_struct"></a>
The `jobs` block supports:

* `id` - Indicates the job ID.

* `instance_id` - Indicates the instance ID.

* `resource_id` - Indicates the resource ID.

* `resource_name` - Indicates the resource name.

* `type` - Indicates the job type.

* `status` - Indicates the job status.

* `reason` - Indicates the failed reason.

* `created_at` - Indicates the create time.

* `updated_at` - Indicates the update time.
