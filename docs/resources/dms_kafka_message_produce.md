---
subcategory: "Distributed Message Service (DMS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dms_kafka_message_produce"
description: |-
  Manages a DMS kafka message produce resource within HuaweiCloud.
---

# huaweicloud_dms_kafka_message_produce

Manages a DMS kafka message produce resource within HuaweiCloud.

## Example Usage

```hcl
variable "instance_id" {}
variable "topic_name" {}

resource "huaweicloud_dms_kafka_message_produce" "test" {
  instance_id = var.instance_id
  topic       = var.topic_name
  body        = "test"

  property_list {
    name  = "KEY"
    value = "testKey"
  }

  property_list {
    name  = "PARTITION"
    value = "1"
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.
  Changing this creates a new resource.

* `instance_id` - (Required, String, ForceNew) Specifies the instance ID.
  Changing this creates a new resource.

* `topic` - (Required, String, ForceNew) Specifies the topic name.
  Changing this creates a new resource.

* `body` - (Required, String, ForceNew) Specifies the message content.
  Changing this creates a new resource.

* `property_list` - (Optional, List, ForceNew) Specifies the topic partition information.
  Changing this creates a new resource.
  The [property_list](#block--property_list) structure is documented below.

<a name="block--property_list"></a>
The `property_list` block supports:

* `name` - (Required, String, ForceNew) Specifies the feature name.
  + **KEY**: Specifies the message key.
  + **PARTITION** : Specifies the partition to which the message will be sent.
  Changing this creates a new resource.

* `value` - (Required, String, ForceNew) Specifies the feature value.
  Changing this creates a new resource.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID in UUID format.
