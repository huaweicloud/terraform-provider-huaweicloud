---
subcategory: "EventGrid (EG)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_eg_connection"
description: ""
---

# huaweicloud_eg_connection

Manages an EG connection resource within HuaweiCloud.

## Example Usage

### Connection with WEBHOOK

```hcl
variable "vpc_id" {}
variable "subnet_id" {}

resource "huaweicloud_eg_connection" "test" {
  name        = "test"
  vpc_id      = var.vpc_id
  subnet_id   = var.subnet_id
  description = "created by terraform"
  type        = "WEBHOOK"
}
```

### Connection with KAFKA

```hcl
variable "vpc_id" {}
variable "subnet_id" {}
variable "kafka_instance_id"
variable "kafka_connect_address"

resource "huaweicloud_eg_connection" "test" {
  name        = "test"
  vpc_id      = var.vpc_id
  subnet_id   = var.subnet_id
  description = "created by terraform"
  type        = "KAFKA"

  kafka_detail {
    instance_id     = var.kafka_instance_id
    connect_address = var.kafka_connect_address
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `name` - (Required, String, ForceNew) Specifies the name of the connection.
  The value can contain no more than 128 characters, including letters, digits, underscores (_), hyphens (-),
  and periods (.), and must start with a character or letter.

  Changing this parameter will create a new resource.

* `vpc_id` - (Required, String, ForceNew) Specifies the ID of the VPC to which the connection belongs.

  Changing this parameter will create a new resource.

* `subnet_id` - (Required, String, ForceNew) Specifies the ID of the subnet to which the connection belongs.

  Changing this parameter will create a new resource.

* `description` - (Optional, String) Specifies the description of the connection.

* `type` - (Optional, String, ForceNew) Specifies the type of the connection.
  The value can be **WEBHOOK** and **KAFKA**. Defaults to **WEBHOOK**.

  Changing this parameter will create a new resource.

* `kafka_detail` - (Optional, List, ForceNew) Specifies the configuration details of the kafka instance.
  This parameter is required when the `type` is set to **KAFKA**.

  Changing this parameter will create a new resource.
The [KafkaDetail](#Connection_KafkaDetail) structure is documented below.

<a name="Connection_KafkaDetail"></a>
The `KafkaDetail` block supports:

* `instance_id` - (Required, String, ForceNew) Specifies the ID of the kafka instance.

  Changing this parameter will create a new resource.

* `connect_address` - (Required, String, ForceNew) Specifies the IP address of the kafka instance.

  Changing this parameter will create a new resource.

* `user_name` - (Optional, String, ForceNew) Specifies the user name of the kafka instance.

  Changing this parameter will create a new resource.

* `password` - (Optional, String, ForceNew) Specifies the password of the kafka instance.

  Changing this parameter will create a new resource.

* `acks` - (Optional, String, ForceNew) Specifies the number of confirmation signals the prouder needs to receive
  to consider the message sent successfully. The acks represents the availability of data backup.
  The value can be:
  + **0**: Indicates that the producer does not need to wait for any confirmation of received information,
    the backup will be immediately added to the socket buffer and considered to have been sent.
    There is no guarantee that the server has successfully received the data in this case,
    and the retry configuration will not take effect and the feedback offset will always be set to -1.
  
  + **1**: Indicates that at least waiting for the leader to successfully write the data to the local log,
    but not waiting for all followers to successfully write the data. If the follower fails to successfully
    backup the data and the leader cannot provide services at this time, the message will be lost.

  + **all**: Indicates that the leader needs to wait for all backups in the ISR to be successfully written to the log.
    As long as any backup survives, the data will not be lost.

  Defaults to **1**.

  Changing this parameter will create a new resource.

* `security_protocol` - (Optional, String, ForceNew) Specifies the security protocol of the kafka instance.

  Changing this parameter will create a new resource.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `status` - Indicates the status of the connection.

* `agency` - Indicates the user-delegated name used for private network target connection.

* `created_at` - The creation time of the connection.

* `updated_at` - The last update time of the connection.

* `flavor` - The configuration details of the kafka instance.  
  The [flavor](#connection_flavor) structure is documented below.

<a name="connection_flavor"></a>
The `flavor` block supports:

* `name` - The name of the kafka instance.

* `bandwidth_type` - The bandwidth type of the kafka instance.

* `concurrency` - The concurrency number of the kafka instance.

* `concurrency_type` - The concurrency type of the kafka instance.

## Import

The connection can be imported using the `id`, e.g.

```bash
$ terraform import huaweicloud_eg_connection.test 3ea117f5-1ea3-4c27-af7f-c12c737f2ca4
```
