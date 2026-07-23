---
subcategory: "Data Security Center (DSC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dsc_test_scan_rule"
description: |-
  Manages a resource to test a DSC scan rule within HuaweiCloud.
---

# huaweicloud_dsc_test_scan_rule

Manages a resource to test a DSC scan rule within HuaweiCloud.

-> 1. This resource is a one-time action resource used to test a DSC scan rule. Deleting this resource will not undo
  the test action, but will only remove the resource information from the tf state file.
  <br/>2. The execution result of this operation is based on the values of the `is_match` and `match_group` fields.

## Example Usage

```hcl
variable "category" {}
variable "data" {}
variable "effective_mode" {}
variable "location" {}
variable "rule_content" {
  type = list(string)
}
variable "rule_id" {}

resource "huaweicloud_dsc_test_scan_rule" "test" {
  category       = var.category
  data           = var.data
  effective_mode = var.effective_mode
  location       = var.location
  rule_content   = var.rule_content
  rule_id        = var.rule_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `category` - (Optional, String, NonUpdatable) Specifies the rule category.
  The valid values are **BUILT_IN** and **BUILT_SELF**.

* `data` - (Optional, String, NonUpdatable) Specifies the data to be tested.

* `effective_mode` - (Optional, String, NonUpdatable) Specifies the effective mode of the rule.

* `location` - (Optional, String, NonUpdatable) Specifies the location where the rule is applied.

* `rule_content` - (Optional, List, NonUpdatable) Specifies the rule content list to be tested.

* `rule_id` - (Optional, String, NonUpdatable) Specifies the rule ID.

* `rule_name` - (Optional, String, NonUpdatable) Specifies the rule name.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `is_match` - Whether the test data matches the rule.

* `match_group` - The matched rule group.
