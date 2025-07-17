---
subcategory: "CodeArts Pipeline"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_codearts_pipeline_tags"
description: |-
  Use this data source to get a list of CodeArts pipeline tags.
---

# huaweicloud_codearts_pipeline_tags

Use this data source to get a list of CodeArts pipeline tags.

## Example Usage

```hcl
variable "codearts_project_id" {}

data "huaweicloud_codearts_pipeline_tags" "test" {
  project_id = var.codearts_project_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `project_id` - (Required, String) Specifies the CodeArts project ID.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `tags` - Indicates the tag list.
  The [tags](#attrblock--tags) structure is documented below.

<a name="attrblock--tags"></a>
The `tags` block supports:

* `id` - Indicates the tag ID.

* `name` - Indicates the tag name.

* `color` - Indicates the tag color.
