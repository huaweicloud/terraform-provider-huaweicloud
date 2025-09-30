---
subcategory: "Distributed Message Service (DMS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dms_kafka_instance_monitoring_dimensions"
description: |-
  Use this data source to get the monitoring dimensions of Kafka instance within HuaweiCloud.
---

# huaweicloud_dms_kafka_instance_monitoring_dimensions

Use this data source to get the monitoring dimensions of Kafka instance within HuaweiCloud.

## Example Usage

```hcl
variable "instance_id" {}

data "huaweicloud_dms_kafka_instance_monitoring_dimensions" "test" {
  instance_id = var.instance_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region where the Kafka instance monitoring dimensions are located.
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the ID of the Kafka instance.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `dimensions` - The list of monitoring dimensions.  
  The [dimensions](#Kafka_instance_monitoring_dimensions_struct) structure is documented below.

* `instance_ids` - The list of instance IDs.  
  The [instance_ids](#Kafka_instance_monitoring_instance_ids_struct) structure is documented below.

* `nodes` - The list of nodes.  
  The [nodes](#Kafka_instance_monitoring_nodes_struct) structure is documented below.

* `queues` - The list of queues.  
  The [queues](#Kafka_instance_monitoring_queues_struct) structure is documented below.

* `groups` - The list of consumer groups.  
  The [groups](#Kafka_instance_monitoring_groups_struct) structure is documented below.

<a name="Kafka_instance_monitoring_dimensions_struct"></a>
The `dimensions` block supports:

* `name` - The name of the monitoring dimension.

* `metrics` - The list of monitoring metric names.

* `key_name` - The list of keys used for monitoring queries.

* `dim_router` - The list of monitoring dimension routes.

* `children` - The list of child dimensions.  
  The [children](#Kafka_instance_monitoring_children_struct) structure is described below.

<a name="Kafka_instance_monitoring_children_struct"></a>
The `children` block supports:

* `name` - The name of the child dimension.

* `metrics` - The list of monitoring metric names.

* `key_name` - The list of keys used for monitoring queries.

* `dim_router` - The list of monitoring dimension routes.

<a name="Kafka_instance_monitoring_instance_ids_struct"></a>
The `instance_ids` block supports:

* `name` - The ID of the instance.

<a name="Kafka_instance_monitoring_nodes_struct"></a>
The `nodes` block supports:

* `name` - The name of the node.

<a name="Kafka_instance_monitoring_queues_struct"></a>
The `queues` block supports:

* `name` - The name of the topic.

* `partitions` - The list of partitions.  
  The [partitions](#Kafka_instance_monitoring_partitions_struct) structure is described below.

<a name="Kafka_instance_monitoring_partitions_struct"></a>
The `partitions` block supports:

* `name` - The name of the partition.

<a name="Kafka_instance_monitoring_groups_struct"></a>
The `groups` block supports:

* `name` - The name of the consumer group.

* `queues` - The list of queues.  
  The [queues](#Kafka_instance_monitoring_queues_struct) structure is described below.
