---
subcategory: "Cloud Phone (CPH)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cph_server_restart"
description: |-
  Manages a CPH server restart resource within HuaweiCloud.
---

# huaweicloud_cph_server_restart

Manages a CPH server restart resource within HuaweiCloud.

-> The current resource is a one-time resource, and destroying this resource is only removed from the state.

## Example Usage

```hcl
variable "server_id" {}

resource "huaweicloud_cph_server_restart" "test" {
  server_id = var.server_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.
  Changing this creates a new resource.

* `server_id` - (Optional, String, NonUpdatable) Specifies the ID of CPH server.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 30 minutes.
