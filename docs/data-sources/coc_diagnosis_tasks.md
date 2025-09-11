---
subcategory: "Cloud Operations Center (COC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_coc_diagnosis_tasks"
description: |-
  Use this data source to get the list of COC diagnosis tasks.
---

# huaweicloud_coc_diagnosis_tasks

Use this data source to get the list of COC diagnosis tasks.

## Example Usage

```hcl
data "huaweicloud_coc_diagnosis_tasks" "test" {}
```

## Argument Reference

The following arguments are supported:

* `task_id` - (Optional, String) Specifies the diagnostic task work order ID.

* `type` - (Optional, String) Specifies the instance types supported by the diagnostic task.
  The value can be **ECS**, **RDS**, **DCS**, **DMS** or **ELB**.

* `status` - (Optional, String) Specifies the execution status of the diagnostic task.
  Values can be as follows:
  + **cancel**: Cancelled.
  + **executing**: Executing.
  + **waiting**: Waiting to be executed.
  + **failed**: Exception.
  + **finish**: Completed.

* `region` - (Optional, String) Specifies the region where the diagnosed instance is located.

* `creator` - (Optional, String) Specifies the IAM user ID of the person who created the diagnostic task ticket.

* `start_time` - (Optional, Int) Specifies the start time of the diagnostic work order.

* `end_time` - (Optional, Int) Specifies the end time of the diagnostic work order.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `data` - Indicates the diagnostic record structure.

  The [data](#data_data_struct) structure is documented below.

<a name="data_data_struct"></a>
The `data` block supports:

* `id` - Indicates the diagnostic task node ID.

* `code` - Indicates the step code.

* `domain_id` - Indicates the account ID to which the diagnostic record belongs.

* `project_id` - Indicates the project ID to which the diagnosed instance belongs.

* `user_id` - Indicates the user ID to which the diagnostic record belongs.

* `user_name` - Indicates the user name to which the diagnostic record belongs.

* `progress` - Indicates the progress of the diagnostic task execution.

* `work_order_id` - Indicates the diagnostic task ticket ID.

* `instance_id` - Indicates the ID of the instance being diagnosed.

* `instance_name` - Indicates the name of the instance being diagnosed.

* `type` - Indicates the type of the diagnosed instance.

* `status` - Indicates the execution status of the diagnostic task.

* `start_time` - Indicates the start time.

* `end_time` - Indicates the end time.

* `instance_num` - Indicates the number of instances included in the diagnostic task.

* `os_type` - Indicates the operating system type of the diagnosed instance.

* `region` - Indicates the region where the diagnosed instance is located.
