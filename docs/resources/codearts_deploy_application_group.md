---
subcategory: "CodeArts Deploy"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_codearts_deploy_application_group"
description: |-
  Manages a CodeArts deploy application group resource within HuaweiCloud.
---

# huaweicloud_codearts_deploy_application_group

Manages a CodeArts deploy application group resource within HuaweiCloud.

## Example Usage

### Create a top-level group

```hcl
variable "project_id" {}
variable "name" {}

resource "huaweicloud_codearts_deploy_application_group" "test" {
  project_id = var.project_id
  name       = var.name
}
```

### Create a sub-level group

```hcl
variable "project_id" {}
variable "name" {}
variable "parent_id" {}

resource "huaweicloud_codearts_deploy_application_group" "test" {
  project_id = var.project_id
  name       = var.name
  parent_id  = var.parent_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.
  Changing this creates a new resource.

* `project_id` - (Required, String, ForceNew) Specifies the project ID for CodeArts service.
  Changing this creates a new resource.

* `parent_id` - (Optional, String, ForceNew) Specifies the parent application group ID.
  Changing this creates a new resource.

* `name` - (Required, String) Specifies the application group name.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `ordinal` - Indicates the group sorting field.

* `path` - Indicates the group path.

* `application_count` - Indicates the total number of applications in the group.

* `created_by` - Indicates the ID of the group creator.

* `updated_by` - Indicates the ID of the user who last updates the group.

* `children` - Indicates the child group name list.

## Import

The application group can be imported using the `project_id` and `id`, separated by a slash, e.g.

```bash
$ terraform import huaweicloud_codearts_deploy_application_group.test <project_id>/<id>
```
