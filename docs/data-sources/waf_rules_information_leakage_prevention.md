---
subcategory: "Web Application Firewall (WAF)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_waf_rules_information_leakage_prevention"
description: |-
  Use this data source to get a list of information leakage prevention rules.
---

# huaweicloud_waf_rules_information_leakage_prevention

Use this data source to get a list of information leakage prevention rules.

## Example Usage

```hcl
variable "policy_id" {}
variable "rule_id" {}

data "huaweicloud_waf_rules_information_leakage_prevention" "test" {
  policy_id = var.policy_id
  rule_id   = var.rule_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `policy_id` - (Required, String) Specifies the ID of the policy to which the information leakage prevention rules belong.

* `rule_id` - (Optional, String) Specifies the ID of the information leakage prevention rule.

* `status` - (Optional, String) Specifies the status of the information leakage prevention rule.
  The value can be **0** or **1**.
  + **0**: The rule is disabled.
  + **1**: The rule is enabled.

* `type` - (Optional, String) Specifies the type of the information leakage prevention rule.
  The value can be **code** for response code or **sensitive** for sensitive information.

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project ID to which the protection policy belongs.
  If omitted, will query the information leakage prevention rules under the default enterprise project for enterprise users.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `rules` - The list of information leakage prevention rules.

  The [rules](#rules_struct) structure is documented below.

<a name="rules_struct"></a>
The `rules` block supports:

* `id` - The ID of the information leakage prevention rule.

* `policy_id` - The ID of the policy to which the information leakage prevention rule belongs.

* `status` - The status of the information leakage prevention rule.

* `description` - The description of the information leakage prevention rule.

* `path` - The path to which the information leakage prevention rule applies.

* `type` - The type of the information leakage prevention rule.

* `contents` - The contents of the information leakage prevention rule.

* `protection_action` - The protection action of the information leakage prevention rule.

* `created_at` - The creation time of the information leakage prevention rule.
