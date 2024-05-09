---
subcategory: "Web Application Firewall (WAF)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_waf_rules_cc_protection"
description: |-
  Use this data source to get a list of cc protection rules.
---

# huaweicloud_waf_rules_cc_protection

Use this data source to get a list of cc protection rules.

## Example Usage

```hcl
variable "policy_id" {}
variable "rule_id" {}

data "huaweicloud_waf_rules_cc_protection" "test" {
  policy_id = var.policy_id
  rule_id   = var.rule_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `policy_id` - (Required, String) Specifies the ID of the policy to which the cc protection rules belong.

* `rule_id` - (Optional, String) Specifies the ID of the cc protection rule.

* `name` - (Optional, String) Specifies the name of the cc protection rule.

* `status` - (Optional, String) Specifies the status of the cc protection rule.
  The valid values are as follows:
  + **0**: The cc protection rule is disabled.
  + **1**: The cc protection  rule is active.

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project ID to which the protection policy belongs.
  If omitted, will query the cc protection rules under the default enterprise project for enterprise users.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `rules` - The list of cc protection rules.

  The [rules](#rules_struct) structure is documented below.

<a name="rules_struct"></a>
The `rules` block supports:

* `id` - The ID of the cc protection rule.

* `policy_id` - The ID of the policy to which the cc protection rule belongs.

* `name` - The name of the cc protection rule.

* `status` - The status of the cc protection rule.

* `description` - The description of the cc protection rule.

* `conditions` - The matching condition list of the cc protection rule.

  The [conditions](#rules_conditions_struct) structure is documented below.

* `action` - The protective action taken when the number of requests reaches the upper limit.

  The [action](#rules_action_struct) structure is documented below.

* `rate_limit_mode` - The rate limit mode.

* `user_identifier` - The user identifier.

* `other_user_identifier` - The other user identifier.

* `limit_num` - The number of requests allowed from a web visitor in a rate limiting period.

* `limit_period` - The rate limiting period.

* `unlock_num` - The allowable frequency.

* `lock_time` - The lock time for resuming normal page access after blocking can be set.

* `request_aggregation` - Whether to enable domain aggregation statistics.

* `all_waf_instances` - Whether to enable global counting.

* `created_at` - The creation time of the cc protection rule.

<a name="rules_conditions_struct"></a>
The `conditions` block supports:

* `field` - The field of the condition.

* `subfield` - The subfield of the condition.

* `logic` - The condition matching logic.

* `content` - The content of the match condition.

* `reference_table_id` - The reference table ID.

<a name="rules_action_struct"></a>
The `action` block supports:

* `protective_action` - The protective action type.

* `detail` - The block page detail information.

  The [detail](#action_detail_struct) structure is documented below.

<a name="action_detail_struct"></a>
The `detail` block supports:

* `block_page_type` - The type of the returned page.

* `page_content` - The content of the returned page.
