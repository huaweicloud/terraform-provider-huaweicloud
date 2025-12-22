---
subcategory: "Workspace"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_workspace_desktop_pool_expand"
description: |-
  Use this resource to expand a desktop pool within HuaweiCloud.
---

# huaweicloud_workspace_desktop_pool_expand

Use this resource to expand a desktop pool within HuaweiCloud.

~> If you use this resource, please use `lifecycle.ignore_changes` to ignore the changes of `size`
   in `huaweicloud_workspace_desktop_pool`.

-> This resource is a one-time action resource for expanding a desktop pool. Deleting this resource will
   not clear the corresponding request record, but will only remove the resource information from the tfstate file.

## Example Usage

### Basic usage of desktop pool expansion

```hcl
variable "pool_id" {}

resource "huaweicloud_workspace_desktop_pool_expand" "test" {
  pool_id = var.pool_id
  size    = 5
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region where the desktop pool is located.  
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `pool_id` - (Required, String, NonUpdatable) Specifies the ID of the desktop pool to be expanded.

* `size` - (Required, Int, NonUpdatable) Specifies the number of desktops to be added to the pool.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is `15` minutes.
