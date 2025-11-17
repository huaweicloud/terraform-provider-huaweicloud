---
subcategory: "Software Repository for Container (SWR)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_swr_enterprise_replication_policy_execution_record_tasks"
description: |-
   Use this data source to get the list of SWR enterprise replication policy execution record tasks.
---

# huaweicloud_swr_enterprise_replication_policy_execution_record_tasks

Use this data source to get the list of SWR enterprise replication policy execution record tasks.

## Example Usage

```hcl
variable "instance_id" {}
variable "execution_id" {}

data "huaweicloud_swr_enterprise_replication_policy_execution_record_tasks" "test" {
  instance_id  = var.instance_id
  execution_id = var.execution_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the enterprise instance ID.

* `execution_id` - (Required, String) Specifies the execution record ID.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `tasks` - Indicates the execution records.

  The [tasks](#tasks_struct) structure is documented below.

* `total` - Indicates the total records.

<a name="tasks_struct"></a>
The `tasks` block supports:

* `id` - Indicates the task ID.

* `execution_id` - Indicates the execution ID.

* `operation` - Indicates the operation.

* `resource_type` - Indicates the resource type.

* `job_id` - Indicates the job ID.

* `src_resource` - Indicates the source resource.

* `dst_resource` - Indicates the destination resource.

* `status` - Indicates the status.

* `status_revision` - Indicates the status revision.

* `start_time` - Indicates the start time.

* `end_time` - Indicates the end time.
