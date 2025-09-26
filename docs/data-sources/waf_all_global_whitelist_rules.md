---
subcategory: "Web Application Firewall (WAF)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_waf_all_global_whitelist_rules"
description: |-
  Use this data source to get list of the WAF global whitelist rules under all policies.
---

# huaweicloud_waf_all_global_whitelist_rules

Use this data source to get list of the WAF global whitelist rules under all policies.

## Example Usage

```hcl
data "huaweicloud_waf_all_global_whitelist_rules" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `policyids` - (Optional, String) Specifies the ID of the policy.

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project ID.
  If you want to query resources under all enterprise projects, set this parameter to **all_granted_eps**.
  Defaults to **0**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `items` - The list of the WAF global whitelist rules.

  The [items](#items_struct) structure is documented below.

<a name="items_struct"></a>
The `items` block supports:

* `id` - The ID of the rule.

* `policyid` - The ID of the policy.

* `timestamp` - The creation time of the rule, in milliseconds.

* `description` - The description of the rule.

* `status` - The status of the rule.
  The valid values are as follows:
  + `0`: Disabled.
  + `1`: Enabled.

* `url` - The false alarm rule blocks the path.
  Only valid when `mode` value is `0`.

* `rule` - The rule that need to be blocked.

* `mode` - The version.
  The value `0` indicates old version v1, value `1` indicates new version v2.

* `url_logic` - The matching logic.

* `conditions` - The condition list.

  The [conditions](#items_condition_struct) structure is documented below.

* `domain` - The protection domain name or website.

* `advanced` - The condition list.

  The [advanced](#items_advanced_struct) structure is documented below.

<a name="items_condition_struct"></a>
The `conditions` block supports:

* `category` - The field type.
  The value can be **ip**, **url**, **params**, **cookie** or **header**.

* `contents` - The content.

* `logic_operation` - The matching logic.

* `index` - The subfield.

<a name="items_advanced_struct"></a>
The `advanced` block supports:

* `index` - The field type.
  The value can be **params**, **cookie**, **header**, **body** or **multipart**.

* `contents` - The subfield of the field type.
