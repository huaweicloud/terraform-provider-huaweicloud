---
subcategory: "Web Application Firewall (WAF)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_waf_rules_known_attack_source"
description: |-
  Use this data source to get a list of known attack source rules.
---

# huaweicloud_waf_rules_known_attack_source

Use this data source to get a list of known attack source rules.

## Example Usage

```hcl
data "huaweicloud_waf_rules_known_attack_source" "test" {
  policy_id = var.policy_id
  rule_id   = var.rule_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `policy_id` - (Required, String) Specifies the ID of the policy to which the known attack source rules belong.

* `rule_id` - (Optional, String) Specifies the ID of  the known attack source rule.

* `block_type` - (Optional, String) Specifies the block type of  the known attack source rule.
  The value can be **long_ip_block**, **long_cookie_block**, **long_params_block**, **short_ip_block**,
  **short_cookie_block** or **short_params_block**.

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project ID to which the protection policy belongs.
  If omitted, will query the known attack source rules under the default enterprise project for enterprise users.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `rules` - The list of known attack source rules.

  The [rules](#rules_struct) structure is documented below.

<a name="rules_struct"></a>
The `rules` block supports:

* `id` - The ID of the known attack source rule.

* `policy_id` - The ID of the policy to which the known attack source rule belongs.

* `block_type` - The block type of the known attack source rule.

* `block_time` - The block time of the known attack source rule.

* `description` - The description of the known attack source rule.

* `created_at` - The creation time of the known attack source rule.
