---
subcategory: "Web Application Firewall (WAF)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_waf_batch_create_antileakage_rules"
description: |-
  Manages a resource to batch create anti-leakage rules within HuaweiCloud WAF.
---

# huaweicloud_waf_batch_create_antileakage_rules

Manages a resource to batch create anti-leakage rules within HuaweiCloud WAF.

-> All WAF resources depend on WAF instances, and the WAF instances need to be purchased before they can be used.

-> This resource is a one-time action resource for batch creating anti-leakage rules. Deleting this resource
   will not remove the created rules, but will only remove the resource information from the tf state file.

## Example Usage

```hcl
variable "enterprise_project_id" {}
variable "policy_ids" {
  type = list(string)
}

resource "huaweicloud_waf_batch_create_antileakage_rules" "test" {
  url                   = "/admin"
  category              = "sensitive"
  contents              = ["id_card", "phone"]
  policy_ids            = var.policy_ids
  description           = "test description"
  enterprise_project_id = var.enterprise_project_id

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

* `url` - (Required, String, NonUpdatable) Specifies the URL to which the rule applies.

* `category` - (Required, String, NonUpdatable) Specifies the type of the rule.
  The value can be:
  + **code**: Response code
  + **sensitive**: Sensitive information

* `contents` - (Required, List, NonUpdatable) Specifies the content of the rule.
  + When `category` is **code**, the value can be: `400`, `401`, `402`, `403`, `404`, `405`, `500`, `501`, `502`, `503`,
    `504`, `507`.
  + When `category` is **sensitive**, the value can be: **phone** (mobile number), **id_card** (ID card number),
    **email** (email address).

* `policy_ids` - (Required, List, NonUpdatable) Specifies the list of policy IDs to which the rule will be applied.

* `action` - (Optional, List, NonUpdatable) Specifies the action to take when the rule is matched.
  The [action](#Rule_action) structure is documented below.

* `description` - (Optional, String, NonUpdatable) Specifies the description of the rule.

* `enterprise_project_id` - (Optional, String, NonUpdatable) Specifies the enterprise project ID.
  Value `0` indicates the default enterprise project.
  Value **all_granted_eps** indicates all enterprise projects to which the user has been granted access.
  Defaults to `0`.

<a name="Rule_action"></a>
The `action` block supports:

* `category` - (Required, String, NonUpdatable) Specifies the action to take when the rule is matched.
  The value can be:
  + **block**: Block the request
  + **log**: Log the request only

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.
