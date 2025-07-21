---
subcategory: "CodeArts Pipeline"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_codearts_pipeline_parameter_group"
description: |-
  Manages a CodeArts pipeline parameter group resource within HuaweiCloud.
---

# huaweicloud_codearts_pipeline_parameter_group

Manages a CodeArts pipeline parameter group resource within HuaweiCloud.

## Example Usage

```hcl
variable "codearts_project_id" {}
variable "name" {}

resource "huaweicloud_codearts_pipeline_parameter_group" "test" {
  project_id = var.codearts_project_id
  name       = var.name
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.
  Changing this creates a new resource.

* `project_id` - (Required, String, NonUpdatable) Specifies the CodeArts project ID.

* `name` - (Required, String) Specifies the parameter group name.

* `description` - (Optional, String) Specifies the parameter group description.

* `variables` - (Optional, List) Specifies the parameter list.
  The [variables](#block--variables) structure is documented below.

<a name="block--variables"></a>
The `variables` block supports:

* `description` - (Optional, String) Specifies the parameter description.

* `is_secret` - (Optional, Bool) Specifies whether it is a private parameter.

* `name` - (Optional, String) Specifies the custom variable name.

* `sequence` - (Optional, Int) Specifies the parameter sequence, starting from 1.

* `type` - (Optional, String) Specifies the custom parameter type.

* `value` - (Optional, String) Specifies the custom parameter default value.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `related_pipelines` - Indicates the associated pipeline.
  The [related_pipelines](#attrblock--related_pipelines) structure is documented below.

* `creator_id` - Indicates the creator ID.

* `updater_id` - Indicates the updater ID.

* `create_time` - Indicates the create time.

* `update_time` - Indicates the update time.

* `creator_name` - Indicates the creator name.

* `updater_name` - Indicates the updater name.

<a name="attrblock--related_pipelines"></a>
The `related_pipelines` block supports:

* `id` - Indicates the pipeline ID.

* `name` - Indicates the pipeline name.

## Import

The parameter group can be imported using `project_id` and `id`, e.g.

```bash
$ terraform import huaweicloud_codearts_pipeline_parameter_group.test <project_id>/<id>
```
