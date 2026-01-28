---
subcategory: "Resource Access Manager (RAM)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_ram_resource_permission"
description: |-
  Use this data source to get the RAM resource permission details within HuaweiCloud.
---

# huaweicloud_ram_resource_permission

Use this data source to get the RAM resource permission details within HuaweiCloud.

## Example Usage

```hcl
variable "permission_id" {}

data "huaweicloud_ram_resource_permission" "test" {
  permission_id = var.permission_id
}
```

## Argument Reference

The following arguments are supported:

* `permission_id` - (Required, String) Specifies the ID of RAM permission.

* `permission_version` - (Optional, String) Specifies the version of RAM permission.
  Not setting this value will return the default version permission.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `permission` - Indicates the permission details.

  The [permission](#permission_struct) structure is documented below.

<a name="permission_struct"></a>
The `permission` block supports:

* `id` - Indicates the ID of RAM permission.

* `name` - Indicates the name of RAM permission.

* `resource_type` - Indicates the resource type of RAM permission.

* `content` - Impact and actions allowed by the permission.

* `is_resource_type_default` - Whether the RAM permission resource type is default.

* `created_at` - Indicates the RAM permission creation time.

* `updated_at` - Indicates the RAM permission last update time.

* `permission_urn` - Indicates the URN for the permission.

* `permission_type` - Indicates the permission type.

* `default_version` - Indicates whether the current version is the default version.

* `version` - Indicates the version of the permission.

* `status` - Indicates the status of the permission.
