---
subcategory: "Distributed Message Service (DMS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dms_kafka_instance_upgrade"
description: |-
  Use this resource to upgrade Kafka instance within HuaweiCloud.
---

# huaweicloud_dms_kafka_instance_upgrade

Use this resource to upgrade Kafka instance within HuaweiCloud.

-> This resource is a one-time action resource for upgrading the Kafka instance. Deleting this resource will not
   recover the upgrade Kafka instance, but will only remove the resource information from the tfstate file.

## Example Usage

### Immediate upgrade

```hcl
variable "instance_id" {}

resource "huaweicloud_dms_kafka_instance_upgrade" "test" {
  instance_id = var.instance_id
}
```

### Scheduled upgrade

```hcl
variable "instance_id" {}
variable "execute_at" {}

resource "huaweicloud_dms_kafka_instance_upgrade" "test" {
  instance_id = var.instance_id
  is_schedule = true
  execute_at  = var.execute_at
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region where the Kafka instance to be upgraded is located.
  If omitted, the provider-level region will be used.
  Changing this creates a new resource.

* `instance_id` - (Required, String, NonUpdatable) Specifies the ID of the Kafka instance.

* `is_schedule` - (Optional, Bool, NonUpdatable) Specifies whether to execute as a scheduled task.  
  Defaults to **false**.

* `execute_at` - (Optional, Int, NonUpdatable) Specifies the scheduled time in Unix timestamp format, in milliseconds.
  This parameter is **required** and available only when `is_schedule` is **true**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID in UUID format.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 30 minutes.
