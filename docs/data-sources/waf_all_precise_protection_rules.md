---
subcategory: "Web Application Firewall (WAF)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_waf_all_precise_protection_rules"
description: |-
  Use this data source to get list of the WAF precise protection rules under all policies.
---

# huaweicloud_waf_all_precise_protection_rules

Use this data source to get list of the WAF precise protection rules under all policies.

## Example Usage

```hcl
data "huaweicloud_waf_all_precise_protection_rules" "test" {}
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

* `items` - The list of the WAF blacklist and whitelist protective rules.

  The [items](#items_struct) structure is documented below.

<a name="items_struct"></a>
The `items` block supports:

* `id` - The ID of the rule.

* `policyid` - The ID of the policy.

* `name` - The name of the rule.

* `description` - The description of the rule.

* `status` - The status of the rule.
  The valid values are as follows:
  + `0`: Disabled.
  + `1`: Enabled.

* `conditions` - The matching condition list of the rule.

  The [conditions](#items_conditions_struct) structure is documented below.

* `action` - The operation target after precise protection rule is triggered.

  The [action](#items_action_struct) structure is documented below.

* `priority` - The priority of executing the rule.

* `timestamp` - The creation time of the rule, in milliseconds.

* `time` - The effective time of the rule.
  The valid values are as follows:
  + **true**: Indicates the rule takes effect immediately.
  + **false**: Indicates a custom effective time.

* `start` - The start time when the rule takes effect.
  This parameter is valid only when the `time` is set to **true**.

* `terminal` - The end time when the rule takes effect.
  This parameter is valid only when the `time` is set to **true**.

<a name="items_conditions_struct"></a>
The `conditions` block supports:

* `category` - The field type.
  The value can be **url**, **user-agent**, **ip**, **params**, **cookie**, **referer**, **header**,
  **request_line**, **method** or **request**.

* `index` - The subfield.

* `logic_operation` - The condition matching logic.

* `contents` - The condition matching contents.

* `value_list_id` - The reference table ID.

<a name="items_action_struct"></a>
The `action` block supports:

* `category` - The operation type.
  The value can be **block**, **pass** or **log**.

* `followed_action_id` - The known attack source rule ID.
