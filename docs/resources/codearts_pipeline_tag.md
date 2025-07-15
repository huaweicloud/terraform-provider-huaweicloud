---
subcategory: "CodeArts Pipeline"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_codearts_pipeline_tag"
description: |-
  Manages a CodeArts pipeline tag resource within HuaweiCloud.
---

# huaweicloud_codearts_pipeline_tag

Manages a CodeArts pipeline tag resource within HuaweiCloud.

## Example Usage

```hcl
variable "codearts_project_id" {}
variable "name" {}
variable "color" {}

resource "huaweicloud_codearts_pipeline_tag" "test" {
  project_id = var.codearts_project_id
  name       = var.name
  color      = var.color
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.
  Changing this creates a new resource.

* `project_id` - (Required, String, NonUpdatable) Specifies the CodeArts project ID.

* `name` - (Required, String) Specifies the tag name.

* `color` - (Required, String) Specifies the tag color.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `project_name` - Indicates the CodeArts project name.

## Import

The tag can be imported using `project_id` and `id`, e.g.

```bash
$ terraform import huaweicloud_codearts_pipeline_tag.test <project_id>/<id>
```
