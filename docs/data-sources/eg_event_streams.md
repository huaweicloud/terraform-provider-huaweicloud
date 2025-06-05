---
subcategory: "EventGrid (EG)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_eg_event_streams"
description: |-
  Use this data source to query EG event streams within HuaweiCloud.
---

# huaweicloud_eg_event_streams

Use this data source to query EG event streams within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_eg_event_streams" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region where the event sources are located.  
  If omitted, the provider-level region will be used.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `event_streams` - List of event streams.  
  The [event_streams](#eg_event_streams_attr) structure is documented below.

<a name="eg_event_streams_attr"></a>
The `streams` block supports:

* `id` - The ID of the event stream.

* `name` - The name of the event stream.

* `status` - The status of the event stream.

* `description` - The description of the event stream.

* `created_time` - The creation time of the event stream, in RFC3339 format.

* `updated_time` - The latest update time of the event stream, in RFC3339 format.

* `source` - The event source configuration.  
  The [source](#eg_event_streams_source) structure is documented below.

* `sink` - The event sink configuration.  
  The [sink](#eg_event_streams_sink) structure is documented below.

* `rule_config` - The configuration of event rules.  
  The [rule_config](#eg_event_streams_rule_config) structure is documented below.

* `option` - The running configuration.  
  The [option](#eg_event_streams_run_option) structure is documented below.

<a name="eg_event_streams_source"></a>
The `source` block supports:

* `name` - The name of the event source type.

* `source_kafka` - The configuration of Kafka event source, in JSON format.

* `source_mobile_rocketmq` - The configuration of mobile cloud RocketMQ event source, in JSON format.

* `source_community_rocketmq` - The configuration of community RocketMQ event source, in JSON format.  

* `source_dms_rocketmq` - The configuration of DMS RocketMQ event source, in JSON format.

<a name="eg_event_streams_sink"></a>
The `sink` block supports:

* `name` - The name of the event sink type.

* `sink_fg` - The configuration of function graph event sink, in JSON format.

* `sink_kafka` - The configuration of Kafka event sink, in JSON format.

* `sink_obs` - The configuration of OBS event sink, in JSON format.

<a name="eg_event_streams_rule_config"></a>
The `rule_config` block supports:

* `transform` - The transformation rules.  
  The [transform](#eg_event_streams_rule_config_transform) structure is documented below.

* `filter` - The filter rules.

<a name="eg_event_streams_rule_config_transform"></a>
The `transform` block supports:

* `type` - The type of transformation rule.

* `value` - The value of transformation rule.

* `template` - The template of transformation rule.

<a name="eg_event_streams_run_option"></a>
The `option` block supports:

* `thread_num` - The number of concurrent threads.

* `batch_window` - The batch push configuration.  
  The [batch_window](#eg_event_streams_run_option_batch_window) structure is documented below.

<a name="eg_event_streams_run_option_batch_window"></a>
The `batch_window` block supports:

* `count` - The number of batch push messages.

* `time` - The number of retries.

* `interval` - The batch push interval in seconds.
