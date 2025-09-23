---
subcategory: "Cloud Native Anti-DDoS Advanced"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cnad_advanced_policy_ip_binding"
description: |-
  Manages a CNAD advanced policy IP binding resource within HuaweiCloud.
---

# huaweicloud_cnad_advanced_policy_ip_binding

Manages a CNAD advanced policy IP binding resource within HuaweiCloud.

-> This resource is a one-time action resource for binding IP addresses to a protection policy. Deleting this resource
   will not remove the IP bindings from the policy, but will only remove the resource information from the tfstate file.

## Example Usage

```hcl
variable "policy_id" {}
variable "ip_list" {
  type = list(string)
}

resource "huaweicloud_cnad_advanced_policy_ip_binding" "test" {
  policy_id = var.policy_id
  ip_list   = var.ip_list
}
```

## Argument Reference

The following arguments are supported:

* `policy_id` - (Required, String, NonUpdatable) Specifies the ID of the protection policy.

* `ip_list` - (Required, List, NonUpdatable) Specifies the list of IP addresses to bind the policy.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID in UUID format.
