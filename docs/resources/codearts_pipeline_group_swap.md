---
subcategory: "CodeArts Pipeline"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_codearts_pipeline_group_swap"
description: |-
  Manages a CodeArts pipeline group swap resource within HuaweiCloud.
---
# huaweicloud_codearts_pipeline_group_swap

Manages a CodeArts pipeline group swap resource within HuaweiCloud.

## Example Usage

```hcl
variable "codearts_project_id" {}
variable "group_id1" {}
variable "group_id2" {}

resource "huaweicloud_codearts_pipeline_group_swap" "test" {
  project_id = var.codearts_project_id
  group_id1  = var.group_id1
  group_id2  = var.group_id2
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.
  Changing this creates a new resource.

* `project_id` - (Required, String, NonUpdatable) Specifies the project ID for CodeArts service.

* `group_id1` - (Required, String, NonUpdatable) Specifies the pipeline group ID1.

* `group_id2` - (Required, String, NonUpdatable) Specifies the pipeline group ID2.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.
