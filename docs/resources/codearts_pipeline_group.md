---
subcategory: "CodeArts Pipeline"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_codearts_pipeline_group"
description: |-
  Manages a CodeArts pipeline group resource within HuaweiCloud.
---

# huaweicloud_codearts_pipeline_group

Manages a CodeArts pipeline group resource within HuaweiCloud.

## Example Usage

```hcl
variable "codearts_project_id" {}
variable "name" {}

resource "huaweicloud_codearts_pipeline_group" "test" {
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

* `name` - (Required, String) Specifies the group name.

* `parent_id` - (Optional, String, NonUpdatable) Specifies the group parent ID.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `children` - Indicates the child group name list.

* `ordinal` - Indicates the group ordinal

* `path_id` - Indicates the path ID.

* `create_time` - Indicates the create time.

* `update_time` - Indicates the update time.

* `creator` - Indicates the creator ID.

* `updater` - Indicates the updater ID.

## Import

The group can be imported using `project_id` and `id` separated by a slash, e.g.

```bash
$ terraform import huaweicloud_codearts_pipeline_group.test <project_id>/<id>
```
