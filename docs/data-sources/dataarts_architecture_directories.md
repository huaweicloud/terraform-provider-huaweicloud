---
subcategory: "DataArts Studio"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dataarts_architecture_directories"
description: |-
  Use this data source to query DataArts Architecture directories within HuaweiCloud.
---

# huaweicloud_dataarts_architecture_directories

Use this data source to query DataArts Architecture directories within HuaweiCloud.

## Example Usage

```hcl
variable "workspace_id" {}

data "huaweicloud_dataarts_architecture_directories" "test" {
  workspace_id = var.workspace_id
  type         = "STANDARD_ELEMENT"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region where the directories are located.  
  If omitted, the provider-level region will be used.

* `workspace_id` - (Required, String) Specifies the ID of the workspace to which the directories belong.

* `type` - (Required, String) Specifies the type of directory to be queried.  
  The valid values are as follows:
  + **STANDARD_ELEMENT**
  + **CODE**

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `directories` - The list of directories that matched filter parameters.  
  The [directories](#dataarts_architecture_directories_attr) structure is documented below.

<a name="dataarts_architecture_directories_attr"></a>
The `directories` block supports:

* `id` - The ID of the directory.

* `name` - The name of the directory.

* `type` - The type of the directory.

* `description` - The description of the directory.

* `parent_id` - The parent ID of the directory.

* `root_id` - The root directory ID.

* `qualified_name` - The qualified name of the directory.

* `created_at` - The creation time of the directory, in RFC3339 format.

* `updated_at` - The latest update time of the directory, in RFC3339 format.

* `created_by` - The user who created the directory.

* `updated_by` - The user who updated the directory.

* `children` - The name list of sub-directories.
