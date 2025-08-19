---
subcategory: "Advanced Anti-DDoS"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_aad_unblock_ip"
description: |-
  Use this resource to unblock an IP address that has been blocked by the Advanced Anti-DDoS service.
---

# huaweicloud_aad_unblock_ip

Use this resource to unblock an IP address that has been blocked by the Advanced Anti-DDoS service.

-> This resource is a one-time action resource for unblocking IP addresses. Deleting this resource will not
re-block the IP address, but will only remove the resource information from the tfstate file.

## Example Usage

```hcl
variable "ip" {}

resource "huaweicloud_aad_unblock_ip" "test" {
  ip = var.ip
}
```

## Argument Reference

The following arguments are supported:

* `ip` - (Required, String, NonUpdatable) Specifies the IP address to be unblocked.

* `blocking_time` - (Optional, Int, NonUpdatable) Specifies the blocking timestamp in milliseconds.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID in UUID format.
