---
subcategory: "Resource Access Manager (RAM)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_ram_resource_share_associated_permissions"
description: |-
  Use this data source to get the list of shared resource permissions associated with the RAM shared resources.
---

# huaweicloud_ram_resource_share_associated_permissions

Use this data source to get the list of shared resource permissions associated with the RAM shared resources.

## Example Usage

```hcl
variable "resource_share_id " {}

data "huaweicloud_ram_resource_share_associated_permissions" "test" {
  resource_share_id  = var.resource_share_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `resource_share_id` - (Required, String) Specifies the ID of the resource share.

* `permission_name` - (Optional, String) Specifies the name of the RAM managed permission.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `associated_permissions` - The list of RAM managed permissions associated with the resource share.

  The [associated_permissions](#associated_permissions_struct) structure is documented below.

<a name="associated_permissions_struct"></a>
The `associated_permissions` block supports:

* `permission_id` - The permission ID.

* `permission_name` - The name of the RAM managed permission.

* `resource_type` - The resource type to which the permission applies.

* `status` - The status of the permission.

* `created_at` - The time when the permission was created.

* `updated_at` - The time when the permission was last updated.
