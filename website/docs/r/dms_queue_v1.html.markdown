---
layout: "huaweicloud"
page_title: "huaweicloud: huaweicloud_dms_queue_v1"
sidebar_current: "docs-huaweicloud-resource-dms-queue-v1"
description: |-
  Manages a DMS queue in the huaweicloud DMS Service
---

# huaweicloud\_dms\_queue_v1

Manages a DMS queue in the huaweicloud DMS Service.

## Example Usage

### Automatically detect the correct network

```hcl
resource "huaweicloud_dms_queue_v1" "queue_1" {
  name  = "queue_1"
  description  = "test create dms queue"
  queue_mode  = "FIFO"
  redrive_policy  = "enable"
  max_consume_count = 80
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Indicates the unique name of a queue. A string of 1 to 64
    characters that contain a-z, A-Z, 0-9, hyphens (-), and underscores (_).
    The name cannot be modified once specified.

* `queue_mode` - (Optional) Indicates the queue type. It only support 'NORMAL' and 'FIFO'.
    NORMAL: Standard queue. Best-effort ordering. Messages might be retrieved in an order
    different from which they were sent. Select standard queues when throughput is important.
    FIFO: First-ln-First-out (FIFO) queue. FIFO delivery. Messages are retrieved in the
    order they were sent. Select FIFO queues when the order of messages is important.
    Default value: NORMAL.

* `description` - (Optional) Indicates the basic information about a queue. The queue
    description must be 0 to 160 characters in length, and does not contain angle
    brackets (<) and (>).

* `redrive_policy` - (Optional) Indicates whether to enable dead letter messages.
    Dead letter messages indicate messages that cannot be normally consumed.
    The redrive_policy should be set to 'enable' or 'disable'. The default value is 'disable'.

* `max_consume_count` - (Optional) This parameter is mandatory only when redrive_policy is
    set to enable. This parameter indicates the maximum number of allowed message consumption
    failures. When a message fails to be consumed after the number of consumption attempts of
    this message reaches this value, DMS stores this message into the dead letter queue.
    The max_consume_count value range is 1–100.


## Attributes Reference

The following attributes are exported:


* `name` - See Argument Reference above.
* `queue_mode` - See Argument Reference above.
* `description` - See Argument Reference above.
* `redrive_policy` - See Argument Reference above.
* `max_consume_count` - See Argument Reference above.
* `created` - Indicates the time when a queue is created.
* `reservation` - Indicates the retention period (unit: min) of a message in a queue.
* `max_msg_size_byte` - Indicates the maximum message size (unit: byte) that is allowed in queue.
* `produced_messages` - Indicates the total number of messages (not including the messages that have expired and been deleted) in a queue.
* `group_count` - Indicates the total number of consumer groups in a queue.
