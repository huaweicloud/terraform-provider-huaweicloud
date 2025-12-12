---
subcategory: "Resource Access Manager (RAM)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_ram_resource_permission"
description: |
  Use this data source to get the permission for the specified version.
---

# huaweicloud_ram_resource_permission

Use this data source to get the permission for the specified version.

## Example Usage

```hcl
variable "permission_id" {}
variable "permission_version" {}

data "huaweicloud_ram_resource_permission" "test" {
  resource_type      = var.permission_id
  permission_version = var.permission_version
}
```

## Argument Reference

The following arguments are supported:

* `permission_id` - (Required, String) Specifies the id of RAM permission.

* `permission_version` - (Optional, String) Specifies the version of RAM permission.
  Not setting this value will return the default version permission.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `permission` - Indicates the list of the RAM permissions
  The [permission](#RAM_Permission) structure is documented below.

<a name="RAM_Permission"></a>
The `permission` block supports:

* `id` - Indicates the id of RAM permission.

* `name` - Indicates the name of RAM permission.

* `resource_type` - Indicates the resource type of RAM permission.

* `is_resource_type_default` - Whether the RAM permission resource type is default.

* `content` - Impact and actions allowed by the permission.

* `created_at` - Indicates the RAM permission create time.

* `updated_at` - Indicates the RAM permission last update time.

* `permission_urn` - Indicates the URN for the permission.

* `permission_type` - Indicates the permission type.

* `default_version` - Indicates whether the current version is the default version.

* `version` - Indicates the version of the permission.

* `status` - Indicates the status of the permission.
