---
subcategory: "Application Operations Management (AOM)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_aom_alarm_group_rules"
description: |-
  Use this data source to get the list of AOM alarm group rules.
---

# huaweicloud_aom_alarm_group_rules

Use this data source to get the list of AOM alarm group rules.

## Example Usage

```hcl
data "huaweicloud_aom_alarm_group_rules" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project ID to which the rules belong.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `rules` - The rules list.
  The [rules](#attrblock--rules) structure is documented below.

<a name="attrblock--rules"></a>
The `rules` block supports:

* `name` - Indicates the alarm group rule name.

* `detail` - Indicates the grouping conditions list.
  The [detail](#attrblock--rules--detail) structure is documented below.

* `group_by` - Indicates the combine notifications.

* `group_wait` - Indicates the initial wait time.

* `group_interval` - Indicates the batch processing interval.

* `group_repeat_waiting` - Indicates the repeat interval.

* `description` - Indicates the alarm group rule description.

* `enterprise_project_id` - Indicates the enterprise project ID to which the rule belongs.

* `created_at` - Indicates the rule create time.

* `updated_at` - Indicates the rule update time.

<a name="attrblock--rules--detail"></a>
The `detail` block supports:

* `bind_notification_rule_ids` - Indicates the action rule IDs.

* `match` - Indicates the matching conditions list.
  The [match](#attrblock--rules--detail--match) structure is documented below.

<a name="attrblock--rules--detail--match"></a>
The `match` block supports:

* `key` - Indicates the matching condition key.

* `operate` - Indicates the matching condition operator.

* `value` - Indicates the matching condition value.
