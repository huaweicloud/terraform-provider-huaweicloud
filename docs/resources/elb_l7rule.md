---
subcategory: "Dedicated Load Balance (Dedicated ELB)"
---

# huaweicloud_elb_l7rule

Manages an ELB L7 Rule resource within HuaweiCloud.

## Example Usage

```hcl
variable l7policy_id {}

resource "huaweicloud_elb_l7rule" "l7rule_1" {
  l7policy_id  = var.l7policy_id
  type         = "PATH"
  compare_type = "EQUAL_TO"
  value        = "/api"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) The region in which to create the L7 Rule resource. If omitted, the
  provider-level region will be used. Changing this creates a new L7 Rule.

* `type` - (Required, String, ForceNew) The L7 Rule type - can either be HOST_NAME or PATH. Changing this creates a new
  L7 Rule.

* `compare_type` - (Required, String) The comparison type for the L7 rule - can either be STARTS_WITH, EQUAL_TO or REGEX

* `l7policy_id` - (Required, String, ForceNew) The ID of the L7 Policy. Changing this creates a new L7 Rule.

* `value` - (Required, String) The value to use for the comparison.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The unique ID for the L7 Rule.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 10 minute.
* `update` - Default is 10 minute.
* `delete` - Default is 10 minute.

## Import

ELB L7 rule can be imported using the L7 policy ID and L7 rule ID separated by a slash, e.g.

```
$ terraform import huaweicloud_elb_rule.rule_1 e0bd694a-abbe-450e-b329-0931fd1cc5eb/4086b0c9-b18c-4d1c-b6b8-4c56c3ad2a9e
```
