---
subcategory: "DataArts Studio"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dataarts_security_data_recognition_rules"
description: |-
  Use this data source to get the list of DataArts Security data recognition rules within HuaweiCloud.
---

# huaweicloud_dataarts_security_data_recognition_rules

Use this data source to get the list of DataArts Security data recognition rules within HuaweiCloud.

## Example Usage

### Query all data recognition rules under a specified workspace

```hcl
variable "workspace_id" {}

data "huaweicloud_dataarts_security_data_recognition_rules" "test" {
  workspace_id = var.workspace_id
}
```

### Query the data recognition rules under a specified workspace and using rule name to filter

```hcl
variable "workspace_id" {}
variable "rule_name" {}

data "huaweicloud_dataarts_security_data_recognition_rules" "test" {
  workspace_id = var.workspace_id
  rule_name    = var.rule_name
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region where the data recognition rules are located.  
  If omitted, the provider-level region will be used.

* `workspace_id` - (Required, String) Specifies the ID of the workspace to which the data recognition rules belong.

* `secrecy_level` - (Optional, String) Specifies the secrecy level to which the data recognition rules belong.

* `rule_name` - (Optional, String) Specifies the name of the specified data recognition rule to be queried.

* `creator` - (Optional, String) Specifies the creator of the data recognition rules to be queried.

* `enable` - (Optional, Bool) Specifies whether the data recognition rules are enabled to be queried.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `rules` - The list of data recognition rules that matched filter parameters.  
  The [rules](#dataarts_security_data_recognition_rules_attr) structure is documented below.

<a name="dataarts_security_data_recognition_rules_attr"></a>
The `rules` block supports:

* `id` - The ID of the data recognition rule.

* `type` - The type of the data recognition rule.

* `secrecy_level` - The secrecy level name of the data recognition rule.

* `secrecy_level_num` - The secrecy level number of the data recognition rule.

* `name` - The name of the data recognition rule.

* `enable` - Whether the rule is enabled.

* `method` - The method of the data recognition rule.

* `content_expression` - The content expression of the data recognition rule.

* `column_expression` - The column expression of the data recognition rule.

* `comment_expression` - The comment expression of the data recognition rule.

* `combine_expression` - The combine expression of the data recognition rule.

* `description` - The description of the data recognition rule.

* `created_by` - The creator of the data recognition rule.

* `created_at` - The creation time of the data recognition rule, in RFC3339 format.

* `updated_by` - The updater of the data recognition rule.

* `updated_at` - The latest update time of the data recognition rule, in RFC3339 format.

* `builtin_rule_id` - The builtin rule ID of the data recognition rule.

* `category_id` - The category ID to which the data recognition rule belongs.

* `instance_id` - The instance ID to which the data recognition rule belongs.

* `match_type` - The match type of the data recognition rule.
