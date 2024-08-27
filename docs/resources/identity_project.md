---
subcategory: "Identity and Access Management (IAM)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_identity_project"
description: ""
---

# huaweicloud_identity_project

Manages an IAM project resource within HuaweiCloud.

-> **NOTE:** You *must* have admin privileges to use this resource.

!>  Deleting projects is not supported. The project is only removed from the state, but it remains in the cloud.

## Example Usage

```hcl
resource "huaweicloud_identity_project" "project_1" {
  name        = "cn-north-1_project1"
  description = "This is a test project"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required, String) Specifies the name of the project. it must start with an existing *region* and be less
  than or equal to 64 characters. Example: cn-north-1_project1.

* `status` - (Optional, String) Specifies the status of the project.
  Valid values are **normal** and **suspended**, default is **normal**.

* `description` - (Optional, String) Specifies the description of the project.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - A resource ID in UUID format.

* `parent_id` - The parent of the IAM project.

* `enabled` - Whether the IAM project is enabled.

## Import

IAM projects can be imported using the `id`, e.g.

```bash
$ terraform import huaweicloud_identity_project.project_1 89c60255-9bd6-460c-822a-e2b959ede9d2
```
