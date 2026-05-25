---
subcategory: "DataArts Studio"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dataarts_security_data_recognition_rule_groups"
description: |-
  Use this data source to get the list of DataArts Security data recognition rule groups within HuaweiCloud.
---

# huaweicloud_dataarts_security_data_recognition_rule_groups

Use this data source to get the list of DataArts Security data recognition rule groups within HuaweiCloud.

## Example Usage

### Query all data recognition rule groups under a specified workspace

```hcl
variable "workspace_id" {}

data "huaweicloud_dataarts_security_data_recognition_rule_groups" "test" {
  workspace_id = var.workspace_id
}
```

### Query the data recognition rule groups under a specified workspace and using group name to filter

```hcl
variable "workspace_id" {}
variable "group_name" {}

data "huaweicloud_dataarts_security_data_recognition_rule_groups" "test" {
  workspace_id = var.workspace_id
  name         = var.group_name
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region where the data recognition rule groups are located.  
  If omitted, the provider-level region will be used.

* `workspace_id` - (Required, String) Specifies the ID of the workspace to which the data recognition rule groups belong.

* `name` - (Optional, String) Specifies the name of the data recognition rule groups to be queried.

* `creator` - (Optional, String) Specifies the creator of the data recognition rule groups to be queried.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `groups` - The list of data recognition rule groups that matched filter parameters.  
  The [groups](#dataarts_security_data_recognition_rule_groups_attr) structure is documented below.

<a name="dataarts_security_data_recognition_rule_groups_attr"></a>
The `groups` block supports:

* `id` - The ID of the data recognition rule group.

* `name` - The name of the data recognition rule group.

* `description` - The description of the data recognition rule group.

* `created_by` - The creator of the data recognition rule group.

* `created_at` - The creation time of the data recognition rule group, in RFC3339 format.

* `updated_by` - The updater of the data recognition rule group.

* `updated_at` - The latest update time of the data recognition rule group, in RFC3339 format.

* `project_id` - The project ID to which the data recognition rule group belongs.

* `rules` - The list of data recognition rules that the group contains.  
  The [rules](#dataarts_security_data_recognition_rule_groups_rules_attr) structure is documented below.

<a name="dataarts_security_data_recognition_rule_groups_rules_attr"></a>
The `rules` block supports:

* `id` - The ID of the data recognition rule.

* `type` - The type of the data recognition rule.

* `secrecy_level` - The secrecy level name of the data recognition rule.

* `secrecy_level_num` - The secrecy level number of the data recognition rule.

* `name` - The name of the data recognition rule.

* `enable` - Whether the data recognition rule is enabled.

* `method` - The method of the data recognition rule.

* `content_expression` - The content expression of the data recognition rule.

* `column_expression` - The column expression of the data recognition rule.

* `comment_expression` - The comment expression of the data recognition rule.

* `description` - The description of the data recognition rule.

* `created_by` - The creator of the data recognition rule.

* `created_at` - The creation time of the data recognition rule, in RFC3339 format.

* `updated_by` - The updater of the data recognition rule.

* `updated_at` - The latest update time of the data recognition rule, in RFC3339 format.

* `builtin_rule_id` - The builtin rule ID of the data recognition rule.

* `category_id` - The category ID to which the data recognition rule belongs.
