---
subcategory: "Resource Access Manager (RAM)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_ram_resource_permissions"
description: ""
---

# huaweicloud_ram_resource_permissions

Use this data source to get the list of RAM permissions.

## Example Usage

```hcl
variable "resource_type" {}
variable "name" {}

data "huaweicloud_ram_resource_permissions" "test" {
  resource_type = var.resource_type
  name          = var.name
}
```

## Argument Reference

The following arguments are supported:

* `resource_type` - (Optional, String) Specifies the resource type of RAM permission in which to query the data source.
  Valid values are **vpc:subnets**, **dns:zone** and **dns:resolverRule**.

* `name` - (Optional, String) Specifies the name of RAM permission in which to query the data source.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `permissions` - Indicates the list of the RAM permissions
  The [Permissions](#RAM_Permissions) structure is documented below.

<a name="RAM_Permissions"></a>
The `Permissions` block supports:

* `id` - Indicates the id of RAM permission.

* `name` - Indicates the name of RAM permission.

* `resource_type` - Indicates the resource type of RAM permission.

* `is_resource_type_default` - Whether the RAM permission resource type is default.

* `created_at` - Indicates the RAM permission create time.

* `updated_at` - Indicates the RAM permission last update time.
