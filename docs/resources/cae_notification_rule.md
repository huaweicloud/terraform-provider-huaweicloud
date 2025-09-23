---
subcategory: "Cloud Application Engine (CAE)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cae_notification_rule"
description: |-
  Manage an event notification rule resource within HuaweiCloud.
---

# huaweicloud_cae_notification_rule

Manage an event notification rule resource within HuaweiCloud.

## Example Usage

```hcl
variable "notification_rule_name" {}
variable "notification_event_name" {}
variable "application_ids" {
  type = list(string)
}
variable "notification_email" {}

resource "huaweicloud_cae_notification_rule" "test" {
  name       = var.notification_rule_name
  event_name = var.notification_event_name

  scope {
    type         = "applications"
    applications = var.application_ids
  }

  trigger_policy {
    type = "immediately"
  }

  notification {
    protocol = "email"
    endpoint = var.notification_email
    template = "EN"
  }

  enabled = true
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.
  Changing this creates a new resource.

* `name` - (Required, String, ForceNew) Specifies the name of the event notification rule. The name must be unique.  
  Changing this creates a new resource.  
  The valid length is limited from `1` to `64`, only English letters, digits, underscores (_) and hyphens (-) are
  allowed.  
  The name must start and end with an English letter or a digit.
  
* `event_name` - (Required, String) Specifies the trigger event of the event notification.  
  Multiple events are separated by commas (,). e.g. **Healthy,Pulled**.  
  The valid values are as follows:
  + **Healthy**: Healthy checked.
  + **Unhealthy**: Healthy check failed.
  + **Pulled**: Image pulled.
  + **FailedPullImage**: Pull image failed.
  + **Started**: Container started up.
  + **BackOffStart**: Container startup failed.
  + **SuccessfulMountVolume**: Volume mounted.
  + **FailedMount**: Attach volume failed.

* `scope` - (Required, List) Specifies the scope in which event notification rule takes effect.  
  The [scope](#notification_rule_scope) structure is documented below.

* `trigger_policy` - (Required, List) Specifies the trigger policy of the event notification rule.  
  The [trigger_policy](#notification_rule_trigger_policy) structure is documented below.

* `notification` - (Required, List, ForceNew) Specifies the configuration of the event notification.
  Changing this creates a new resource.  
  The [notification](#notification_rule_notification) structure is documented below.

* `enabled` - (Optional, Bool) Specifies whether to enable the event notification rule. Defaults to **false**.

<a name="notification_rule_notification"></a>
The `notification` block supports:

* `protocol` - (Required, String, ForceNew) Specifies the protocol of the event notification.
  Changing this creates a new resource.  
  The valid values are as follows:
  + **sms**
  + **email**
  + **wechat**

* `endpoint` - (Required, String, ForceNew) Specifies the endpoint of the event notification.
  Changing this creates a new resource.
  + If `notification.type` is set to **sms**, the endpoint is a phone number.
  + If `notification.type` is set to **email**, the endpoint is a email address.
  + If `notification.type` is set to **wechat**, the endpoint is a webhook address starting with
    `https://qyapi.weixin.qq.com/cgi-bin/webhook/send`.
    you want to use this parameter, please refer to the [documentation](https://support.huaweicloud.com/intl/en-us/usermanual-ticket/topic_0065264094.html)
    to submit a service ticket to apply for it.
    For details about how to obtain a WeCom subscription endpoint, please refer to the [documentation](https://support.huaweicloud.com/intl/en-us/smn_faq/smn_faq_0027.html).
  
* `template` - (Required, String, ForceNew) Specifies the template language of the event notification.
  Changing this creates a new resource.  
  The valid values are as follows:
  + **EN**
  + **ZH**

<a name="notification_rule_scope"></a>
The `scope` block supports:

* `type` - (Required, String) Specifies the type to which the event notification rule takes effect.  
  The valid values are as follows:
  + **environments**: The rule takes effect for all components in the environment.
  + **applications**: The rule takes effect for all components in the application.
  + **components**: The rule takes effect for the specified components.

* `environments` - (Optional, List) Specifies the list of the environment IDs.  
  This parameter is required and available only when the `scope.type` parameter is set to **environments**.

* `applications` - (Optional, List) Specifies the list of the applications IDs.  
  This parameter is required and available only when the `scope.type` parameter is set to **applications**.

* `components` - (Optional, List) Specifies the list of the components IDs.  
  This parameter is required and available only when the `scope.type` parameter is set to **components**.

<a name="notification_rule_trigger_policy"></a>
The `trigger_policy` block supports:

* `type` - (Required, String) Specifies the type of the trigger.  
  The valid values are as follows:
  + **accumulative**
  + **immediately**

* `period` - (Optional, Int) Specifies the trigger period of the event. The unit is second.  
  This parameter is required and available only when the `trigger_policy.type` parameter is set to **accumulative**.  
  The valid values are as follows:
  + **300**
  + **1200**
  + **3600**
  + **14400**
  + **86400**

* `count` - (Optional, Int) Specifies the number of times the event occurred.  
  The valid value ranges from `1` to `100`.  
  This parameter is required and available only when the `trigger_policy.type` parameter is set to **accumulative**.

* `operator` - (Optional, String) Specifies the condition of the event notification.  
  The valid values are **>** and **>=**.  
  This parameter is required and available only when the `trigger_policy.type` parameter is set to **accumulative**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID, also event notification rule ID.

## Import

The event notification rule resource can be imported using `name`, e.g.

```bash
$ terraform import huaweicloud_cae_notification_rule.test <name>
```
