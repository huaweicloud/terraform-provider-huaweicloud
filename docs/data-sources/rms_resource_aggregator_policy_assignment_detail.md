---
subcategory: "Config"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_rms_resource_aggregator_policy_assignment_detail"
description: |-
  Use this data source to get the detail about a specified aggregated rule.
---

# huaweicloud_rms_resource_aggregator_policy_assignment_detail

Use this data source to get the detail about a specified aggregated rule.

## Example Usage

```hcl
variable "aggregator_id" {}
variable "account_id" {}
variable "policy_assignment_id" {}

data "huaweicloud_rms_resource_aggregator_policy_assignment_detail" "test" {
  aggregator_id        = var.aggregator_id
  account_id           = var.account_id
  policy_assignment_id = var.policy_assignment_id
}
```

## Argument Reference

The following arguments are supported:

* `aggregator_id` - (Required, String) Specifies the resource aggregator ID.

* `account_id` - (Required, String) Specifies the source account ID.

* `policy_assignment_id` - (Required, String) Specifies the rule ID.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `policy_assignment_type` - Indicates the rule type, which can be builtin or custom.

* `name` - Indicates the rule name.

* `description` - Indicates the rule description.

* `policy_filter` - Indicates the policy filter of a rule.
  The [policy_filter](#policy_filter_struct) structure is documented below.

* `policy_filter_v2` - Indicates the policy filter of a rule.
  The [policy_filter_v2](#policy_filter_v2_struct) structure is documented below.

* `period` - Indicates how often the rule is triggered, which can be **One_Hour**, **Three_Hours**, **Six_Hours**,
  **Twelve_Hours**, or **TwentyFour_Hours**.

* `state` - Indicates the rule status.

* `created` - Indicates the time when the rule was added.

* `updated` - Indicates the time when the rule was modified.

* `policy_definition_id` - Indicates the ID of the policy associated with a rule.

* `custom_policy` - Indicates the custom rule.
  The [custom_policy](#custom_policy_struct) structure is documented below.

* `parameters` -  Indicates rule parameters.

* `tags` - Indicates the tags.
  The [tags](#tags_struct) structure is documented below.

* `created_by` - Indicates the rule creator.

* `target_type` - Indicates the execution method of remediation.

* `target_id` - Indicates the ID of a remediation object.

<a name="policy_filter_struct"></a>
The `policy_filter` block supports:

* `region_id` - Indicates the region ID.

* `resource_provider` - Indicates the cloud service name.

* `resource_type` - Indicates the resource type.

* `resource_id` - Indicates the resource ID.

* `tag_key` - Indicates the tag key.

* `tag_value` - Indicates the tag value.

<a name="policy_filter_v2_struct"></a>
The `policy_filter_v2` block supports:

* `region_ids` -  Indicates the region IDs.

* `resource_types` - Indicates the cloud services.

* `resource_ids` - Indicates the resource list.

* `tag_key_logic` - Indicates the logical relationship when parameter `tags` takes multiple values, for example: When the
  `tags` is **"tags.1.key":"a", "tags.1.values":"a", "tags.2.key":"b", "tags.2.values":"b"**, if this parameter is set to
  **AND**, it means that the rule only applies to resources bound with both tags **a:a** and **b:b**. If not specified,
  the default logic is **OR**.

* `tags` - Indicates the tags.
  The [tags](#policy_filter_v2_tags_struct) structure is documented below.

* `exclude_tag_key_logic` - Indicates the logical relationship when parameter `exclude_tags` takes multiple values, for
  example: When the `exclude_tags` is **"exclude_tags.1.key":"a", "exclude_tags.1.values":"a", "exclude_tags.2.key":"b",
  "exclude_tags.2.values":"b"**, if this parameter is set to **AND**, it means that the rule excludes resources that are
  bound with the tags **a:a** and **b:b**. If not specified, the default logic is **OR**.

* `exclude_tags` - Indicates the exclude tags.
  The [exclude_tags](#policy_filter_v2_tags_struct) structure is documented below.

<a name="policy_filter_v2_tags_struct"></a>
The `tags` and `exclude_tags` block supports:

* `key` - Indicates the tag key.

* `values` - Indicates the tag values.

<a name="custom_policy_struct"></a>
The `custom_policy` block supports:

* `function_urn` - Indicates the URN of a custom function.

* `auth_type` - Indicates the method used by a custom rule to call a function.

* `auth_value` - Indicates the value of the method used by a custom rule to call a function.

<a name="tags_struct"></a>
The `tags` block supports:

* `key` - Indicates the tag key.

* `value` - Indicates the tag value.
