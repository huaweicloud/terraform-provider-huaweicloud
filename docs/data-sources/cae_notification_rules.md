---
subcategory: "Cloud Application Engine (CAE)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cae_notification_rules"
description: |-
  Use this data source to get the list of the event notification rules within HuaweiCloud.
---

# huaweicloud_cae_notification_rules

Use this data source to get the list of the event notification rules within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_cae_notification_rules" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `rules` - The list of the event notification rules.

  The [rules](#rules_struct) structure is documented below.

<a name="rules_struct"></a>
The `rules` block supports:

* `id` - The ID of the event notification rule.

* `name` - The name of the event notification rule.

* `event_name` - The trigger event of the event notification.  
  The valid values are as follows:
  + **Healthy**: Healthy checked.
  + **Unhealthy**: Healthy check failed.
  + **Pulled**: Image pulled.
  + **FailedPullImage**: Pull image failed.
  + **Started**: Container started up.
  + **BackOffStart**: Container startup failed.
  + **SuccessfulMountVolume**: Volume mounted.
  + **FailedMount**: Attach volume failed.

* `scope` - The scope in which event notification rule takes effect.

  The [scope](#rules_scope_struct) structure is documented below.

* `trigger_policy` - The trigger policy of the event notification rule.

  The [trigger_policy](#rules_trigger_policy_struct) structure is documented below.

* `notification` - The configuration of the event notification.

  The [notification](#rules_notification_struct) structure is documented below.

* `enabled` - Whether the event notification rule is enabled.

<a name="rules_scope_struct"></a>
The `scope` block supports:

* `type` - The type to which the event notification rule takes effect.  
  The valid values are as follows:
  + **environments**: The rule takes effect for all components in the environment.
  + **applications**: The rule takes effect for all components in the application.
  + **components**: The rule takes effect for the specified components.

* `environments` - The list of the environment IDs.

* `applications` - The list of the application IDs.

* `components` - The list of the component IDs.

<a name="rules_trigger_policy_struct"></a>
The `trigger_policy` block supports:

* `type` - The type of the trigger.  
  The valid values are as follows:
  + **accumulative**
  + **immediately**

* `period` - The trigger period of the event.  
  The valid values are as follows:
  + **300**
  + **1200**
  + **3600**
  + **14400**
  + **86400**

* `count` - The number of times the event occurred.

* `operator` - The condition of the event notification.  
  The valid values are **>** and **>=**.

<a name="rules_notification_struct"></a>
The `notification` block supports:

* `protocol` - The protocol of the event notification.  
  The valid values are as follows:
  + **sms**
  + **email**
  + **wechat**

* `endpoint` - The endpoint of the event notification.

* `template` - The template language of the event notification.  
  The valid values are as follows:  
  + **EN**
  + **ZH**
