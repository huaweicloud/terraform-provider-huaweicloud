---
subcategory: "Application Operations Management (AOM 2.0)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_aom_alarm_action_rule"
description: ""
---

# huaweicloud_aom_alarm_action_rule

Manages an AOM alarm action rule resource within HuaweiCloud.

## Example Usage

```hcl
variable "topic_urn" {}

resource "huaweicloud_aom_alarm_action_rule" "test" {
  name                  = "test_rule"
  description           = "terraform test"
  type                  = "1"
  notification_template = "aom.built-in.template.zh"

  smn_topics {
    topic_urn = var.topic_urn
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `name` - (Required, String, ForceNew) Specifies the action rule name. The value can be a string of 1 to 100
  characters that can consist of letters, digits, underscores (_), hyphens (-) and Chinese characters,
  and it must start and end with letters, digits or Chinese characters.

  Changing this parameter will create a new resource.

* `type` - (Required, String) Specifies the action rule type. The value can be **1**, which indicates notification.

* `smn_topics` - (Required, List) Specifies the SMN topic configurations. A maximum of 5 topics are allowed.
  The [SmnTopics](#AlarmActionRule_SmnTopics) structure is documented below.

* `notification_template` - (Required, String) Specifies the notification template.

* `description` - (Optional, String) Specifies the action rule description.
  The value can be a string of 0 to 1024 characters.

<a name="AlarmActionRule_SmnTopics"></a>
The `SmnTopics` block supports:

* `topic_urn` - (Required, String) Specifies the SMN topic URN.

* `name` - (Optional, String) Specifies the SMN topic name.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID. The value is the rule name.

* `created_at` - The creation time.

* `updated_at` - The last update time.

## Import

The application operations management can be imported using the `id` (name), e.g.

```bash
$ terraform import huaweicloud_aom_alarm_action_rule.test test_rule
```
