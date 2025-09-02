---
subcategory: "Distributed Message Service (DMS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dms_rocketmq_instance_diagnosis"
description: |-
  Manages DMS RocketMQ instance diagnosis resource within HuaweiCloud.
---

# huaweicloud_dms_rocketmq_instance_diagnosis

Manages DMS RocketMQ instance diagnosis resource within HuaweiCloud.

## Example Usage

```hcl
variable "instance_id" {}
variable "group_name" {}

resource "huaweicloud_dms_rocketmq_instance_diagnosis" "test" {
  instance_id = var.instance_id
  group_name  = var.group_name
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.  
  If omitted, the provider-level region will be used.  
  Changing this parameter will create a new resource.

* `instance_id` - (Required, String, ForceNew) Specifies the ID of the RocketMQ instance to be diagnosed.  
  Changing this parameter will create a new resource.

* `group_name` - (Required, String, ForceNew) Specifies the name of the consumer group to be diagnosed.  
  Changing this parameter will create a new resource.

* `node_ids` - (Optional, List, ForceNew) Specifies the list of node IDs to be diagnosed.  
  Changing this parameter will create a new resource.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The report ID of the diagnosis.

* `report_id` - The ID of the diagnosis report.

* `consumer_nums` - The number of consumers.

* `status` - The status of the diagnosis task.

* `created_at` - The creation time of the diagnosis report.

* `abnormal_item_sum` - The number of abnormal items.

* `faulted_node_sum` - The number of abnormal nodes.

* `online` - Whether the consumer group is online.

* `message_accumulation` - The number of accumulated messages.

* `subscription_consistency` - Whether the subscription is consistent.

* `subscriptions` - The list of subscribers.

* `diagnosis_node_reports` - The list of diagnosis node report.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 10 minutes.

## Import

The RocketMQ instance diagnosis can be imported using the `id`, e.g.

```bash
$ terraform import huaweicloud_dms_rocketmq_instance_diagnosis.test <id>
```
