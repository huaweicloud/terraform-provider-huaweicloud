---
subcategory: "Web Application Firewall (WAF)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_waf_rules_threat_intelligence"
description: |-
  Use this data source to get a list of threat intelligence rules.
---

# huaweicloud_waf_rules_threat_intelligence

Use this data source to get a list of threat intelligence rules.

## Example Usage

```hcl
variable "policy_id" {}

data "huaweicloud_waf_rules_threat_intelligence" "test" {
  policy_id = var.policy_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `policy_id` - (Required, String) Specifies the protection policy ID.

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project ID.
  This parameter is only valid for enterprise users.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID in UUID format.

* `items` - The list of threat intelligence rules.

  The [items](#items_struct) structure is documented below.

<a name="items_struct"></a>
The `items` block supports:

* `name` - The name of the threat intelligence rule.

* `id` - The ID of the threat intelligence rule.

* `policyid` - The protection policy ID.

* `type` - The reputation type.

* `description` - The description of the threat intelligence rule.

* `tags` - The list of IDC data center tags for the reputation type.

* `status` - The status of the threat intelligence rule. The value can be `0` (disabled) or `1` (enabled).

* `action` - The protective action taken when the rule is triggered.

  The [action](#rules_action_struct) structure is documented below.

<a name="rules_action_struct"></a>
The `action` block supports:

* `category` - The action category.  
  The valid values are as follows:
  + **pass**: Pass.
  + **block**: Block.
  + **log**: Only record.
