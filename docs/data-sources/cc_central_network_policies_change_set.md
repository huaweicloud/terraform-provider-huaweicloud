---
subcategory: "Cloud Connect (CC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cc_central_network_policies_change_set"
description: |-
  Use this data source to get the list of central network policy change sets.
---

# huaweicloud_cc_central_network_policies_change_set

Use this data source to get the list of central network policy change sets.

## Example Usage

```hcl
variable "central_network_id" {}
variable "central_network_policy_id" {}

data "huaweicloud_cc_central_network_policy_change_sets" "test" {
  central_network_id        = var.central_network_id
  central_network_policy_id = var.central_network_policy_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `policy_id` - (Required, String) Specifies central network policy ID.

* `central_network_id` - (Required, String) Specifies central network ID.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `central_network_policy_change_set` - The central network policy change set.

  The [central_network_policy_change_set](#central_network_policy_change_set_struct) structure is documented below.

<a name="central_network_policy_change_set_struct"></a>
The `central_network_policy_change_set` block supports:

* `change_content` - The central network policy change set content.
