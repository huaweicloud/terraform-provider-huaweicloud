---
subcategory: "EventGrid (EG)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_eg_event_subscriptions"
description: |-
  Use this data source to get the list of EG event subscriptions within HuaweiCloud.
---

# huaweicloud_eg_event_subscriptions

Use this data source to get the list of EG event subscriptions within HuaweiCloud.

## Example Usage

### Query all subscriptions

```hcl
data "huaweicloud_eg_event_subscriptions" "test" {}
```

### Query subscriptions by name

```hcl
variable "subscription_name" {}

data "huaweicloud_eg_event_subscriptions" "test" {
  name = var.subscription_name
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region where the subscriptions are located.  
  If omitted, the provider-level region will be used.

* `channel_id` - (Optional, String) Specifies the ID of the event channel to filter subscriptions.

* `name` - (Optional, String) Specifies the exact name of the subscription to be queried.

* `fuzzy_name` - (Optional, String) Specifies the name of the subscription to be queried for fuzzy matching.

* `connection_id` - (Optional, String) Specifies the ID of the target connection to filter subscriptions.

* `sort` - (Optional, String) Specifies the sorting method for query results.  
  The format is `field:order`, where `field` is the field name and `order` is `ASC` or `DESC`. e.g. `created_time:DESC`.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `subscriptions` - The list of the subscriptions that matched filter parameters.  
  The [subscriptions](#eg_subscriptions_attr) structure is documented below.

<a name="eg_subscriptions_attr"></a>
The `subscriptions` block supports:

* `id` - The ID of the subscription.

* `name` - The name of the subscription.

* `description` - The description of the subscription.

* `type` - The type of the subscription.
  + **EVENT**
  + **SCHEDULED**

* `status` - The status of the subscription.
  + **CREATED**
  + **ENABLED**
  + **DISABLED**
  + **FROZEN**
  + **ERROR**

* `channel_id` - The ID of the event channel.

* `channel_name` - The name of the event channel.

* `used` - The associated resources of the subscription.  
  The [used](#eg_subscriptions_used) structure is documented below.

* `sources` - The list of subscription sources.  
  The [sources](#eg_subscriptions_sources) structure is documented below.

* `targets` - The list of subscription targets.  
  The [targets](#eg_subscriptions_targets) structure is documented below.

* `created_time` - The creation time of the subscription, in RFC3339 format.

* `updated_time` - The update time of the subscription, in RFC3339 format.

<a name="eg_subscriptions_used"></a>
The `used` block supports:

* `resource_id` - The ID of the associated resource.

* `owner` - The management tenant account to which the associated resource belongs.

* `description` - The description of the associated resource.

<a name="eg_subscriptions_sources"></a>
The `sources` block supports:

* `id` - The ID of the subscription source.

* `name` - The name of the subscription source.

* `provider_type` - The provider type of the subscription source.
  + **OFFICIAL**
  + **CUSTOM**
  + **PARTNER**

* `detail` - The parameter list of the subscription source, in JSON format.

* `filter` - The matching filter rules of the subscription source, in JSON format.

* `created_time` - The creation time of the subscription source, in RFC3339 format.

* `updated_time` - The update time of the subscription source, in RFC3339 format.

<a name="eg_subscriptions_targets"></a>
The `targets` block supports:

* `id` - The ID of the subscription target.

* `name` - The name of the subscription target.

* `provider_type` - The provider type of the subscription target.
  + **OFFICIAL**
  + **CUSTOM**
  + **PARTNER**

* `connection_id` - The target connection ID used by the subscription target.

* `detail` - The parameter list of the subscription target, in JSON format.

* `kafka_detail` - The Kafka target parameter list of the subscription, in JSON format.

* `smn_detail` - The SMN target parameter list of the subscription, in JSON format.

* `eg_detail` - The EG channel target parameter list of the subscription, in JSON format.

* `apigw_detail` - The APIGW target parameter list of the subscription, in JSON format.

* `retry_times` - The number of retry times.

* `transform` - The transform rules of the subscription target, in JSON format.

* `dead_letter_queue` - The dead letter queue parameters of the subscription, in JSON format.

* `created_time` - The creation time of the subscription target, in RFC3339 format.

* `updated_time` - The update time of the subscription target, in RFC3339 format.
