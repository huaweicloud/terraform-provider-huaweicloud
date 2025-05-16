---
subcategory: "Advanced Anti-DDoS"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_aad_black_white_list"
description: |-
  Manages a black white list resource of Advanced Anti-DDos service within HuaweiCloud.
---

# huaweicloud_aad_forward_rule

Manages a black white list resource of Advanced Anti-DDos service within HuaweiCloud.

## Example Usage

### Add whitelist

```hcl
variable "instance_id" {}

resource "huaweicloud_aad_black_white_list" "white_test" {
  instance_id = var.instance_id
  type        = "white"
  ips         = ["11.1.2.114", "11.1.2.115"]
}
```

### Add blacklist

```hcl
variable "instance_id" {}

resource "huaweicloud_aad_black_white_list" "black_test" {
  instance_id = var.instance_id
  type        = "black"
  ips         = ["11.1.2.112", "11.1.2.113"]
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String, NonUpdatable) Specifies the AAD instance ID.

* `type` - (Required, String, NonUpdatable) Specifies the rule type. Valid values are **black** and **white**.

* `ips` - (Required, List) Specifies the IP address list.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID (also `instance_id`).

## Import

The AAD black white list resource can be imported using the `instance_id` and `type`, separated by a slash (/), e.g.

```bash
$ terraform import huaweicloud_aad_black_white_list.test <instance_id>/<type>
```
