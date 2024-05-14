---
subcategory: "Web Application Firewall (WAF)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_waf_rules_data_masking"
description: |-
  Use this data source to get a list of data masking rules.
---

# huaweicloud_waf_rules_data_masking

Use this data source to get a list of data masking rules.

## Example Usage

```hcl
variable "policy_id" {}
variable "rule_id" {}

data "huaweicloud_waf_rules_data_masking" "test" {
  policy_id = var.policy_id
  rule_id   = var.rule_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `policy_id` - (Required, String) Specifies the ID of the policy to which the data masking rules belong.

* `rule_id` - (Optional, String) Specifies the ID of the data masking rule.

* `status` - (Optional, String) Specifies the status of the data masking rule.
  The value can be **0** or **1**.
  + **0**: The rule is disabled.
  + **1**: The rule is enabled.

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project ID to which the protection policy belongs.
  If omitted, will query the data masking rules under the default enterprise project for enterprise users.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `rules` - The list of the data masking rules.

  The [rules](#rules_struct) structure is documented below.

<a name="rules_struct"></a>
The `rules` block supports:

* `id` - The ID of the data masking rule.

* `policy_id` - The ID of the policy to which the data masking rule belongs.

* `status` - The status of the data masking rule.

* `description` - The description of the data masking rule.

* `path` - The URL protected by the data masking rule.

* `field` - The position where the masked field stored

* `subfield` - The name of the masked field.

* `created_at` - The creation time of the data masking rule.
