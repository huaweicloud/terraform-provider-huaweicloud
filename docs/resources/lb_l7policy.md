---
subcategory: "Elastic Load Balance (ELB)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_lb_l7policy"
description: ""
---

# huaweicloud_lb_l7policy

Manages an ELB L7 Policy resource within HuaweiCloud.

## Example Usage

```hcl
resource "huaweicloud_lb_loadbalancer" "loadbalancer_1" {
  name          = "loadbalancer_1"
  vip_subnet_id = var.subnet_id
}

resource "huaweicloud_lb_listener" "listener_1" {
  name            = "listener_1"
  protocol        = "HTTP"
  protocol_port   = 8080
  loadbalancer_id = huaweicloud_lb_loadbalancer.loadbalancer_1.id
}

resource "huaweicloud_lb_pool" "pool_1" {
  name            = "pool_1"
  protocol        = "HTTP"
  lb_method       = "ROUND_ROBIN"
  loadbalancer_id = huaweicloud_lb_loadbalancer.loadbalancer_1.id
}

resource "huaweicloud_lb_l7policy" "l7policy_1" {
  name             = "test"
  action           = "REDIRECT_TO_POOL"
  description      = "test l7 policy"
  position         = 1
  listener_id      = huaweicloud_lb_listener.listener_1.id
  redirect_pool_id = huaweicloud_lb_pool.pool_1.id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) The region in which to create the L7 Policy resource. If omitted, the
  provider-level region will be used. Changing this creates a new L7 Policy.

* `name` - (Optional, String) Human-readable name for the L7 Policy. Does not have to be unique.

* `description` - (Optional, String) Human-readable description for the L7 Policy.

* `listener_id` - (Required, String, ForceNew) Specifies the ID of the listener for which the forwarding policy is added.
  Changing this creates a new L7 Policy.

* `action` - (Required, String, ForceNew) Specifies whether requests are forwarded to another backend server group
  or redirected to an HTTPS listener. Changing this creates a new L7 Policy. The value ranges:
  + **REDIRECT_TO_POOL**: Requests are forwarded to the backend server group specified by `redirect_pool_id`.
  + **REDIRECT_TO_LISTENER**: Requests are redirected from the HTTP listener specified by `listener_id` to the
    HTTPS listener specified by `redirect_listener_id`.

* `position` - (Optional, Int, ForceNew) The position of this policy on the listener. Positions start at 1.
  Changing this creates a new L7 Policy.

* `redirect_pool_id` - (Optional, String) Specifies the ID of the backend server group to which traffic is forwarded.
  This parameter is mandatory when `action` is set to **REDIRECT_TO_POOL**. The backend server group must meet the
  following requirements:
  + Cannot be the default backend server group of the listener.
  + Cannot be the backend server group used by forwarding policies of other listeners.

* `redirect_listener_id` - (Optional, String) Specifies the ID of the listener to which the traffic is redirected.
  This parameter is mandatory when `action` is set to **REDIRECT_TO_LISTENER**. The listener must meet the
  following requirements:
  + Can only be an HTTPS listener.
  + Can only be a listener of the same load balancer.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The unique ID for the L7 policy.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 10 minutes.
* `update` - Default is 10 minutes.
* `delete` - Default is 10 minutes.

## Import

Load Balancer L7 Policy can be imported using the L7 Policy ID, e.g.:

```bash
$ terraform import huaweicloud_lb_l7policy.l7policy_1 8a7a79c2-cf17-4e65-b2ae-ddc8bfcf6c74
```
