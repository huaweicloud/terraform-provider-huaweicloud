---
subcategory: "DataArts Studio"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dataarts_quality_rule_templates"
description: |-
  Use this data source to query the quality rule templates within HuaweiCloud.
---

# huaweicloud_dataarts_quality_rule_templates

Use this data source to query the quality rule templates within HuaweiCloud.

## Example Usage

### Basic Usage

```hcl
variable "workspace_id" {}
variable "name" {}

data "huaweicloud_dataarts_quality_rule_templates" "test" {
  workspace_id = var.workspace_id
  name         = var.name
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region where the quality rule templates are located.  
  If omitted, the provider-level region will be used.

* `workspace_id` - (Required, String) Specifies the ID of the workspace to which the quality rule templates belong.

* `name` - (Optional, String) Specifies the name of the quality rule template.

* `category_id` - (Optional, Int) Specifies the category ID of the quality rule template.

* `system_template` - (Optional, Bool) Specifies whether to query only system templates.

* `creator` - (Optional, String) Specifies the creator of the quality rule template.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `templates` - The list of quality rule templates that matched filter parameters.  
  The [templates](#quality_rule_templates) structure is documented below.

<a name="quality_rule_templates"></a>
The `templates` block supports:

* `id` - The ID of the quality rule template.

* `name` - The name of the quality rule template.

* `category_id` - The category ID of the quality rule template.

* `dimension` - The dimension of the quality rule template.
  + **Completeness**: Completeness dimension.
  + **Uniqueness**: Uniqueness dimension.
  + **Timeliness**: Timeliness dimension.
  + **Validity**: Validity dimension.
  + **Accuracy**: Accuracy dimension.
  + **Consistency**: Consistency dimension.

* `type` - The type of the quality rule template.
  + **Field**: Field-level rule.
  + **Table**: Table-level rule.
  + **Database**: Database-level rule.
  + **Cross-field**: Cross-field level rule.
  + **Customize**: Custom rule.

* `system_template` - Whether the quality rule template is a system template.

* `sql_info` - The definition relationship of the quality rule template.

* `abnormal_table_template` - The abnormal table template of the quality rule template.

* `result_description` - The result description of the quality rule template.

* `created_at` - The creation time of the quality rule template, in RFC3339 format.

* `creator` - The creator of the quality rule template. **System** represents the system built-in template.
