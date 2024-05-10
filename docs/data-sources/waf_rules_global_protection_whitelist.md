---
subcategory: "Web Application Firewall (WAF)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_waf_rules_global_protection_whitelist"
description: |-
  Use this data source to get a list of global protection whitelist rules.
---

# huaweicloud_waf_rules_global_protection_whitelist

Use this data source to get a list of global protection whitelist rules.

## Example Usage

```hcl
variable "policy_id" {}
variable "rule_id" {}

data "huaweicloud_waf_rules_global_protection_whitelist" "test" {
  policy_id = var.policy_id
  rule_id   = var.rule_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `policy_id` - (Required, String) Specifies the ID of the policy to which the global protection whitelist rules belong.

* `rule_id` - (Optional, String) Specifies the ID of the global protection whitelist rule.

* `status` - (Optional, String) Specifies the status of the global protection whitelist rule.
  The valid values are as follows:
  + **0**: The global protection whitelist rule is disabled.
  + **1**: The global protection whitelist rule is active.

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project ID to which the protection policy belongs.
  If omitted, will query the global protection whitelist rules under the default enterprise project for enterprise users.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `rules` - The list of the global protection whitelist rules.

  The [rules](#rules_struct) structure is documented below.

<a name="rules_struct"></a>
The `rules` block supports:

* `id` - The ID of the global protection whitelist rule.

* `policy_id` - The ID of the policy to which the global protection whitelist rule belongs.

* `status` - The status of the global protection whitelist rule.

* `description` - The description of the global protection whitelist rule.

* `ignore_waf_protection` - The rules that need to be ignored.

* `conditions` - The matching condition list of the global protection whitelist rule.

  The [conditions](#rules_conditions_struct) structure is documented below.

* `domains` - The protected domain name or website bound with the policy.

* `advanced_field` - The filed type of the advanced configuration.

* `advanced_content` - The subfiled of the advanced configuration.

* `created_at` - The creation time of the global protection whitelist rule.

<a name="rules_conditions_struct"></a>
The `conditions` block supports:

* `field` - The field type of the condition.

* `subfield` - The subfield of the condition.

* `logic` - The condition matching logic.

* `content` - The content of the match condition.
