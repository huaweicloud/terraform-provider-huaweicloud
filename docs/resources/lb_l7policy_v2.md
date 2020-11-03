---
subcategory: "Elastic Load Balance (ELB)"
---

# huaweicloud\_lb\_l7policy\_v2

Manages a Load Balancer L7 Policy resource within HuaweiCloud.

## Example Usage

```hcl
resource "huaweicloud_lb_loadbalancer_v2" "loadbalancer_1" {
  name          = "loadbalancer_1"
  vip_subnet_id = "SUBNET_ID"
}

resource "huaweicloud_lb_listener_v2" "listener_1" {
  name            = "listener_1"
  protocol        = "HTTP"
  protocol_port   = 8080
  loadbalancer_id = huaweicloud_lb_loadbalancer_v2.loadbalancer_1.id
}

resource "huaweicloud_lb_pool_v2" "pool_1" {
  name            = "pool_1"
  protocol        = "HTTP"
  lb_method       = "ROUND_ROBIN"
  loadbalancer_id = huaweicloud_lb_loadbalancer_v2.loadbalancer_1.id
}

resource "huaweicloud_lb_l7policy_v2" "l7policy_1" {
  name             = "test"
  action           = "REDIRECT_TO_POOL"
  description      = "test l7 policy"
  position         = 1
  listener_id      = huaweicloud_lb_listener_v2.listener_1.id
  redirect_pool_id = huaweicloud_lb_pool_v2.pool_1.id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional) The region in which to obtain the V2 Networking client.
    A Networking client is needed to create a L7 Policy resource. If omitted, the
    `region` argument of the provider is used. Changing this creates a new
    L7 Policy.

* `tenant_id` - (Optional) Required for admins. The UUID of the tenant who owns
    the L7 Policy.  Only administrative users can specify a tenant UUID
    other than their own. Changing this creates a new L7 Policy.

* `name` - (Optional) Human-readable name for the L7 Policy. Does not have
    to be unique.

* `description` - (Optional) Human-readable description for the L7 Policy.

* `action` - (Required) The L7 Policy action - can either be REDIRECT\_TO\_POOL,
    or REDIRECT\_TO\_LISTENER. Changing this creates a new L7 Policy.

* `listener_id` - (Required) The Listener on which the L7 Policy will be associated with.
    Changing this creates a new L7 Policy.

* `position` - (Optional) The position of this policy on the listener. Positions start at 1. Changing this creates a new L7 Policy.

* `redirect_pool_id` - (Optional) Requests matching this policy will be redirected to the
    pool with this ID. Only valid if action is REDIRECT\_TO\_POOL.

* `redirect_listener_id` - (Optional) Requests matching this policy will be redirected to the listener with this ID. Only valid if action is REDIRECT\_TO\_LISTENER.

* `admin_state_up` - (Optional) The administrative state of the L7 Policy.
    This value can only be true (UP).

## Attributes Reference

The following attributes are exported:

* `id` - The unique ID for the L7 {olicy.
* `region` - See Argument Reference above.
* `tenant_id` - See Argument Reference above.
* `name` - See Argument Reference above.
* `description` - See Argument Reference above.
* `action` - See Argument Reference above.
* `listener_id` - See Argument Reference above.
* `position` - See Argument Reference above.
* `redirect_pool_id` - See Argument Reference above.
* `redirect_listener_id` - See Argument Reference above.
* `admin_state_up` - See Argument Reference above.

## Import

Load Balancer L7 Policy can be imported using the L7 Policy ID, e.g.:

```
$ terraform import huaweicloud_lb_l7policy_v2.l7policy_1 8a7a79c2-cf17-4e65-b2ae-ddc8bfcf6c74
```
