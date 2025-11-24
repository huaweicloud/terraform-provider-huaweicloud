---
subcategory: "Software Repository for Container (SWR)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_swr_enterprise_retention_policy_execution_record_sub_tasks"
description: |-
  Use this data source to get the list of SWR enterprise retention policy execution record sub tasks.
---

# huaweicloud_swr_enterprise_retention_policy_execution_record_sub_tasks

Use this data source to get the list of SWR enterprise retention policy execution record sub tasks.

## Example Usage

```hcl
variable "instance_id" {}
variable "execution_id" {}
variable "task_id" {}

data "huaweicloud_swr_enterprise_retention_policy_execution_record_sub_tasks" "test" {
  instance_id  = var.instance_id
  execution_id = var.execution_id
  task_id      = var.task_id
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

* `task_id` - (Required, String) Specifies the execution record task ID.

* `status` - (Optional, String) Specifies the execution record task status.
  Valid values are **Initialized**, **Pending**, **InProgress**, **Succeed**, **Failed**, **Stopped**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `sub_tasks` - Indicates the execution record sub tasks.

  The [sub_tasks](#sub_tasks_struct) structure is documented below.

* `total` - Indicates the total sub tasks.

<a name="sub_tasks_struct"></a>
The `sub_tasks` block supports:

* `id` - Indicates the execution sub task ID.

* `job_id` - Indicates the execution sub task job ID.

* `namespace` - Indicates the namespace.

* `tag` - Indicates the image version tags.

* `action` - Indicates the retention action.

* `repository` - Indicates the repository name.

* `digest` - Indicates the image sha256.

* `status` - Indicates the sub task status.

* `status_text` - Indicates the sub task status text.

* `op_time` - Indicates the operation time.
