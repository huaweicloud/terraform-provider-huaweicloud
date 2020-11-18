---
subcategory: "Elastic Load Balance (ELB)"
---

# huaweicloud\_lb\_l7rule

Manages an ELB L7 Rule resource within HuaweiCloud.
This is an alternative to `huaweicloud_lb_l7rule_v2`

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
  name         = "test"
  action       = "REDIRECT_TO_URL"
  description  = "test description"
  position     = 1
  listener_id  = huaweicloud_lb_listener.listener_1.id
  redirect_url = "http://www.example.com"
}

resource "huaweicloud_lb_l7rule" "l7rule_1" {
  l7policy_id  = huaweicloud_lb_l7policy.l7policy_1.id
  type         = "PATH"
  compare_type = "EQUAL_TO"
  value        = "/api"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional) The region in which to create the L7 Rule resource.
    If omitted, the provider-level region will be used as default.
    Changing this creates a new L7 Rule.

* `tenant_id` - (Optional) Required for admins. The UUID of the tenant who owns
    the L7 Rule.  Only administrative users can specify a tenant UUID
    other than their own. Changing this creates a new L7 Rule.

* `description` - (Optional) Human-readable description for the L7 Rule.

* `type` - (Required) The L7 Rule type - can either be HOST\_NAME or PATH. Changing this creates a new L7 Rule.

* `compare_type` - (Required) The comparison type for the L7 rule - can either be
    STARTS\_WITH, EQUAL_TO or REGEX

* `l7policy_id` - (Required) The ID of the L7 Policy to query. Changing this creates a new
    L7 Rule.

* `value` - (Required) The value to use for the comparison. For example, the file type to
    compare.

* `key` - (Optional) The key to use for the comparison. For example, the name of the cookie to
    evaluate. Valid when `type` is set to COOKIE or HEADER. Changing this creates a new L7 Rule.

* `admin_state_up` - (Optional) The administrative state of the L7 Rule.
    The value can only be true (UP).

## Attributes Reference

The following attributes are exported:

* `id` - The unique ID for the L7 Rule.
* `region` - See Argument Reference above.
* `tenant_id` - See Argument Reference above.
* `type` - See Argument Reference above.
* `compare_type` - See Argument Reference above.
* `l7policy_id` - See Argument Reference above.
* `value` - See Argument Reference above.
* `key` - See Argument Reference above.
* `invert` - See Argument Reference above.
* `admin_state_up` - See Argument Reference above.
* `listener_id` - The ID of the Listener owning this resource.

## Timeouts
This resource provides the following timeouts configuration options:
- `create` - Default is 10 minute.
- `update` - Default is 10 minute.
- `delete` - Default is 10 minute.

## Import

Load Balancer L7 Rule can be imported using the L7 Policy ID and L7 Rule ID
separated by a slash, e.g.:

```
$ terraform import huaweicloud_lb_l7rule.l7rule_1 e0bd694a-abbe-450e-b329-0931fd1cc5eb/4086b0c9-b18c-4d1c-b6b8-4c56c3ad2a9e
```
