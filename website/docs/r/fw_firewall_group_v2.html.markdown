---
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_fw_firewall_group_v2"
sidebar_current: "docs-huaweicloud-resource-fw-firewall-group-v2"
description: |-
  Manages a v2 firewall group resource within HuaweiCloud.
---

# huaweicloud\_fw\_firewall_group_v2

Manages a v2 firewall group resource within HuaweiCloud.

## Example Usage

```hcl
resource "huaweicloud_fw_rule_v2" "rule_1" {
  name             = "my-rule-1"
  description      = "drop TELNET traffic"
  action           = "deny"
  protocol         = "tcp"
  destination_port = "23"
  enabled          = "true"
}

resource "huaweicloud_fw_rule_v2" "rule_2" {
  name             = "my-rule-2"
  description      = "drop NTP traffic"
  action           = "deny"
  protocol         = "udp"
  destination_port = "123"
  enabled          = "false"
}

resource "huaweicloud_fw_policy_v2" "policy_1" {
  name = "my-policy"

  rules = ["${huaweicloud_fw_rule_v2.rule_1.id}",
    "${huaweicloud_fw_rule_v2.rule_2.id}",
  ]
}

resource "huaweicloud_fw_firewall_group_v2" "firewall_group_1" {
  name      = "my-firewall-group"
  ingress_policy_id = "${huaweicloud_fw_policy_v2.policy_1.id}"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional) The region in which to obtain the v2 networking client.
    A networking client is needed to create a firewall group. If omitted, the
    `region` argument of the provider is used. Changing this creates a new
    firewall group.

* `ingress_policy_id` - The ingress policy resource id for the firewall group. Changing
    this updates the `ingress_policy_id` of an existing firewall group.

* `egress_policy_id` - The egress policy resource id for the firewall group. Changing
    this updates the `egress_policy_id` of an existing firewall group.

* `name` - (Optional) A name for the firewall group. Changing this
    updates the `name` of an existing firewall group.

* `description` - (Required) A description for the firewall group. Changing this
    updates the `description` of an existing firewall group.

* `admin_state_up` - (Optional) Administrative up/down status for the firewall group
    (must be "true" or "false" if provided - defaults to "true").
    Changing this updates the `admin_state_up` of an existing firewall group.

* `tenant_id` - (Optional) The owner of the floating IP. Required if admin wants
    to create a firewall group for another tenant. Changing this creates a new
    firewall group.

* `ports` - (Optional) Port(s) to associate this firewall group instance
    with. Must be a list of strings. Changing this updates the associated routers
    of an existing firewall group.

* `value_specs` - (Optional) Map of additional options.

## Attributes Reference

The following attributes are exported:

* `region` - See Argument Reference above.
* `policy_id` - See Argument Reference above.
* `name` - See Argument Reference above.
* `description` - See Argument Reference above.
* `admin_state_up` - See Argument Reference above.
* `tenant_id` - See Argument Reference above.
* `ports` - See Argument Reference above.

## Import

Firewall Groups can be imported using the `id`, e.g.

```
$ terraform import huaweicloud_fw_firewall_group_v2.firewall_group_1 c9e39fb2-ce20-46c8-a964-25f3898c7a97
```
