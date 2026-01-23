---
subcategory: "Resource Access Manager (RAM)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_ram_permission_versions"
description: |
  Use this data source to get the list of all versions with specified permission in Resource Access Manager.
---

# huaweicloud_ram_permission_versions

Use this data source to get the list of all versions with specified permission in Resource Access Manager.

## Example Usage

```hcl
variable "permission_id" {}

data "huaweicloud_ram_permission_versions" "test" {
  permission_id = var.permission_id
}
```

## Argument Reference

The following arguments are supported:

* `permission_id` - (Required, String) Specifies the ID of the shared resource permission.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `permissions` - The detailed information list of shared resource permissions.

  The [permissions](#RAM_permissions) structure is documented below.

<a name="RAM_permissions"></a>
The `permissions` block supports:

* `id` - Indicates the ID of RAM permission.

* `name` - Indicates the name of RAM permission.

* `resource_type` - Indicates the resource type of RAM permission.

* `is_resource_type_default` - Whether the RAM permission resource type is default.

* `created_at` - Indicates the RAM permission creation time.

* `updated_at` - Indicates the RAM permission last update time.

* `permission_urn` - Indicates the URN for the permission.

* `permission_type` - Indicates the permission type.

* `default_version` - Indicates whether the current version is the default version.

* `version` - Indicates the version of the permission.

* `status` - Indicates the status of the permission.
