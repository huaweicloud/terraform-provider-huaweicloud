---
subcategory: "Web Application Firewall (WAF)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_waf_rules_precise_protection"
description: |-
  Use this data source to get a list of precise protection rules.
---

# huaweicloud_waf_rules_precise_protection

Use this data source to get a list of precise protection rules.

## Example Usage

```hcl
variable "policy_id" {}
variable "rule_id" {}

data "huaweicloud_waf_rules_precise_protection" "test" {
  policy_id = var.policy_id
  rule_id   = var.rule_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `policy_id` - (Required, String) Specifies the ID of the policy to which the precise protection rules belong.

* `rule_id` - (Optional, String) Specifies the ID of the precise protection rule.

* `name` - (Optional, String) Specifies the name of the precise protection rule.

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project ID to which the protection policies belong.
  If omitted, will query precise protection rules under the default enterprise project for enterprise users.

* `status` - (Optional, String) Specifies the status of the precise protection rule.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `rules` - The list of the precise protection rules.

  The [rules](#rules_struct) structure is documented below.

<a name="rules_struct"></a>
The `rules` block supports:

* `id` - The ID of the precise protection rule.

* `policy_id` - The ID of the policy to which the precise protection rule belongs.

* `name` - The name of the precise protection rule.

* `description` - The description of the precise protection rule.

* `status` - The status of the precise protection rule.

* `priority` - The priority of the precise protection rule.

* `time` - The effective time of the precision protection rule.

* `start_time` - The start time for the implementation of precision protection rule.
  This parameter will only be returned when the `time` value is **true**.

* `end_time` - The end time for the implementation of precision protection rule.
  This parameter will only be returned when the `time` value is **true**.

* `conditions` - The matching condition list of the precision protection rule.

  The [conditions](#rules_conditions_struct) structure is documented below.

* `action` - The protective action of the precise protection rule.

* `known_attack_source_id` - The known attack source ID.

* `created_at` - The creation time of the precise protection rule.

<a name="rules_conditions_struct"></a>
The `conditions` block supports:

* `field` - The field of the condition.

* `subfield` - The subfield of the condition.

* `logic` - The condition matching logic.

* `content` - The content of the match condition.

* `reference_table_id` - The reference table id.
