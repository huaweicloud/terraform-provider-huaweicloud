---
subcategory: "CodeArts Pipeline"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_codearts_pipeline_modify_histories"
description: |-
  Use this data source to get a list of CodeArts pipeline modify histories.
---

# huaweicloud_codearts_pipeline_modify_histories

Use this data source to get a list of CodeArts pipeline modify histories.

## Example Usage

```hcl
variable "codearts_project_id" {}
variable "pipeline_id" {}

data "huaweicloud_codearts_pipeline_modify_histories" "test" {
  project_id  = var.codearts_project_id
  pipeline_id = var.pipeline_id
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

* `histories` - Indicates the history list.
  The [histories](#attrblock--histories) structure is documented below.

<a name="attrblock--histories"></a>
The `histories` block supports:

* `modify_type` - Indicates the modify type.

* `creator_name` - Indicates the creator name.

* `creator_nick_name` - Indicates the creator nick name.

* `create_time` - Indicates the create time.
