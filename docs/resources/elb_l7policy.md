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

* `redirect_pool_id` - (Required, String) Requests matching this policy will be redirected to the pool with this ID.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The unique ID for the L7 policy.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 10 minute.
* `update` - Default is 10 minute.
* `delete` - Default is 10 minute.

## Import

ELB policy can be imported using the policy ID, e.g.

```
$ terraform import huaweicloud_elb_policy.policy_1 5c20fdad-7288-11eb-b817-0255ac10158b
```
