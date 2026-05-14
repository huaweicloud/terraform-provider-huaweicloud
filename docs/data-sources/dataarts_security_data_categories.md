---
subcategory: "DataArts Studio"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dataarts_security_data_categories"
description: |-
  Use this data source to query DataArts Security data category list within HuaweiCloud.
---

# huaweicloud_dataarts_security_data_categories

Use this data source to query DataArts Security data category list within HuaweiCloud.

## Example Usage

```hcl
variable "workspace_id" {}

data "huaweicloud_dataarts_security_data_categories" "test" {
  workspace_id = var.workspace_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region where the data categories are located.  
  If omitted, the provider-level region will be used.

* `workspace_id` - (Required, String) Specifies the ID of the workspace to which the data categories belong.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `categories` - The list of data categories.  
  The [categories](#dataarts_security_data_categories_attr) structure is documented below.

<a name="dataarts_security_data_categories_attr"></a>
The `categories` block supports:

* `category_id` - The ID of the data category.

* `category_name` - The name of the data category.

* `description` - The description of the data category.

* `category_level` - The level of the data category in the tree.

* `root_id` - The ID of the root node of the category tree.

* `parent_id` - The ID of the parent category.

* `category_path` - The path of the category in the tree.

* `instance_id` - The instance ID to which the data category belongs.

* `synchronize` - Whether the data category is synchronized with assets.

* `children` - The children information of the data category, in JSON format.

* `rules` - The list of data classification rules associated with the category.  
  The [rules](#dataarts_security_data_categories_rules) structure is documented below.

* `create_by` - The creator of the data category.

* `create_time` - The creation time of the data category, in RFC3339 format.

* `update_by` - The updater of the data category.

* `update_time` - The latest update time of the data category, in RFC3339 format.

<a name="dataarts_security_data_categories_rules"></a>
The `rules` block supports:

* `id` - The ID of the rule.

* `name` - The name of the rule.

* `type` - The type of the rule.

* `secrecy_level` - The secrecy level name.

* `secrecy_level_num` - The secrecy level number.

* `enable` - Whether the rule is enabled.

* `method` - The method of the rule.

* `content_expression` - The content expression of the rule.

* `column_expression` - The column expression of the rule.

* `comment_expression` - The comment expression of the rule.

* `combine_expression` - The combine expression of the rule.

* `description` - The description of the rule.

* `builtin_rule_id` - The builtin rule ID.

* `match_type` - The match type of the rule.

* `created_by` - The creator of the rule.

* `created_at` - The creation time of the rule, in RFC3339 format.

* `updated_by` - The updater of the rule.

* `updated_at` - The latest update time of the rule, in RFC3339 format.
