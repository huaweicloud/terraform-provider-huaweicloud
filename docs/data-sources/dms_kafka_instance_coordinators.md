---
subcategory: "Distributed Message Service (DMS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dms_kafka_instance_coordinators"
description: |-
  Use this data source to get the coordinator list of the Kafka instance within HuaweiCloud.
---

# huaweicloud_dms_kafka_instance_coordinators

Use this data source to get the coordinator list of the Kafka instance within HuaweiCloud.

## Example Usage

```hcl
variable "instance_id" {}

data "huaweicloud_dms_kafka_instance_coordinators" "test" {
  instance_id = var.instance_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region where the coordinators are located.  
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the ID of the Kafka instance to which the coordinators belong.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `coordinators` - The list of coordinators corresponding to all consumer groups.  
  The [coordinators](#kafka_instance_coordinators_attr) structure is documented below.

<a name="kafka_instance_coordinators_attr"></a>
The `coordinators` block supports:

* `id` - The ID of the broker of the corresponding to the coordinator.

* `group_id` - The ID of the consumer group.

* `host` - The address of the broker of the corresponding to the coordinator.

* `port` - The port number of the corresponding to the coordinator.
