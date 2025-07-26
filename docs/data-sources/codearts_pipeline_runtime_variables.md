---
subcategory: "CodeArts Pipeline"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_codearts_pipeline_runtime_variables"
description: |-
  Use this data source to get a list of CodeArts pipeline custom parameters required for pipeline running.
---

# huaweicloud_codearts_pipeline_runtime_variables

Use this data source to get a list of CodeArts pipeline custom parameters required for pipeline running.

## Example Usage

```hcl
variable "codearts_project_id" {}
variable "pipeline_id" {}

data "huaweicloud_codearts_pipeline_runtime_variables" "test" {
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

* `name` - Indicates the custom variable name.

* `sequence` - Indicates the parameter sequence, starting from 1.

* `type` - Indicates the custom parameter type.

* `value` - Indicates the custom parameter default value.

* `description` - Indicates the parameter description.

* `is_reset` - Indicates whether to reset.

* `is_runtime` - Indicates whether to set parameters at runtime.

* `is_secret` - Indicates whether it is a private parameter.

* `latest_value` - Indicates the last parameter value.

* `limits` - Indicates the list of enumerated values.
