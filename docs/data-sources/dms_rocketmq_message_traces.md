---
subcategory: "Distributed Message Service (DMS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dms_rocketmq_message_traces"
description: |-
  Use this data source to get the list of RocketMQ message traces.
---

# huaweicloud_dms_rocketmq_message_traces

Use this data source to get the list of RocketMQ message traces.

## Example Usage

```hcl
variable "instance_id" {}
variable "message_id" {}

data "huaweicloud_dms_rocketmq_message_traces" "test" {
  instance_id = var.instance_id
  message_id  = var.message_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the instance ID.

* `message_id` - (Required, String) Specifies the message ID.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `traces` - Specifies the message trace list.

  The [traces](#traces_struct) structure is documented below.

<a name="traces_struct"></a>
The `traces` block supports:

* `message_id` - Specifies the message ID.

* `consume_status` - Specifies the consumption status.
  + **0**: successful
  + **1**: timeout
  + **2**: abnormal
  + **3**: null
  + **5**: failed

* `message_type` - Specifies the message type.

* `keys` - Specifies the message keys.

* `body_length` - Specifies the message body length.

* `offset_message_id` - Specifies the offset message ID.

* `time` - Specifies the time.

* `cost_time` - Specifies the time spent.

* `topic` - Specifies the topic name.

* `group_name` - Specifies the producer group or consumer group.

* `client_host` - Specifies the IP address of the host that generates the message.

* `store_host` - Specifies the IP address of the host that stores the message.

* `transaction_id` - Specifies the transaction ID.

* `transaction_state` - Specifies the transaction status.

* `from_transaction_check` - Specifies whether the response is a transaction check response.

* `trace_type` - Specifies the trace type.

* `retry_times` - Specifies the number of retry times.

* `tags` - Specifies the message tag.

* `request_id` - Specifies the request ID.

* `success` - Specifies whether the request is successful.
