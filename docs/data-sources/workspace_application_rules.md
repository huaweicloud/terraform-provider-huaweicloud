---
subcategory: "Workspace"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_workspace_application_rules"
description: |-
  Use this data source to get the list of the Workspace application rules within HuaweiCloud.
---

# huaweicloud_workspace_application_rules

Use this data source to get the list of the Workspace application rules within HuaweiCloud.

## Example Usage

### Basic Usage

```hcl
data "huaweicloud_workspace_application_rules" "test" {}
```

### Filter application rules by name

```hcl
variable "rule_name" {}

data "huaweicloud_workspace_application_rules" "test" {
  name = var.rule_name
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region where the application rules are located.

* `name` - (Optional, String) Specifies the name of the application rule to be queried.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `rules` - The list of application rules that match the filter parameters.  
  The [rules](#workspace_application_rules_attr) structure is documented below.

<a name="workspace_application_rules_attr"></a>
The `rules` block supports:

* `id` - The ID of the application rule.

* `name` - The name of the application rule.

* `description` - The description of the application rule.

* `rule_source` - The source of the application rule.

* `create_time` - The create time of the application rule, in RFC3339 format.

* `update_time` - The update time of the application rule, in RFC3339 format.

* `detail` - The detail of the application rule.  
  The [detail](#workspace_application_rule_detail) structure is documented below.

<a name="workspace_application_rule_detail"></a>
The `detail` block supports:

* `scope` - The scope of the rule.  
  The valid values are as follows:
  + **PRODUCT**
  + **PATH**

* `product_rule` - The product rule configuration.  
  The [product_rule](#workspace_application_rule_product_rule_attr) structure is documented below.

* `path_rule` - The path rule configuration.  
  The [path_rule](#workspace_application_rule_path_rule_attr) structure is documented below.

<a name="workspace_application_rule_product_rule_attr"></a>
The `product_rule` block supports:

* `identify_condition` - The identification condition.

* `publisher` - The publisher name.

* `product_name` - The product name.

* `process_name` - The process name.

* `support_os` - The supported operating system type, defaults to **Windows**

* `version` - The version number.

* `product_version` - The product version number.

<a name="workspace_application_rule_path_rule_attr"></a>
The `path_rule` block supports:

* `path` - The complete path.
