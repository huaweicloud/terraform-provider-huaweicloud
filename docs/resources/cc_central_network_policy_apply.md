---
subcategory: "Cloud Connect (CC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cc_central_network_policy_apply"
description: ""
---

# huaweicloud_cc_central_network_policy_apply

Apply a central network policy of Cloud Connect within HuaweiCloud.
Only one policy can be applied to a central network. If you need to change the policy, apply a new policy.
The previously applied policy will be automatically canceled.

## Example Usage

```hcl
variable "central_network_id" {}
variable "policy_id" {}

resource "huaweicloud_cc_central_network_policy_apply" "test" {
  central_network_id = var.central_network_id
  policy_id          = var.policy_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `central_network_id` - (Required, String, ForceNew) Central network ID.

  Changing this parameter will create a new resource.

* `policy_id` - (Required, String) Policy ID.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID which equals the `central_network_id`.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 10 minutes.
* `update` - Default is 10 minutes.
* `delete` - Default is 10 minutes.

## Import

The central network policy apply can be imported using `central_network_id` and `policy_id`, separated by a slash, e.g.

```bash
$ terraform import huaweicloud_cc_central_network_policy_apply.test <central_network_id>/<policy_id>
```
