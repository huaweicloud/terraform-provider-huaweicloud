---
subcategory: "Distributed Message Service (DMS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dms_kafka_instance_batch_action"
description: |-
  Use this resource to batch operate Kafka instances within HuaweiCloud.
---

# huaweicloud_dms_kafka_instance_batch_action

Use this resource to batch operate Kafka instances within HuaweiCloud.

-> This resource is only a one-time action resource for restarting or deleting Kafka instances. Deleting this
   resource will not clear the corresponding request record, but will only remove the resource information from the
   tfstate file.

## Example Usage

### Restart Kafka instances

```hcl
variable "instance_ids" {
  type = list(string)
}

resource "huaweicloud_dms_kafka_instance_batch_action" "test" {
  instance_ids = var.instance_ids
  action       = "restart"
}
```

### Delete Kafka instances

```hcl
variable "instance_ids" {
  type = list(string)
}

resource "huaweicloud_dms_kafka_instance_batch_action" "test" {
  instance_ids = var.instance_ids
  action       = "delete"
}
```

### Delete all failed Kafka instances

```hcl
resource "huaweicloud_dms_kafka_instance_batch_action" "test" {
  action      = "delete"
  all_failure = "kafka"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region where the instances to be operated are located.  
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `action` - (Required, String, NonUpdatable) Specifies the type of the action.

* `instances` - (Optional, List, NonUpdatable) Specifies the list of instance IDs to be operated.

* `all_failure` - (Optional, String, NonUpdatable) Specifies whether to delete all instances that failed
  to be created.  
  The valid values are as follows:
  + **kafka**: Delete all failed Kafka instances.

* `force_delete` - (Optional, Bool, NonUpdatable) Specifies whether to force delete instances.  
  Defaults to **false**. Force delete instances do not enter the recycle bin.
