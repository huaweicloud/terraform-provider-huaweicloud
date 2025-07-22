---
subcategory: "CodeArts Pipeline"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_codearts_pipeline_groups"
description: |-
  Use this data source to get a list of CodeArts pipeline groups.
---

# huaweicloud_codearts_pipeline_groups

Use this data source to get a list of CodeArts pipeline groups.

## Example Usage

```hcl
variable "codearts_project_id" {}

data "huaweicloud_codearts_pipeline_groups" "test" {
  project_id = var.codearts_project_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `project_id` - (Required, String) Specifies the CodeArts project ID.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `groups` - Indicates the pipeline groups list.
  The [groups](#attrblock--groups) structure is documented below.

<a name="attrblock--groups"></a>
The `groups` block supports:

* `id` - Indicates the group ID.

* `name` - Indicates the group name.

* `children` - Indicates the child group name list.

* `ordinal` - Indicates the group sorting field.

* `parent_id` - Indicates the parent group ID.

* `path_id` - Indicates the group path.

* `create_time` - Indicates the create time.

* `creator` - Indicates the ID of the group creator.

* `update_time` - Indicates the update time.

* `updater` - Indicates the ID of the user who last updates the group.
