---
subcategory: "Distributed Message Service (DMS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dms_kafka_volume_auto_expand_configuration"
description: |-
  Use this data source to query the volume auto-expansion configuration of the specified Kafka instance within HuaweiCloud.
---

# huaweicloud_dms_kafka_volume_auto_expand_configuration

Use this data source to query the volume auto-expansion configuration of the specified Kafka instance within HuaweiCloud.

## Example Usage

```hcl
variable "instance_id" {}

data "huaweicloud_dms_kafka_volume_auto_expand_configuration" "test" {
  instance_id = var.instance_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.  
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the ID of the Kafka instance.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `auto_volume_expand_enable` - Whether disk auto-expansion is enabled.

* `expand_threshold` - The threshold that triggers disk auto-expansion.

* `max_volume_size` - The maximum volume size for disk auto-expansion.

* `expand_increment` - The percentage of each disk auto-expansion.
