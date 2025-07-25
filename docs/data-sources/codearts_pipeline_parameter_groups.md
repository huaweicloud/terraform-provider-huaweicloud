---
subcategory: "CodeArts Pipeline"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_codearts_pipeline_parameter_groups"
description: |-
  Use this data source to get a list of CodeArts pipeline parameter groups.
---

# huaweicloud_codearts_pipeline_parameter_groups

Use this data source to get a list of CodeArts pipeline parameter groups.

## Example Usage

```hcl
variable "codearts_project_id" {}

data "huaweicloud_codearts_pipeline_parameter_groups" "test" {
  project_id = var.codearts_project_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `project_id` - (Required, String) Specifies the CodeArts project ID.

* `name` - (Optional, String) Specifies the parameter group name.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `groups` - Indicates the parameter group list.
  The [groups](#attrblock--groups) structure is documented below.

<a name="attrblock--groups"></a>
The `groups` block supports:

* `id` - Indicates the parameter group ID.

* `name` - Indicates the parameter group name.

* `description` - Indicates the parameter group description.

* `related_pipelines` - Indicates the associated pipeline.
  The [related_pipelines](#attrblock--groups--related_pipelines) structure is documented below.

* `variables` - Indicates the parameter list.
  The [variables](#attrblock--groups--variables) structure is documented below.

* `create_time` - Indicates the create time.

* `creator_id` - Indicates the creator ID.

* `creator_name` - Indicates the creator name.

* `update_time` - Indicates the update time.

* `updater_id` - Indicates the updater ID.

* `updater_name` - Indicates the updater name.

<a name="attrblock--groups--related_pipelines"></a>
The `related_pipelines` block supports:

* `id` - Indicates the pipeline ID.

* `name` - Indicates the pipeline name.

<a name="attrblock--groups--variables"></a>
The `variables` block supports:

* `description` - Indicates the parameter description.

* `is_secret` - Indicates whether it is a private parameter.

* `name` - Indicates the custom variable name.

* `sequence` - Indicates the parameter sequence, starting from 1.

* `type` - Indicates the custom parameter type.

* `value` - Indicates the custom parameter default value.
