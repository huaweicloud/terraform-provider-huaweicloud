---
subcategory: "Config"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_rms_policy_definitions"
description: ""
---

# huaweicloud_rms_policy_definitions

Use this data source to query policy definition list.

## Example Usage

```hcl
variable "trigger_type" {}

data "huaweicloud_rms_policy_definitions" "test" {
  trigger_type = var.trigger_type
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Optional, String) Specifies the name of the policy definitions used to query definition list.

* `policy_type` - (Optional, String) Specifies the policy type used to query definition list.  
  The valid value is **builtin**.

* `policy_rule_type` - (Optional, String) Specifies the policy rule type used to query definition list.

* `trigger_type` - (Optional, String) Specifies the trigger type used to query definition list.  
  The valid values are **resource** and **period**.

* `keywords` - (Optional, List) Specifies the keyword list used to query definition list.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `definitions` - The policy definition list.
  The [object](#policy_definitions) structure is documented below.

<a name="policy_definitions"></a>
The `definitions` block supports:

* `id` - The ID of the policy definition.

* `name` - The name of the policy definition.

* `policy_type` - The policy type of the policy definition.

* `description` - The description of the policy definition.

* `policy_rule_type` - The policy rule type of the policy definition.

* `policy_rule` - The policy rule of the policy definition.

* `trigger_type` - The trigger type of the policy definition.

* `keywords` - The keyword list of the policy definition.

* `parameters` - The parameter reference map of the policy definition.
