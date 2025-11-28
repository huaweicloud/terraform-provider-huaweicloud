---
subcategory: "Software Repository for Container (SWR)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_swr_enterprise_trigger_jobs"
description: |-
  Use this data source to get the list of SWR enterprise trigger policy jobs.
---

# huaweicloud_swr_enterprise_trigger_jobs

Use this data source to get the list of SWR enterprise trigger policy jobs.

## Example Usage

```hcl
variable "instance_id" {}
variable "namespace_name" {}
variable "policy_id" {}

data "huaweicloud_swr_enterprise_trigger_jobs" "test" {
  instance_id    = var.instance_id
  namespace_name = var.namespace_name
  policy_id      = var.policy_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the enterprise instance ID.

* `namespace_name` - (Required, String) Specifies the namespace name.

* `policy_id` - (Required, String) Specifies the trigger policy ID.

* `status` - (Optional, String) Specifies the job status.
  Valid values can be **Initialized**, **Pending**, **InProgress**, **Succeed**, **Failed**, **Stopped**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `jobs` - Indicates the trigger jobs.

  The [jobs](#jobs_struct) structure is documented below.

* `total` - Indicates the total count of the trigger jobs.

<a name="jobs_struct"></a>
The `jobs` block supports:

* `id` - Indicates the job ID.

* `policy_id` - Indicates the trigger policy ID.

* `job_detail` - Indicates the job detail.

* `event_type` - Indicates the event type.

* `notify_type` - Indicates the notify type.

* `status` - Indicates the job status.

* `status_text` - Indicates the status text.

* `created_at` - Indicates the create time.

* `updated_at` - Indicates the update time.
