---
subcategory: "Application Operations Management (AOM 2.0)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_aom_alarm_silence_rules"
description: |-
  Use this data source to get the list of alarm silence rules.
---

# huaweicloud_aom_alarm_silence_rules

Use this data source to get the list of alarm silence rules.

## Example Usage

```hcl
data "huaweicloud_aom_alarm_silence_rules" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `rules` - Indicates the alarm silence rules list.
  The [rules](#attrblock--rules) structure is documented below.

<a name="attrblock--rules"></a>
The `rules` block supports:

* `name` - Indicates the rule name.

* `silence_conditions` - Indicates the silence conditions of the rule.
  The [silence_conditions](#attrblock--rules--silence_conditions) structure is documented below.

* `silence_time` - Indicates the silence time of the rule.
  The [silence_time](#attrblock--rules--silence_time) structure is documented below.

* `time_zone` - Indicates the time zone of the rule.

* `description` - Indicates the description of the rule.

* `created_at` - Indicates the create time of the rule.

* `updated_at` - Indicates the update time of the rule.

<a name="attrblock--rules--silence_conditions"></a>
The `silence_conditions` block supports:

* `conditions` - Indicates the serial conditions.
  The [conditions](#attrblock--rules--silence_conditions--conditions) structure is documented below.

<a name="attrblock--rules--silence_conditions--conditions"></a>
The `conditions` block supports:

* `key` - Indicates the key of the match condition.

* `operate` - Indicates the operate of the match condition.

* `value` - Indicates the value list of the match condition.

<a name="attrblock--rules--silence_time"></a>
The `silence_time` block supports:

* `type` - Indicates the effective time type of the rule.

* `starts_at` - Indicates the start time of the rule.

* `ends_at` - Indicates the end time of the rule.

* `scope` - Indicates the silence time of the rule.
