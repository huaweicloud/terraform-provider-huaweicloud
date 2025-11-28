---
subcategory: "EventGrid (EG)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_eg_event_stream"
description: ""
---

# huaweicloud_eg_event_stream

Using this resource to manage an EG event stream within Huaweicloud.

## Example Usage

### Create an event stream with event source type Kafka

```hcl
variable "stream_name" {}
variable "kafka_connect_address" {}
variable "kafka_consumer_group_id" {}
variable "kafka_instance_name" {}
variable "kafka_instance_id" {}
variable "kafka_topic_name" {}
variable "eg_connection_id" {}

resource "huaweicloud_eg_event_stream" "test" {
  name   = var.stream_name
  action = "START"

  source {
    name  = "HC.Kafka"
    kafka = jsonencode({
      "addr": var.kafka_connect_address,
      "group": var.kafka_consumer_group_id,
      "instance_name": var.kafka_instance_name,
      "instance_id": var.kafka_instance_id,
      "topic": var.kafka_topic_name,
      "seek_to": "latest",
      "security_protocol": "PLAINTEXT",
    })
  }
  sink {
    name  = "HC.Kafka"
    kafka = jsonencode({
      "connection_id": var.eg_connection_id,
      "topic": var.kafka_topic_name,
      "key_transform": {
        "type": "ORIGINAL"
      }
    })
  }
  rule_config {
    transform {
      type = "ORIGINAL"
    }
  }
  option {
    thread_num = 2

    batch_window {
      count    = 5
      time     = 3
      interval = 2
    }
  }
}
```

### Create an event stream with event source type RocketMQ

```hcl
variable "stream_name" {}
variable "rocketmq_instance_id" {}
variable "rocketmq_consumer_group_name" {}
variable "rocketmq_topic_name" {}
variable "function_urn" {}
# The agency must be authorized to the EG service and include the following permissions:
# + EG Publisher
# + FunctionGraph CommonOperations
# + SMN Administrator
# + functiongraph:function:invoke*
variable "agency_name" {}

resource "huaweicloud_eg_event_stream" "test" {
  name   = var.stream_name
  action = "START"

  source {
    name         = "HC.DMS_ROCKETMQ"
    dms_rocketmq = jsonencode({
      "instance_id": var.rocketmq_instance_id,
      "group": var.rocketmq_consumer_group_name,
      "topic": var.rocketmq_topic_name,
      "tag": "lance",
      "access_key": "custom_user_name",
      "secret_key": "User!Password",
      "ssl_enable": false,
      "enable_acl": true,
      "message_type": "NORMAL",
      "consume_timeout": 30000,
      "consumer_thread_nums": 20,
      "consumer_batch_max_size": 2
    })
  }
  sink {
    name          = "HC.FunctionGraph"
    functiongraph = jsonencode({
      "urn": var.function_urn,
      "agency": var.agency_name,
      "invoke_type": "ASYNC"
    })
  }
  rule_config {
    transform {
      type = "ORIGINAL"
    }
  }
  option {
    thread_num = 2

    batch_window {
      count    = 5
      time     = 3
      interval = 2
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region where the event stream is located.  
  If omitted, the provider-level region will be used.
  Changing this will create a new resource.

* `name` - (Required, String) Specifies the name of the event stream.  
  The valid length is limited from `1` to `128`, only letters, digits, hyphens (-), underscores (_) and dots (.) are
  allowed. The name must start with a letter or digit.

* `source` - (Required, List) Specifies the source configuration of the event stream.
  The [source](#stream_source) structure is documented below.

* `sink` - (Required, List) Specifies the target configuration of the event stream.
  The [sink](#stream_sink) structure is documented below.

* `rule_config` - (Required, List) Specifies the rule configuration of the event stream.
  The [rule_config](#stream_rule_config) structure is documented below.

* `option` - (Required, List) Specifies the runtime configuration of the event stream.
  The [targets](#stream_option) structure is documented below.

* `description` - (Optional, String) Specifies the description of the event stream.

* `action` - (Optional, String) Specifies the desired running status of the event stream.
  + **START**
  + **PAUSE**

  Defaults to **PAUSE**.

  ~> The security group used by Kafka or RocketMQ must open the port required by its service documentation display in
  the ingress direction so that the EG service can create resources for it. If not, event stream will fail to start.

<a name="stream_source"></a>
The `source` block supports:

* `name` - (Required, String, ForceNew) Specifies the name of the event source type.  
  The valid values are as follows:
  + **HC.Kafka**
  + **HC.COMMUNITY_ROCKETMQ**
  + **HC.DMS_ROCKETMQ**
  + **HC.MOBILE_ROCKETMQ**

  Changing this will create a new resource.

* `kafka` - (Optional, String) Specifies the event source configuration detail for DMS Kafka type, in JSON format.

* `mobile_rocketmq` - (Optional, String) Specifies the event source configuration detail for mobile RocketMQ type, in
  JSON format.

* `community_rocketmq` - (Optional, String) Specifies the event source configuration detail for community RocketMQ type,
  in JSON format.

* `dms_rocketmq` - (Optional, String) Specifies the event source configuration detail for DMS RocketMQ type, in JSON
  format.

-> Exactly one of `kafka`, `mobile_rocketmq`, `community_rocketmq` and `dms_rocketmq` must be provided.

<a name="stream_sink"></a>
The `sink` block supports:

* `name` - (Required, String) Specifies the name of the event target type.  
  The valid values are as follows:
  + **HC.FunctionGraph**
  + **HC.Kafka**

* `functiongraph` - (Optional, String) Specifies the event target configuration detail for FunctionGraph type, in JSON
  format.

* `kafka` - (Optional, String) Specifies the event target configuration detail for DMS Kafka type, in JSON format.

-> Exactly one of `functiongraph` and `kafka` must be provided.

<a name="stream_rule_config"></a>
The `rule_config` block supports:

* `transform` - (Required, List) Specifies the configuration detail of the transform rule.  
  The [transform](#stream_rule_config_transform) structure is documented below.

* `filter` - (Optional, String) Specifies the configuration detail of the filter rule, in JSON format.

<a name="stream_rule_config_transform"></a>
The `transform` block supports:

* `type` - (Required, String) Specifies the type of transform rule.  
  The valid values are as follows:
  + **ORIGINAL**
  + **CONSTANT**
  + **VARIABLE**

* `value` - (Optional, String) Specifies the rule content definition.
  + When the constant type rule is used, the field is a constant content definition
  + when the variable type rule is used, it is a variable definition, and the content must be a JSON object string.
    - A maximum of `100` variables are supported, and nested structure definitions are not supported.
    - Variable names are composed of letters, numbers, dots, underscores, and dashes. They must start with a letter or
      number and cannot start with `HC.`, and the length should not exceed `64` characters.
    - variable value Expressions support constants or JsonPath expressions, and the string length does not exceed
      `1,024` characters.

* `template` - (Optional, String) Specifies the template definition of the rule content.  
  It's only valid for variable type rules and supports references to defined variables.  
  The string length does not exceed `2,048` characters.

<a name="stream_option"></a>
The `option` block supports:

* `thread_num` - (Required, Int) Specifies the number of concurrent threads.

* `batch_window` - (Required, List) Specifies the configuration of the batch push.
  The [transform](#stream_option_batch_window) structure is documented below.

<a name="stream_option_batch_window"></a>
The `batch_window` block supports:

* `count` - (Required, Int) Specifies the number of items pushed in batches.  
  The valid value is range from `1` to `10,000`.

* `time` - (Optional, Int) Specifies the number of retries.

* `interval` - (Optional, Int) Specifies the interval of the batch push.  
  The valid value is range from `1` to `15`.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - A resource ID in UUID format.

* `status` - The status of the event stream.

* `created_at` - The (UTC) creation time of the event stream, in RFC3339 format.

* `updated_at` - The (UTC) update time of the event stream, in RFC3339 format.

## Import

Streams can be imported using their `id`, e.g.

```bash
$ terraform import huaweicloud_eg_event_stream.test <id>
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response, security or some other reason.
The missing attributes include: `source`, `sink`.
It is generally recommended running `terraform plan` after importing a stream.
You can then decide if changes should be applied to the stream, or the resource definition should be updated to
align with the stream. Also you can ignore changes as below.

```hcl
resource "huaweicloud_eg_event_stream" "test" {
  ...

  lifecycle {
    ignore_changes = [
      source, sink,
    ]
  }
}
```
