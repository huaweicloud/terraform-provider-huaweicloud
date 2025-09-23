---
subcategory: "CodeArts Pipeline"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_codearts_pipeline_system_variables"
description: |-
  Use this data source to get the CodeArts pipeline available predefined parameters.
---

# huaweicloud_codearts_pipeline_system_variables

Use this data source to get the CodeArts pipeline available predefined parameters.

## Example Usage

```hcl
variable "codearts_project_id" {}
variable "pipeline_id" {}

data "huaweicloud_codearts_pipeline_system_variables" "test" {
  project_id  = var.codearts_project_id
  pipeline_id = var.pipeline_id"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `project_id` - (Required, String) Specifies the CodeArts project ID.

* `pipeline_id` - (Required, String) Specifies the pipeline ID.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `variables` - Indicates the pipeline variables list.
  The [variables](#attrblock--variables) structure is documented below.

<a name="attrblock--variables"></a>
The `variables` block supports:

* `ordinal` - Indicates the parameter ordinal.

* `name` - Indicates the system variable name.

* `type` - Indicates the system parameter type.

* `value` - Indicates the system parameter value.

* `description` - Indicates the parameter description.

* `is_show` - Indicates whether it is showed.

* `is_alias` - Indicates whether the name is alias.

* `kind` - Indicates the parameter context type.

* `context_name` - Indicates the context name.

* `source_identifier` - Indicates the source identifier.
