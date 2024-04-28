---
subcategory: "Distributed Message Service (DMS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dms_rocketmq_migration_task"
description: ""
---

# huaweicloud_dms_rocketmq_migration_task

Manages a DMS RocketMQ migration task resource within HuaweiCloud.

## Example Usage

### RoecktMQ migration task from RocketMQ to RocketMQ

```hcl
variable "instance_id" {}
variable "name" {}
variable "topic_name" {}
variable "group_name" {}

resource "huaweicloud_dms_rocketmq_migration_task" "test" {
  instance_id = var.instance_id
  overwrite   = "true"
  name        = var.name
  type        = "rocketmq"

  topic_configs {
    order             = false
    perm              = 6
    read_queue_num    = 16
    topic_filter_type = "SINGLE_TAG"
    topic_name        = var.topic_name
    topic_sys_flag    = 0
    write_queue_num   = 16
  }

  subscription_groups  {     
    consume_broadcast_enable          = true
    consume_enable                    = true
    consume_from_min_enable           = true
    group_name                        = var.group_name
    notify_consumerids_changed_enable = true
    retry_max_times                   = 16
    retry_queue_num                   = 1
    which_broker_when_consume_slow    = 1        
  }
}
```

### RoecktMQ migration task from RabbitMQ to RocketMQ

```hcl
variable "instance_id" {}
variable "name" {}
variable "vhost_name" {}
variable "queue_name" {}
variable "exchange_name" {}

resource "huaweicloud_dms_rocketmq_migration_task" "test" {
  instance_id = var.instance_id
  overwrite   = "true"
  name        = var.name
  type        = "rabbitToRocket"

  vhosts {
    name = var.vhost_name
  }

  queues {
    name    = var.queue_name
    vhost   = var.vhost_name
    durable = false
  }
  
  exchanges {
    name    = var.exchange_name
    vhost   = var.vhost_name
    type    = "topic"
    durable = false
  }

  bindings {
    source           = var.exchange_name
    vhost            = var.vhost_name
    destination      = var.queue_name
    destination_type = "queue"
    routing_key      = var.queue_name
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this creates a new resource.

* `instance_id` - (Required, String, ForceNew) Specifies the ID of the RocketMQ instance.
  Changing this creates a new resource.

* `name` - (Required, String, ForceNew) Specifies the name of the migration task.
  Changing this creates a new resource.

* `overwrite` - (Required, String, ForceNew) Specifies whether to overwrite configurations with the same name.
  Value options:
  + **true**: Configurations in the destination metadata with the same name as the source metadata will be overwritten.
  + **false**: An error is reported when a topic or group already exists.
  Changing this creates a new resource.

* `type` - (Required, String, ForceNew) Specifies the migration task type.
  Value options: **rocketmq** or **rabbitToRocket**. Changing this creates a new resource.

* `topic_configs` - (Optional, List, ForceNew) Specifies the topic metadata.
  The [topic_configs](#RocketMQ_migration_task_topic_configs) structure is documented below.
  Changing this creates a new resource.

* `subscription_groups` - (Optional, List, ForceNew) Specifies the consumer group metadata.
  The [subscription_groups](#RocketMQ_migration_task_subscription_groups) structure is documented below.
  Changing this creates a new resource.

-> **NOTE:** Parameters `topic_configs` and `subscription_groups` are required when `type` is set to **rocketmq**.

* `vhosts` - (Optional, List, ForceNew) Specifies the virtual hosts metadata.
  The [vhosts](#RocketMQ_migration_task_vhosts) structure is documented below.
  Changing this creates a new resource.

* `queues` - (Optional, List, ForceNew) Specifies the queue metadata.
  The [queues](#RocketMQ_migration_task_queues) structure is documented below.
  Changing this creates a new resource.

* `exchanges` - (Optional, List, ForceNew) Specifies the exchange metadata.
  The [exchanges](#RocketMQ_migration_task_exchanges) structure is documented below.
  Changing this creates a new resource.

* `bindings` - (Optional, List, ForceNew) Specifies the binding metadata.
  The [bindings](#RocketMQ_migration_task_bindings) structure is documented below.
  Changing this creates a new resource.

-> **NOTE:** Parameters `vhosts`, `queues`, `exchanges` and `bindings` are required when `type` is set to **rabbitToRocket**.

<a name="RocketMQ_migration_task_topic_configs"></a>
The `topic_configs` block supports:

* `topic_name` - (Required, String, ForceNew) Specifies the topic name. Changing this creates a new resource.

* `order` - (Optional, Bool, ForceNew) Specifies whether a message is an ordered message.
  Changing this creates a new resource.

* `perm` - (Optional, Int, ForceNew) Specifies the number of permission. Changing this creates a new resource.

* `read_queue_num` - (Optional, Int, ForceNew) Specifies the number of read queues.
  Changing this creates a new resource.

* `topic_filter_type` - (Optional, String, ForceNew) Specifies the filter type of a topic.
  Value options: **SINGLE_TAG**, **MULTI_TAG**. Changing this creates a new resource.

* `topic_sys_flag` - (Optional, Int, ForceNew) Specifies the system flag of a topic.
  Changing this creates a new resource.

* `write_queue_num` - (Optional, Int, ForceNew) Specifies the number of write queues.
  Changing this creates a new resource.

<a name="RocketMQ_migration_task_subscription_groups"></a>
The `subscription_groups` block supports:

* `group_name` - (Required, String, ForceNew) Specifies the name of a consumer group.
  Changing this creates a new resource.

* `consume_broadcast_enable` - (Optional, Bool, ForceNew) Specifies whether to enable broadcast.
  Changing this creates a new resource.

* `consume_enable` - (Optional, Bool, ForceNew) Specifies whether to enable consumption.
  Changing this creates a new resource.

* `consume_from_min_enable` - (Optional, Bool, ForceNew) Specifies whether to enable consumption from the earliest
  offset. Changing this creates a new resource.

* `notify_consumerids_changed_enable` - (Optional, Bool, ForceNew) Specifies whether to notify changes of consumer IDs.
  Changing this creates a new resource.

* `retry_max_times` - (Optional, Int, ForceNew) Specifies the maximum number of consumption retries.
  Changing this creates a new resource.

* `retry_queue_num` - (Optional, Int, ForceNew) Specifies the number of retry queues.
  Changing this creates a new resource.

* `which_broker_when_consume_slow` - (Optional, Int, ForceNew) Specifies the ID of the broker selected for slow
  consumption. Changing this creates a new resource.

<a name="RocketMQ_migration_task_vhosts"></a>
The `vhosts` block supports:

* `name` - (Optional, String, ForceNew) Specifies the virtual host name. Changing this creates a new resource.

<a name="RocketMQ_migration_task_queues"></a>
The `queues` block supports:

* `name` - (Optional, String, ForceNew) Specifies the queue name. Changing this creates a new resource.

* `vhost` - (Optional, String, ForceNew) Specifies the virtual host name. Changing this creates a new resource.

* `durable` - (Optional, Bool, ForceNew) Specifies whether to enable data persistence.
  Changing this creates a new resource.

<a name="RocketMQ_migration_task_exchanges"></a>
The `exchanges` block supports:

* `name` - (Optional, String, ForceNew) Specifies the switch name. Changing this creates a new resource.

* `vhost` - (Optional, String, ForceNew) Specifies the virtual host name. Changing this creates a new resource.

* `durable` - (Optional, Bool, ForceNew) Specifies whether to enable data persistence.
  Changing this creates a new resource.

* `type` - (Optional, String, ForceNew) Specifies the exchange type. Changing this creates a new resource.

<a name="RocketMQ_migration_task_bindings"></a>
The `bindings` block supports:

* `vhost` - (Optional, String, ForceNew) Specifies the virtual host name. Changing this creates a new resource.

* `source` - (Optional, String, ForceNew) Specifies the message source. Changing this creates a new resource.

* `destination` - (Optional, String, ForceNew) Specifies the message target. Changing this creates a new resource.

* `destination_type` - (Optional, String, ForceNew) Specifies the message target type.
  Changing this creates a new resource.

* `routing_key` - (Optional, String, ForceNew) Specifies the routing key. Changing this creates a new resource.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID. The value is the RocketMQ migration task ID.

* `start_date` - Indicates the start time of the migration task.

* `status` - Indicates the status of the migration task. The value can be **finished** or **failed***.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 50 minutes.
* `delete` - Default is 10 minutes.

## Import

The RocketMQ migration task can be imported using the RocketMQ instance ID and the RocketMQ migration task ID
separated by a slash, e.g.

```bash
$ terraform import huaweicloud_dms_rocketmq_migration_task.test <instance_id>/<id>
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response, security or some other reason. The missing attribute includes: `overwrite`.
It is generally recommended running `terraform plan` after importing the task. You can then decide
if changes should be applied to the task, or the resource definition should be updated to align with the task.
Also you can ignore changes as below.

```hcl
resource "huaweicloud_dms_rocketmq_migration_task" "test" {
    ...

  lifecycle {
    ignore_changes = [
      overwrite,
    ]
  }
}
```
