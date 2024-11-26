---
subcategory: "Distributed Message Service (DMS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dms_kafka_message_diagnosis_task"
description: |-
  Manage DMS kafka message diagnosis task resource within HuaweiCloud.
---

# huaweicloud_dms_kafka_message_diagnosis_task

Manage DMS kafka message diagnosis task resource within HuaweiCloud.

## Example Usage

```hcl
variable "instance_id" {}
variable "group_name" {}
variable "topic_name" {}

resource "huaweicloud_dms_kafka_message_diagnosis_task" "test" {
  instance_id = var.instance_id
  group_name  = var.group_name
  topic_name  = var.topic_name
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.
  Changing this creates a new resource.

* `instance_id` - (Required, String, ForceNew) Specifies the kafka instance ID.
  Changing this creates a new resource.

* `group_name` - (Required, String, ForceNew) Specifies the group name.
  Changing this creates a new resource.

* `topic_name` - (Required, String, ForceNew) Specifies the topic name.
  Changing this creates a new resource.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `accumulated_partitions` - Indicates the number of partitions where accumulated messages are found.

* `diagnosis_dimension_list` - Indicates the diagnosis dimensions.
  The [diagnosis_dimension_list](#attrblock--diagnosis_dimension_list) structure is documented below.

* `begin_time` - Indicates the diagnosis task start time.

* `end_time` - Indicates the diagnosis task end time.

* `status` - Indicates the status of a message stack diagnosis task.
  Value can be as follows:
  + **diagnosing**
  + **failed**
  + **finished**

<a name="attrblock--diagnosis_dimension_list"></a>
The `diagnosis_dimension_list` block supports:

* `name` - Indicates the diagnosis dimension name.

* `abnormal_num` - Indicates the total number of abnormal items in this diagnosis dimension.

* `failed_num` - Indicates the total number of failed items in this diagnosis dimension.

* `diagnosis_item_list` - Indicates the diagnosis items.
  The [diagnosis_item_list](#attrblock--diagnosis_dimension_list--diagnosis_item_list) structure is documented below.

<a name="attrblock--diagnosis_dimension_list--diagnosis_item_list"></a>
The `diagnosis_item_list` block supports:

* `name` - Indicates the diagnosis item name.

* `result` - Indicates the diagnosis result.

* `advice_ids` - Indicates the suggestions for diagnosis exceptions.
  The [conclusion](#attrblock--diagnosis_dimension_list--diagnosis_item_list--conclusion) structure is documented below.

* `cause_ids` - Indicates the diagnosis exception causes.
  The [conclusion](#attrblock--diagnosis_dimension_list--diagnosis_item_list--conclusion) structure is documented below.

* `broker_ids` - Indicates the brokers affected by the diagnosis exceptions.

* `failed_partitions` - Indicates the partitions that failed to be diagnosed.

* `partitions` - Indicates the partitions affected by the diagnosis exceptions.

<a name="attrblock--diagnosis_dimension_list--diagnosis_item_list--conclusion"></a>
The `conclusion` block supports:

* `id` - Indicates the diagnosis conclusion ID.

* `params` - Indicates the diagnosis conclusion parameters.

## Timeouts

This resource provides the following timeout configuration options:

* `create` - Default is 30 minutes.

## Import

The kafka smart message diagnosis task can be imported using `instance_id` and `id` separated by a slash, e.g.

```bash
$ terraform import huaweicloud_dms_kafka_message_diagnosis_task.test <instance_id>/<id>
```
