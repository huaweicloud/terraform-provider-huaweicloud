---
subcategory: "Resource Access Manager (RAM)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_ram_resource_share_permission"
description: |-
  Manages a RAM resource share permission resource within HuaweiCloud.
---

# huaweicloud_ram_resource_share_permission

Manages a RAM resource share permission resource within HuaweiCloud.

## Example Usage

```hcl
variable "resource_share_id" {}
variable "permission_id" {}

resource "huaweicloud_ram_resource_share_permission" "test" {
  resource_share_id = var.resource_share_id
  permission_id     = var.permission_id
}
```

## Argument Reference

The following arguments are supported:

* `resource_share_id` - (Required, String, NonUpdatable) Specifies the ID of the resource share.

* `permission_id` - (Required, String, NonUpdatable) Specifies the ID of the permission.

* `replace` - (Optional, Bool, NonUpdatable) Specify a specific permission to replace or bind to the existing resource
  type associated with the shared resource instance.
  Setting it to **true** will replace the current permissions for the same resource type.
  Setting it to **false** will bind the permission to the current resource type.
  The default value is **false**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `permission_name` - Indicates the name of RAM permission.

* `resource_type` - Indicates the resource type of RAM permission.

* `status` - Indicates the status of the permission.

* `created_at` - Indicates the RAM permission create time.

* `updated_at` - Indicates the RAM permission last update time.

## Import

The RAM resource share permission can be imported using the resource share ID and the permission ID
separated by a slash, e.g.

```bash
$ terraform import huaweicloud_ram_resource_share_permission.test <resource_share_id>/<permission_id>
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response, security or some other reason. The missing attribute includes: `replace`.
It is generally recommended running `terraform plan` after importing the permission. You can then decide
if changes should be applied to the permission, or the resource definition should be updated to align with the permission.
Also you can ignore changes as below.

```hcl
resource "huaweicloud_ram_resource_share_permission" "test" {
    ...

  lifecycle {
    ignore_changes = [
      replace,
    ]
  }
}
```
