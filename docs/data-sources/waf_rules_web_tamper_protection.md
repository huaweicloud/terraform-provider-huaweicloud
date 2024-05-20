---
subcategory: "Web Application Firewall (WAF)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_waf_rules_web_tamper_protection"
description: |-
  Use this data source to get a list of web tamper protection rules.
---

# huaweicloud_waf_rules_web_tamper_protection

Use this data source to get a list of web tamper protection rules.

## Example Usage

```hcl
data "huaweicloud_waf_rules_web_tamper_protection" "test" {
  policy_id = var.policy_id
  rule_id   = var.rule_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `policy_id` - (Required, String) Specifies the ID of the policy to which the web tamper protection rules belong.

* `rule_id` - (Optional, String) Specifies the ID of the web tamper protection rule.

* `status` - (Optional, String) Specifies the status of the web tamper protection rule.
  The value can be **0** or **1**.
  + **0**: The rule is disabled.
  + **1**: The rule is enabled.

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project ID to which the protection policy belongs.
  If omitted, will query the web tamper protection rules under the default enterprise project for enterprise users.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `rules` - The list of web tamper protection rules.

  The [rules](#rules_struct) structure is documented below.

<a name="rules_struct"></a>
The `rules` block supports:

* `id` - The ID of the web tamper protection rule.

* `policy_id` - The ID of the policy to which the web tamper protection rule belongs.

* `status` - The status of the web tamper protection rule.

* `description` - The description of the web tamper protection rule.

* `path` - The URL protected by the web tamper protection rule.

* `domain` - The domain name protected by the web tamper protection rule.

* `created_at` - The creation time of the web tamper protection rule.
