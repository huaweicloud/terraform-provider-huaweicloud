---
subcategory: "Distributed Message Service (DMS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dms_rocketmq_brokers"
description: |-
  Use this data source to get the list of DMS rocketMQ brokers within HuaweiCloud.
---

# huaweicloud_dms_rocketmq_brokers

Use this data source to get the list of DMS rocketMQ brokers within HuaweiCloud.

-> This data source can only be used for RocketMQ `4.8.0` version instance.

## Example Usage

```hcl
variable "instance_id" {}

data "huaweicloud_dms_rocketmq_brokers" "test" {
  instance_id = var.instance_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region where the RocketMQ brokers are located.  
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the ID of the rocketMQ instance.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `brokers` - Indicates the list of the brokers.
