---
subcategory: "Distributed Message Service (DMS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dms_kafka_message_diagnosis_tasks"
description: |-
  Use this data source to get the list of Kafka message diagnosis tasks.
---

# huaweicloud_dms_kafka_message_diagnosis_tasks

Use this data source to get the list of Kafka message diagnosis tasks.

## Example Usage

```hcl
variable "instance_id" {}

data "huaweicloud_dms_kafka_message_diagnosis_tasks" "test" {
  instance_id = var.instance_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the kafka instance ID.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `report_list` - Indicates the diagnosis reports.

  The [report_list](#report_list_struct) structure is documented below.

<a name="report_list_struct"></a>
The `report_list` block supports:

* `report_id` - Indicates the diagnosis report ID.

* `group_name` - Indicates the name of the consumer group being diagnosed.

* `topic_name` - Indicates the name of the topic being diagnosed.

* `accumulated_partitions` - Indicates the number of partitions where accumulated messages are found.

* `status` - Indicates the status of a message stack diagnosis task.

* `begin_time` - Indicates the diagnosis task start time.

* `end_time` - Indicates the diagnosis task end time.
