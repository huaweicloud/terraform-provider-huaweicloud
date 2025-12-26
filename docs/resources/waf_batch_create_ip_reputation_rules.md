---
subcategory: "Web Application Firewall (WAF)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_waf_batch_create_ip_reputation_rules"
description: |-
  Manages a resource to batch create WAF IP reputation rules within HuaweiCloud WAF.
---

# huaweicloud_waf_batch_create_ip_reputation_rules

Manages a resource to batch create WAF IP reputation rules within HuaweiCloud WAF.

-> All WAF resources depend on WAF instances, and the WAF instances need to be purchased before they can be used.

-> This resource is a one-time action resource for batch creating IP reputation rules. Deleting this resource
   will not remove the created rules, but will only remove the resource information from the tf state file.

## Example Usage

```hcl
variable "enterprise_project_id" {}
variable "policy_ids" {
  type = list(string)
}

resource "huaweicloud_waf_batch_create_ip_reputation_rules" "test" {
  name                  = "test_ip_reputation_rule"
  type                  = "idc"
  tags                  = ["Tencent"]
  policy_ids            = var.policy_ids
  enterprise_project_id = var.enterprise_project_id
  description           = "test_ip_reputation_rule"

  action {
    category = "block"
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.
  Changing this parameter will create a new resource.

* `name` - (Required, String, NonUpdatable) Specifies the name of the IP reputation rule.

* `type` - (Required, String, NonUpdatable) Specifies the type of the IP reputation rule.
  Only **idc** is supported.

* `tags` - (Required, List, NonUpdatable) Specifies the list of tags for the IP reputation rule.

* `action` - (Required, List, NonUpdatable) Specifies the action to take when the IP reputation rule is triggered.
  The [action](#Rule_action) structure is documented below.

* `policy_ids` - (Required, List, NonUpdatable) Specifies the list of policy IDs to which the IP reputation rule will be
  applied.

* `description` - (Optional, String, NonUpdatable) Specifies the description of the IP reputation rule.

* `enterprise_project_id` - (Optional, String, NonUpdatable) Specifies the enterprise project ID.
  Value `0` indicates the default enterprise project.
  Value **all_granted_eps** indicates all enterprise projects to which the user has been granted access.
  Defaults to `0`.

<a name="Rule_action"></a>
The `action` block supports:

* `category` - (Required, String, NonUpdatable) Specifies the action to take when the IP reputation rule is triggered.
  Valid values are **pass**, **log**, and **block**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.
