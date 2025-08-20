---
subcategory: "Cloud Native Anti-DDoS Advanced"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cnad_advanced_policy_ip_unbinding"
description: |-
  Manages a CNAD advanced policy IP unbinding resource within HuaweiCloud.
---

# huaweicloud_cnad_advanced_policy_ip_unbinding

Manages a CNAD advanced policy IP unbinding resource within HuaweiCloud.

-> This resource is a one-time action resource for unbinding IP addresses from a protection policy. Deleting this resource
   will not change the current CNAD policy IP binding, but will only remove the resource information from the tfstate file.

## Example Usage

```hcl
variable "policy_id" {}
variable "ip_list" {
  type = list(string)
}

resource "huaweicloud_cnad_advanced_policy_ip_unbinding" "test" {
  policy_id = var.policy_id
  ip_list   = var.ip_list
}
```

## Argument Reference

The following arguments are supported:

* `policy_id` - (Required, String, NonUpdatable) Specifies the ID of the protection policy to unbind IPs from.

* `ip_list` - (Required, List, NonUpdatable) Specifies the list of IP addresses to unbind from the policy.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID in UUID format.
