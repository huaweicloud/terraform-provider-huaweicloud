---
subcategory: "Workspace"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_workspace_app_rules"
description: |-
  Use this data source to get the list of the Workspace app rules within HuaweiCloud.
---

# huaweicloud_workspace_app_rules

Use this data source to get the list of the Workspace app rules within HuaweiCloud.

## Example Usage

### Basic Usage

```hcl
data "huaweicloud_workspace_app_rules" "test" {}
```

### Filter app rules by name

```hcl
variable "app_rule_name" {}

data "huaweicloud_workspace_app_rules" "test" {
  name = var.app_rule_name
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region where the app rules are located.

* `name` - (Optional, String) Specifies the name of the app rule to be queried.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `app_rules` - The list of app rules that match the filter parameters.  
  The [app_rules](#workspace_app_rules_attr) structure is documented below.

<a name="workspace_app_rules_attr"></a>
The `app_rules` block supports:

* `id` - The ID of the app rule.

* `name` - The name of the app rule.

* `description` - The description of the app rule.

* `rule_source` - The source of the app rule.

* `create_time` - The create time of the app rule, in RFC3339 format.

* `update_time` - The update time of the app rule, in RFC3339 format.

* `rule` - The rule configuration.  
  The [rule](#workspace_app_rule_attr) structure is documented below.

<a name="workspace_app_rule_attr"></a>
The `rule` block supports:

* `scope` - The scope of the rule.  
  The valid values are as follows:
  + **PRODUCT**
  + **PATH**

* `product_rule` - The product rule configuration.  
  The [product_rule](#workspace_app_rule_product_rule_attr) structure is documented below.

* `path_rule` - The path rule configuration.  
  The [path_rule](#workspace_app_rule_path_rule_attr) structure is documented below.

<a name="workspace_app_rule_product_rule_attr"></a>
The `product_rule` block supports:

* `identify_condition` - The identification condition.

* `publisher` - The publisher name.

* `product_name` - The product name.

* `process_name` - The process name.

* `support_os` - The supported operating system type.

* `version` - The version number.

* `product_version` - The product version number.

<a name="workspace_app_rule_path_rule_attr"></a>
The `path_rule` block supports:

* `path` - The complete path.
