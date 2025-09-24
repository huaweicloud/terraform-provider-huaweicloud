---
subcategory: "Elastic Cloud Server (ECS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_compute_kernel_dump_trigger"
description: |-
  Manages an ECS kernel dump trigger resource within HuaweiCloud.
---

# huaweicloud_compute_kernel_dump_trigger

Manages an ECS kernel dump trigger resource within HuaweiCloud.

## Example Usage

```hcl
variable "server_id" {}

resource "huaweicloud_compute_kernel_dump_trigger" "test" {
  server_id = var.server_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the auto launch group.
  If omitted, the provider-level region will be used. Changing this creates a new resource.

* `server_id` - (Required, String, NonUpdatable) Specifies the ECS ID.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID. The value is the instance ID.
