---
subcategory: "EventGrid (EG)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_eg_event_subscription_target"
description: |-
  Using this resource to manage EG event subscription target within Huaweicloud.
---

# huaweicloud_eg_event_subscription_target

Using this resource to manage EG event subscription target within Huaweicloud.

-> Manage this resource will cause the `targets` param of the `huaweicloud_eg_event_subscription` resource has
   change plan. if both resources are managed together, please add `ignore_changes = [targets]`
   in `huaweicloud_eg_event_subscription` to suppress this diff.

## Example Usage

```hcl
variable "target_name" {}
variable "resource_urn" {}
variable "subscription_id" {}

resource "huaweicloud_eg_event_subscription_target" "test" {
  subscription_id = var.subscription_id
  name            = var.target_name
  provider_type   = "OFFICIAL"
  
  key_transform {
    type = "ORIGINAL"
  }

  detail {
    urn         = var.resource_urn
    agency_name = "EG_TARGET_AGENCY"
    invoke_type = "SYNC"
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region where the event subscription target is located.  
  If omitted, the provider-level region will be used.
  Changing this parameter will create a new resource.

* `subscription_id` - (Required, String, NonUpdatable) Specifies the ID of the event subscription.

* `name` - (Required, String, NonUpdatable) Specifies the name of the event subscription target.  
  The valid values are as follows:
  + **HTTPS**
  + **HC.APIGW**
  + **HC.EG**
  + **HC.FunctionGraph**
  + **HC.Kafka**
  + **HC.OPENAPI**
  + **HC.SMN**

* `provider_type` - (Required, String, NonUpdatable) Specifies the provider type of the event subscription target.  
  if the prefix of `name` is **HC**, then `provider_type` should be set to **OFFICIAL**.
  The valid values are as follows:
  + **OFFICIAL**
  + **CUSTOM**

* `key_transform` - (Required, List) Specifies the transform configuration for event data transformation.  
  The [key_transform](#eg_target_key_transform) structure is documented below.

* `connection_id` - (Optional, String) Specifies the connection ID used by the event subscription target.

* `enterprise_project_id` - (Optional, String) Specifies the ID of enterprise project.  
   This parameter is only valid for enterprise user.

* `detail` - (Optional, String) Specifies the configuration details of the event subscription target, in JSON format.  
  Required if the `name` is **HTTPS**, **HC.FunctionGraph** and **HC.OPENAPI**.

* `kafka_detail` - (Optional, List) Specifies the Kafka target configuration details.  
  The [kafka_detail](#eg_target_kafka_detail) structure is documented below.
  Required if the `name` is **HC.Kafka**.

* `smn_detail` - (Optional, List) Specifies the SMN target configuration details.  
  The [smn_detail](#eg_target_smn_detail) structure is documented below.
  Required if the `name` is **HC.SMN**.

* `eg_detail` - (Optional, List) Specifies the EG channel target configuration details.  
  The [eg_detail](#eg_target_eg_detail) structure is documented below.
  Required if the `name` is **HC.EG**.

* `apigw_detail` - (Optional, List) Specifies the APIGW target configuration details.  
  The [apigw_detail](#eg_target_apigw_detail) structure is documented below.
  Required if the `name` is **HC.APIGW**.

* `retry_times` - (Optional, Int) Specifies the number of retry times for the event subscription target.

* `dead_letter_queue` - (Optional, List) Specifies the dead letter queue configuration of the event subscription
 target.  
  The [dead_letter_queue](#eg_target_dead_letter_queue) structure is documented below.

<a name="eg_target_key_transform"></a>
The `key_transform` block supports:

* `type` - (Required, String) Specifies the type of transform rule.  
  The valid values are as follows:
  + **ORIGINAL**: passthrough variable.
  + **CONSTANT**: pass constant to target.
  + **VARIABLE**: pass variable to target.

* `value` - (Optional, String) Specifies the value of the transform rule.

* `template` - (Optional, String) Specifies the template definition for VARIABLE type transform rules.

<a name="eg_target_kafka_detail"></a>
The `kafka_detail` block supports:

* `topic` - (Required, String) Specifies the topic name of the Kafka instance.

* `key_transform` - (Optional, List) Specifies the transform configuration of the Kafka messages.  
  The [key_transform](#eg_target_key_transform) structure is documented below.

<a name="eg_target_smn_detail"></a>
The `smn_detail` block supports:

* `urn` - (Required, String) Specifies the URN of the SMN topic.

* `agency_name` - (Required, String) Specifies the agency name for cross-account access.

* `key_transform` - (Optional, List) Specifies the subject transform configuration of the Kafka
 messages.  
  The [key_transform](#eg_target_key_transform) structure is documented below.

<a name="eg_target_eg_detail"></a>
The `eg_detail` block supports:

* `target_project_id` - (Required, String) Specifies the target project ID of the EG channel.

* `target_channel_id` - (Required, String) Specifies the target channel ID of the EG channel.

* `target_region` - (Required, String) Specifies the target region of the EG channel.

* `agency_name` - (Required, String) Specifies the agency name for cross-account access.

* `cross_region` - (Optional, Bool) Specifies whether this is a cross-region EG channel target.

* `cross_account` - (Optional, Bool) Specifies whether this is a cross-account EG channel target.

<a name="eg_target_apigw_detail"></a>
The `apigw_detail` block supports:

* `url` - (Required, String) Specifies the URL of the APIGW endpoint.

* `invocation_http_parameters` - (Optional, List) Specifies the HTTP parameters for the APIGW invocation.  
  The [invocation_http_parameters](#eg_target_apigw_parameters) structure is documented below.

<a name="eg_target_apigw_parameters"></a>
The `invocation_http_parameters` block supports:

* `header_parameters` - (Optional, List) Specifies the header parameters for the HTTP request.  
  The [header_parameters](#eg_target_apigw_header_parameter) structure is documented below.

<a name="eg_target_apigw_header_parameter"></a>
The `header_parameters` block supports:

* `key` - (Optional, String) Specifies the key of the header parameter.

* `value` - (Optional, String) Specifies the value of the header parameter.

* `is_value_secret` - (Optional, Bool) Specifies whether the header parameter value is secret.

<a name="eg_target_dead_letter_queue"></a>
The `dead_letter_queue` block supports:

* `type` - (Required, String) Specifies the type of dead letter queue.

* `instance_id` - (Required, String) Specifies the instance ID of the dead letter queue.

* `connection_id` - (Required, String) Specifies the connection ID of the dead letter queue.

* `topic` - (Required, String) Specifies The topic name of the dead letter queue.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `created_at` - The create time, in UTC format.

* `updated_at` - The update time, in UTC format.

## Import

Event subscription target can be imported using `subscription_id` and `id`, separated by a slash (/), e.g.

```bash
$ terraform import huaweicloud_eg_event_subscription_target.test <subscription_id>/<id>
```
