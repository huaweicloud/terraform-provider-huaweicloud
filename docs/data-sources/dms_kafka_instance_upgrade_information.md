---
subcategory: "Distributed Message Service (DMS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dms_kafka_instance_upgrade_information"
description: |-
  Use this data source to get the version information of the Kafka instance within HuaweiCloud.
---

# huaweicloud_dms_kafka_instance_upgrade_information

Use this data source to get the version information of the Kafka instance within HuaweiCloud.

## Example Usage

```hcl
variable "instance_id" {}

data "huaweicloud_dms_kafka_instance_upgrade_information" "test" {
  instance_id = var.instance_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region where the Kafka instance is located.  
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the ID of the Kafka instance.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `current_version` - The current version of the instance.

* `latest_version` - The latest version of the instance.
