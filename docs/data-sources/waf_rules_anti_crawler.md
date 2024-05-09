---
subcategory: "Web Application Firewall (WAF)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_waf_rules_anti_crawler"
description: |-
  Use this data source to get a list of anti crawler rules.
---

# huaweicloud_waf_rules_anti_crawler

Use this data source to get a list of anti crawler rules.

## Example Usage

```hcl
variable "policy_id" {}
variable "rule_id" {}

data "huaweicloud_waf_rules_anti_crawler" "test" {
  policy_id = var.policy_id
  rule_id   = var.rule_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `policy_id` - (Required, String) Specifies the ID of the policy to which the anti crawler rules belong.

* `rule_id` - (Optional, String) Specifies the ID of the anti crawler rule.

* `name` - (Optional, String) Specifies the name of the anti crawler rule.

* `protection_mode` - (Optional, String) Specifies the protection mode of the anti crawler rule.
  The valid values are as follows:
  + **anticrawler_except_url**: All paths are protected except the one specified in the queried anti crawler rules.
  + **anticrawler_specific_url**: The specified path is protected in the queried anti crawler rules.
  
  -> If omitted, the API default query the anti crawler rules in **anticrawler_except_url** protection mode.

* `status` - (Optional, String) Specifies the status of the anti crawler rule.
  The valid values are as follows:
  + **0**: The anti crawler rule is disabled.
  + **1**: The anti crawler rule is active.

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project ID to which the protection policy belongs.
  If omitted, will query the anti crawler rules under the default enterprise project for enterprise users.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `rules` - The list of anti crawler rules.

  The [rules](#rules_struct) structure is documented below.

<a name="rules_struct"></a>
The `rules` block supports:

* `id` - The ID of the anti crawler rule.

* `policy_id` - The ID of the policy to which the anti crawler rule belongs.

* `name` - The name of the anti crawler rule.

* `protection_mode` - The protection mode of the anti crawler rule.

* `status` - The status of the anti crawler rule.

* `description` - The description of the anti crawler rule.

* `priority` - The priority of the anti crawler rule.

* `conditions` - The matching condition list of the anti crawler rule.

  The [conditions](#rules_conditions_struct) structure is documented below.

* `created_at` - The creation time of the anti crawler rule.

<a name="rules_conditions_struct"></a>
The `conditions` block supports:

* `field` - The field type of the condition.

* `logic` - The condition matching logic.

* `content` - The content of the match condition.

* `reference_table_id` - The reference table ID.
