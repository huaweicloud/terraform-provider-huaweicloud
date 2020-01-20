---
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_lb_whitelist_v2"
sidebar_current: "docs-huaweicloud-resource-lb-whitelist-v2"
description: |-
  Manages a Load Balancer whitelist resource within HuaweiCloud.
---

# huaweicloud\_lb\_whitelist\_v2

Manages a Load Balancer whitelist resource within HuaweiCloud.

## Example Usage

```hcl
resource "huaweicloud_lb_listener_v2" "listener_1" {
  name            = "listener_1"
  protocol        = "HTTP"
  protocol_port   = 8080
  loadbalancer_id = var.loadbalancer_id
}

resource "huaweicloud_lb_whitelist_v2" "whitelist_1" {
  enable_whitelist = true
  whitelist        = "192.168.11.1,192.168.0.1/24,192.168.201.18/8"
  listener_id      = huaweicloud_lb_listener_v2.listener_1.id
}
```

## Argument Reference

The following arguments are supported:

* `tenant_id` - (Optional) Required for admins. The UUID of the tenant who owns
    the whitelist. Only administrative users can specify a tenant UUID
    other than their own. Changing this creates a new whitelist.

* `listener_id` - (Required) The Listener ID that the whitelist will be associated with. Changing this creates a new whitelist.

* `enable_whitelist` - (Optional) Specify whether to enable access control.

* `whitelist` - (Optional) Specifies the IP addresses in the whitelist. Use commas(,) to separate
    the multiple IP addresses.

## Attributes Reference

The following attributes are exported:

* `id` - The unique ID for the whitelist.
* `tenant_id` - See Argument Reference above.
* `listener_id` - See Argument Reference above.
* `enable_whitelist` - See Argument Reference above.
* `whitelist` - See Argument Reference above.
