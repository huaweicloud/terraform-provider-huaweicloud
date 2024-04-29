---
subcategory: "EventGrid (EG)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_eg_custom_event_source"
description: ""
---

# huaweicloud_eg_custom_event_source

Using this resource to manage an EG custom event source within Huaweicloud.

## Example Usage

### Create a custom event source for the application type

```hcl
variable "channel_id" {}
variable "source_name" {}

resource "huaweicloud_eg_custom_event_source" "test" {
  channel_id  = var.channel_id
  type        = "APPLICATION"
  name        = var.source_name
  description = "Created by script"
}
```

### Create a custom event source for the RabbitMQ type

```hcl
variable "channel_id" {}
variable "source_name" {}
variable "rabbitmq_instance_id" {}
variable "rabbitmq_user_name" {}
variable "rabbitmq_user_password" {}
variable "rabbitmq_vhost_name" {}
variable "rabbitmq_queue_name" {}

resource "huaweicloud_eg_custom_event_source" "test" {
  channel_id = var.channel_id
  type       = "RABBITMQ"
  name       = var.source_name
  detail     = jsonencode({
    instance_id = var.rabbitmq_instance_id
    user_name   = var.rabbitmq_user_name
    password    = var.rabbitmq_user_password
    vhost_name  = var.rabbitmq_vhost_name
    queue_name  = var.rabbitmq_queue_name
  })
}
```

### Create a custom event source for the RocketMQ type

```hcl
variable "channel_id" {}
variable "source_name" {}
variable "rocketmq_instance_id" {}
variable "rocketmq_instance_name" {}
variable "rocketmq_instance_namesrv_address" {}
variable "rocketmq_consumer_group_id" {}
variable "rocketmq_topic_id" {}

resource "huaweicloud_eg_custom_event_source" "test" {
  channel_id = var.channel_id
  type       = "ROCKETMQ"
  name       = var.source_name
  detail     = jsonencode({
    instance_id     = var.rocketmq_instance_id
    name            = var.rocketmq_instance_name
    namesrv_address = var.rocketmq_instance_namesrv_address
    group           = var.rocketmq_consumer_group_id
    topic           = var.rocketmq_topic_id
    enable_acl      = false
    ssl_enable      = false
  })
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region where the custom event channel and custom event source
  are located. If omitted, the provider-level region will be used.  
  Changing this will create a new resource.

* `channel_id` - (Required, String, ForceNew) Specifies the ID of the custom event channel to which the custom event
  source belongs.  
  Changing this will create a new resource.

* `name` - (Required, String, ForceNew) Specifies the name of the custom event source.  
  The valid length is limited from `1` to `128`, only lowercase letters, digits, hyphens (-), underscores (_) are
  allowed. The name must start with a lowercase letter or digit.  
  Changing this will create a new resource.

* `type` - (Optional, String, ForceNew) Specifies the type of the custom event source.
  The valid values are as follows:
  + **APPLICATION**
  + **RABBITMQ**
  + **ROCKETMQ**

  Defaults to **APPLICATION**.  
  Changing this will create a new resource.

  -> Before creating a **RocketMQ** event source, you need to open ingress rule for TCP `8,100` and `10,100`-`10,103`
     ports for the security group.

* `description` - (Optional, String) Specifies the description of the custom event source.

* `detail` - (Optional, String) Specifies the configuration detail of the event source, in JSON format.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - A resource ID in UUID format.

* `status` - The status of the custom event source.

* `created_at` - The (UTC) creation time of the custom event source, in RFC3339 format.

* `updated_at` - The (UTC) update time of the custom event source, in RFC3339 format.

## Import

Custom event sources can be imported by their `id`, e.g.

```bash
terraform import huaweicloud_eg_custom_event_source.test <id>
```
