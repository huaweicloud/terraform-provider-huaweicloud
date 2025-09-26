---
subcategory: "Web Application Firewall (WAF)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_waf_all_ip_reputation_policy_rules"
description: |-
  Use this data source to get list of the WAF IP reputation policy rules under all policies.
---

# huaweicloud_waf_all_ip_reputation_policy_rules

Use this data source to get list of the WAF IP reputation policy rules under all policies.

## Example Usage

```hcl
data "huaweicloud_waf_all_ip_reputation_policy_rules" "test" {}
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

* `items` - The list of the WAF IP reputation policy rules.

  The [items](#items_struct) structure is documented below.

<a name="items_struct"></a>
The `items` block supports:

* `id` - The ID of the rule.

* `policyid` - The ID of the policy.

* `name` - The name of the rule.

* `type` - The reputation type of the rule.

* `description` - The description of the rule.

* `tags` - The list of IDC data centers for the reputation type.

* `status` - The status of the rule.
  The valid values are as follows:
  + `0`: Disabled
  + `1`: Enabled

* `action` - The protection action of the rule.

  The [action](#action_struct) structure is documented below.

<a name="action_struct"></a>
The `action` block supports:

* `category` - The action type of the rule.
  The valid values are as follows:
  + **pass**: Allow
  + **block**: Block
  + **log**: Log only
