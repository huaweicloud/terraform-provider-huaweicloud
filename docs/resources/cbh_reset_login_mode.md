---
subcategory: "Cloud Bastion Host (CBH)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cbh_reset_login_mode"
description: |-
  Manages a CBH reset login mode resource within HuaweiCloud.
---

# huaweicloud_cbh_reset_login_mode

Resets a CBH instance login mode resource within HuaweiCloud.

-> This resource is only a one-time action resource to reset a CBH instance login mode. Deleting this resource
will not clear the corresponding login mode, but will only remove the resource information from the tfstate file.

## Example Usage

```hcl
variable "server_id" {}

resource "huaweicloud_cbh_reset_login_mode" "test" {
  server_id = var.server_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to reset the login mode.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `server_id` - (Required, String, NonUpdatable) Specifies the ID of the CBH instance to reset the login mode.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID (same as `server_id`).
