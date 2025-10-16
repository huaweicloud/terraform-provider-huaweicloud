---
subcategory: "Elastic Cloud Server (ECS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_compute_recycle_bin_server_recover"
description: |-
  Manages an ECS recycle bin server recover resource within HuaweiCloud.
---

# huaweicloud_compute_recycle_bin_server_recover

Manages an ECS recycle bin server recover resource within HuaweiCloud.

## Example Usage

```hcl
variable "server_id" {}

resource "huaweicloud_compute_recycle_bin_server_recover" "test" {
  server_id = var.server_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this creates a new resource.

* `server_id` - (Required, String, NonUpdatable) Specifies the VM ID.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID. The value is the instance ID.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 10 minutes.
