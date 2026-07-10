---
subcategory: "Data Security Center (DSC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dsc_template_rules"
description: |-
  Use this data source to get the list of DSC template rules within HuaweiCloud.
---

# huaweicloud_dsc_template_rules

Use this data source to get the list of DSC template rules within HuaweiCloud.

## Example Usage

```hcl
variable "template_id" {}

data "huaweicloud_dsc_template_rules" "test" {
  template_id = var.template_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `template_id` - (Required, String) Specifies the template ID.

* `keyword` - (Optional, String) Specifies the keyword for fuzzy search of sensitive information object name.

* `classification_ids` - (Optional, List of String) Specifies the list of classification IDs to filter results.

* `security_level_ids` - (Optional, List of String) Specifies the list of risk level IDs to filter results.

* `is_used` - (Optional, String) Specifies whether the rule is enabled.

* `rule_name` - (Optional, String) Specifies the keyword for fuzzy search of rule name.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID in UUID format.

* `template_rules_list` - The template rules information list.

  The [template_rules_list](#template_rules_list_struct) structure is documented below.

<a name="template_rules_list_struct"></a>
The `template_rules_list` block supports:

* `rule_id` - The rule ID.

* `project_id` - The project ID.

* `rule_name` - The rule name.

* `template_id` - The template ID.

* `classification_id` - The classification ID.

* `security_level_id` - The risk level ID.

* `security_level_name` - The risk level name.

* `security_level_color` - The risk level color.

* `is_used` - Whether the rule is enabled.

* `rule_description` - The rule description.
