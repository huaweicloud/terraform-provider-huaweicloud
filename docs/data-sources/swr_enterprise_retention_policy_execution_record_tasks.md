---
subcategory: "Software Repository for Container (SWR)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_swr_enterprise_retention_policy_execution_record_tasks"
description: |-
  Use this data source to get the list of SWR enterprise retention policy execution record tasks.
---

# huaweicloud_swr_enterprise_retention_policy_execution_record_tasks

Use this data source to get the list of SWR enterprise retention policy execution record tasks.

## Example Usage

```hcl
variable "instance_id" {}
variable "namespace_name" {}
variable "policy_id" {}
variable "execution_id" {}

data "huaweicloud_swr_enterprise_retention_policy_execution_record_tasks" "test" {
  instance_id    = var.instance_id
  namespace_name = var.namespace_name
  policy_id      = var.policy_id
  execution_id   = var.execution_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the enterprise instance ID.

* `namespace_name` - (Required, String) Specifies the namespace name.

* `policy_id` - (Required, String) Specifies the policy ID.

* `execution_id` - (Required, String) Specifies the execution record ID.

* `status` - (Optional, String) Specifies the task status.
  Values can be **Initialized**, **Pending**, **InProgress**, **Succeed**, **Failed**, **Stopped**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `tasks` - Indicates the execution records.

  The [tasks](#tasks_struct) structure is documented below.

* `total` - Indicates the total counts of tasks.

<a name="tasks_struct"></a>
The `tasks` block supports:

* `id` - Indicates the execution task ID.

* `execution_id` - Indicates the execution record ID.

* `job_id` - Indicates the job ID.

* `status` - Indicates the status.

* `status_code` - Indicates the status code.

* `repository` - Indicates the repository name.

* `start_time` - Indicates the start time.

* `end_time` - Indicates the end time.

* `status_revision` - Indicates the status revision.

* `retained` - Indicates the retained counts of version.

* `total` - Indicates the total counts of version.
