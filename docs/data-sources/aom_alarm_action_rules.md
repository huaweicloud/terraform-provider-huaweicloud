---
subcategory: "Application Operations Management (AOM)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_aom_alarm_action_rules"
description: |-
  Use this data source to get the list of AOM alarm action rules.
---

# huaweicloud_aom_alarm_action_rules

Use this data source to get the list of AOM alarm action rules.

## Example Usage

```hcl
data "huaweicloud_aom_alarm_action_rules" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `action_rules` - Indicates the alarm action rules list.

  The [action_rules](#action_rules_struct) structure is documented below.

<a name="action_rules_struct"></a>
The `action_rules` block supports:

* `name` - Indicates the rule name.

* `type` - Indicates the action type.

* `notification_template` - Indicates the message template.

* `description` - Indicates the rule description.

* `created_by` - Indicates the user who created the rule.

* `created_at` - Indicates the create time.

* `updated_at` - Indicates the update time.

* `time_zone` - Indicates the time zone.

* `smn_topics` - Indicates the SMN topics.

  The [smn_topics](#action_rules_smn_topics_struct) structure is documented below.

<a name="action_rules_smn_topics_struct"></a>
The `smn_topics` block supports:

* `topic_urn` - Indicates the unique resource identifier of the topic.

* `name` - Indicates the name of the topic.

* `display_name` - Indicates the topic display name.

* `push_policy` - Indicates the SMN message push policy.

* `status` - Indicates the status of the topic subscriber.
  Value can be:
  + **0**: The topic has been deleted or the subscription list of this topic is empty.
  + **1**: The subscription object is in the subscribed state.
  + **2**: The subscription object is in the unsubscribed or canceled state.
