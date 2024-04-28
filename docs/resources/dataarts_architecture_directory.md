---
subcategory: "DataArts Studio"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dataarts_architecture_directory"
description: ""
---

# huaweicloud_dataarts_architecture_directory

Manages DataArts Architecture directory resource within HuaweiCloud.

## Example Usage

```hcl
variable "workspace_id" {}
variable "name" {}
variable "directory_id" {}

resource "huaweicloud_dataarts_architecture_directory" "test-root" {
  workspace_id = var.workspace_id
  name         = var.name
  type         = "STANDARD_ELEMENT"
}

resource "huaweicloud_dataarts_architecture_directory" "test-sub" {
  workspace_id = var.workspace_id
  name         = var.name
  type         = "STANDARD_ELEMENT"
  parent_id    = var.directory_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to manage the directory.
  Changing this creates a new directory.

* `workspace_id` - (Required, String, ForceNew) Specifies the workspace ID which the directory in.
  Changing this creates a new directory

* `name` - (Required, String) Specifies the directory name.

* `type` - (Required, String) Specifies the directory type. The valid values are **STANDARD_ELEMENT** and **CODE**.

* `description` - (Optional, String) Specifies the description of directory.

* `parent_id` - (Optional, String) Specifies the parent ID of the directory.
  It's **Required** when you created a subordinate directory.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `root_id` - The root directory ID.

* `qualified_name` - The directory path. Format is `<root_directory_name>.<sub_directory_name1>.<sub_directory_name2>...`

* `created_at` - The create time of the directory.

* `updated_at` - The update time of the directory.

* `created_by` - The person creating the directory.

* `updated_by` - The person updating the directory.

* `children` - The name list of sub-directory.

## Import

DataArts Architecture directory can be imported using `<workspace_id>/<type>/<qualified_name>`, e.g.

```sh
terraform import huaweicloud_dataarts_architecture_directory.test b606cd4a47b645108a122857204b360f/STANDARD_ELEMENT/root
```
