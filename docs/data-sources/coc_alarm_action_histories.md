---
subcategory: "Cloud Operations Center (COC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_coc_alarm_action_histories"
description: |-
  Use this data source to get the list of COC alarm action histories.
---

# huaweicloud_coc_alarm_action_histories

Use this data source to get the list of COC alarm action histories.

## Example Usage

```hcl
variable "alarm_id" {}

data "huaweicloud_coc_alarm_action_histories" "test" {
  alarm_id = var.alarm_id
}
```

## Argument Reference

The following arguments are supported:

* `alarm_id` - (Required, String) Specifies the alarm ID.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `alarm_handle_histories` - Indicates the execution details of the alarm history work order.

  The [alarm_handle_histories](#alarm_handle_histories_struct) structure is documented below.

<a name="alarm_handle_histories_struct"></a>
The `alarm_handle_histories` block supports:

* `work_order_id` - Indicates the task execution work order ID.

* `create_name` - Indicates the name of the person who created the work order.

* `create_alias` - Indicates the alias of the person who created the work order.

* `task_type` - Indicates the work order task type.

* `start_time` - Indicates the work order execution start time.

* `end_time` - Indicates the work order execution end time.

* `duration` - Indicates the time taken to execute the work order.

* `status` - Indicates the work order status.
