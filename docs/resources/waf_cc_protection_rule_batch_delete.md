---
subcategory: "Web Application Firewall (WAF)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_waf_cc_protection_rule_batch_delete"
description: |-
  Manages a resource to batch delete the CC protection rules within HuaweiCloud.
---

# huaweicloud_waf_cc_protection_rule_batch_delete

Manages a resource to batch delete the CC protection rules within HuaweiCloud.

-> All WAF resources depend on WAF instances, and the WAF instances need to be purchased before they can be used.

-> This resource is only a one-time action resource using to batch delete CC protection rules. Deleting this resource
  will not clear the corresponding request record, but will only remove the resource information from the tf state
  file.

## Example Usage

```hcl
variable "policy_id" {}
variable "rule_ids"  {
  type = list(string)
}

resource "huaweicloud_waf_cc_protection_rule_batch_delete" "test" {
  policy_rule_ids {
    policy_id = var.policy_id
    rule_ids  = var.rule_ids
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.
  Changing this parameter will create a new resource.

* `policy_rule_ids` - (Required, List, NonUpdatable) Specifies the policy and rule details.
  The [policy_rule_ids](#cc_protection_policy_rules) structure is documented below.

<a name="cc_protection_policy_rules"></a>
The `policy_rule_ids` block supports:

* `policy_id` - (Required, String, NonUpdatable) Specifies the policy ID to which the CC protection rule belongs.

* `rule_ids` - (Required, List, NonUpdatable) Specifies the ID list of the CC protection rule.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.
