---
subcategory: "Web Application Firewall (WAF)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_waf_batch_create_custom_rules"
description: |-
  Manages a resource to batch create WAF custom (precise protection) rules within HuaweiCloud WAF.
---

# huaweicloud_waf_batch_create_custom_rules

Manages a resource to batch create WAF custom rules (precise protection rules) within HuaweiCloud WAF.

-> All WAF resources depend on WAF instances, and the WAF instances need to be purchased before they can be used.

-> This resource is a one-time action resource for batch creating custom rules. Deleting this resource
   will not remove the created rules, but will only remove the resource information from the tf state file.

## Example Usage

```hcl
variable "enterprise_project_id" {}
variable "policy_ids" {
  type = list(string)
}

resource "huaweicloud_waf_batch_create_custom_rules" "test" {
  name        = "test_rule"
  description = "test description"
  time        = false
  priority    = 10
  policy_ids  = var.policy_ids
  
  conditions {
    category        = "url"
    logic_operation = "equal"
    contents        = ["/admin"]
  }

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

* `time` - (Required, Bool, NonUpdatable) Specifies whether to customize the effective time of the rule.
  The value can be:
  + **false**: The rule takes effect immediately
  + **true**: Customize the effective time

* `conditions` - (Required, List, NonUpdatable) Specifies the list of matching conditions.
  The [conditions](#Rule_conditions) structure is documented below.

* `action` - (Required, List, NonUpdatable) Specifies the action to take when the rule is matched.
  The [action](#Rule_action) structure is documented below.

* `priority` - (Required, Int, NonUpdatable) Specifies the priority of the rule.
  The value ranges from `0` to `65,535`. A smaller value indicates a higher priority.

* `name` - (Required, String, NonUpdatable) Specifies the name of the rule.

* `policy_ids` - (Required, List, NonUpdatable) Specifies the list of policy IDs to which the rule will be applied.

* `start` - (Optional, Int, NonUpdatable) Specifies the start timestamp (in seconds) when the rule takes effect.
  This parameter is required when `time` is **true**.

* `terminal` - (Optional, Int, NonUpdatable) Specifies the end timestamp (in seconds) when the rule takes effect.
  This parameter is required when `time` is **true**.

* `description` - (Optional, String, NonUpdatable) Specifies the description of the rule.

* `enterprise_project_id` - (Optional, String, NonUpdatable) Specifies the enterprise project ID.
  Value `0` indicates the default enterprise project.
  Value **all_granted_eps** indicates all enterprise projects to which the user has been granted access.
  Defaults to `0`.

<a name="Rule_conditions"></a>
The `conditions` block supports:

* `category` - (Optional, String) Specifies the field type to match against.
  The value can be:
  + **url**
  + **user-agent**
  + **ip**
  + **params**
  + **cookie**
  + **referer**
  + **header**
  + **request_line**
  + **method**
  + **request**

* `index` - (Optional, String) Specifies the subfield name.
  + When `category` is **params**, **header**, or **cookie**, the value of `index` is custom subfield.
  + For other `category` values, leave this empty.

* `logic_operation` - (Optional, String) Specifies the matching logic operation.
  + When `category` is **url**, **user-agent**, or **referer**, the value can be: **equal**, **not_equal**, **contain**,
    **not_contain**, **prefix**, **not_prefix**, **suffix**, **not_suffix**, **contain_any**, **not_contain_all**,
    **equal_any**, **not_equal_all**, **prefix_any**, **not_prefix_all**, **suffix_any**, **not_suffix_all**,
    **len_greater**, **len_less**, **len_equal**, **len_not_equal**
  + When `category` is **ip**, the value can be: **equal**, **not_equal**, **equal_any**, **not_equal_all**
  + When `category` is **method**, the value can be: **equal**, **not_equal**
  + When `category` is **request_line** or **request**, the value can be: **len_greater**, **len_less**, **len_equal**,
    **len_not_equal**
  + When `category` is **params**, **header**, or **cookie**, the value can be: **contain**, **not_contain**, **equal**,
    **not_equal**, **prefix**, **not_prefix**, **suffix**, **not_suffix**, **contain_any**, **not_contain_all**,
    **equal_any**, **not_equal_all**, **prefix_any**, **not_prefix_all**, **suffix_any**, **not_suffix_all**,
    **len_greater**, **len_less**, **len_equal**, **len_not_equal**, **num_greater**, **num_less**, **num_equal**,
    **num_not_equal**, **exist**, **not_exist**

* `contents` - (Optional, List) Specifies the content to match against.
  This parameter is required when `logic_operation` does not end with **any** or **all**.

* `value_list_id` - (Optional, String) Specifies the reference table ID.
  This parameter is required when `logic_operation` ends with **any** or **all**.

<a name="Rule_action"></a>
The `action` block supports:

* `category` - (Required, String) Specifies the action to take when the rule is matched.
  The value can be:
  + **block**: Block the request
  + **pass**: Allow the request
  + **log**: Log the request only

* `followed_action_id` - (Optional, String) Specifies the ID of the attack penalty rule.
  This parameter is valid only when `category` is set to **block**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.
