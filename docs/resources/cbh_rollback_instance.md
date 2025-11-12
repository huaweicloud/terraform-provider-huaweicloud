---
subcategory: "Cloud Bastion Host (CBH)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cbh_rollback_instance"
description: |-
  Manages a CBH rollback instance resource within HuaweiCloud.
---

# huaweicloud_cbh_rollback_instance

Manages a CBH rollback instance resource within HuaweiCloud.

-> This resource is only a one-time action resource to rollback a CBH instance. Deleting this resource
will not clear the corresponding rollback instance, but will only remove the resource information from the tfstate file.

## Example Usage

```hcl
variable "server_id" {}

resource "huaweicloud_cbh_rollback_instance" "test" {
  server_id = var.server_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to rollback a CBH instance.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `server_id` - (Required, String, NonUpdatable) Specifies the ID of the CBH instance to rollback.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID (same as `server_id`).
