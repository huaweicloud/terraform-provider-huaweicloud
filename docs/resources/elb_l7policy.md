---
subcategory: "Dedicated Load Balance (Dedicated ELB)"
---

# huaweicloud_elb_l7policy

Manages an ELB L7 Policy resource within HuaweiCloud.

## Example Usage

```hcl
variable listener_id {}
variable pool_id {}

resource "huaweicloud_elb_l7policy" "policy_1" {
  name             = "policy_1"
  action           = "REDIRECT_TO_POOL"
  description      = "test description"
  listener_id      = var.listener_id
  redirect_pool_id = var.pool_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) The region in which to create the L7 Policy resource. If omitted, the
  provider-level region will be used. Changing this creates a new L7 Policy.

* `name` - (Optional, String) Human-readable name for the L7 Policy. Does not have to be unique.

* `description` - (Optional, String) Human-readable description for the L7 Policy.

* `listener_id` - (Required, String, ForceNew) The Listener on which the L7 Policy will be associated with. Changing
  this creates a new L7 Policy.

* `action` - (Optional, String, ForceNew) Specifies whether requests are forwarded to another backend server group
  or redirected to an HTTPS listener. Changing this creates a new L7 Policy. The value ranges:
  + **REDIRECT_TO_POOL**: Requests are forwarded to the backend server group specified by `redirect_pool_id`.
  + **REDIRECT_TO_LISTENER**: Requests are redirected from the HTTP listener specified by `listener_id` to the
    HTTPS listener specified by `redirect_listener_id`.
  Defaults to **REDIRECT_TO_POOL**.

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

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The unique ID for the L7 policy.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 10 minutes.
* `update` - Default is 10 minutes.
* `delete` - Default is 10 minutes.

## Import

ELB policy can be imported using the `id`, e.g.

```bash
$ terraform import huaweicloud_elb_policy.policy_1 <id>
```
