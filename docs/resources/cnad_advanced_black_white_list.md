---
subcategory: "Cloud Native Anti-DDoS Advanced"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cnad_advanced_black_white_list"
description: ""
---

# huaweicloud_cnad_advanced_black_white_list

Manages a CNAD advanced policy black and white IP list resource within HuaweiCloud.

## Example Usage

```hcl
variable "policy_id" {}
variable "black_ip_list" {
  type = list(string)
}
variable "white_ip_list" {
  type = list(string)
}

resource "huaweicloud_cnad_advanced_black_white_list" "test" {
  policy_id     = var.policy_id
  black_ip_list = var.black_ip_list
  white_ip_list = var.white_ip_list
}
```

## Argument Reference

The following arguments are supported:

* `policy_id` - (Required, String, ForceNew) Specifies the CNAD advanced policy ID in which to add black and white IP
  list.

  Changing this parameter will create a new resource.

* `black_ip_list` - (Optional, List) Specifies the black IP list.

* `white_ip_list` - (Optional, List) Specifies the white IP list.

  -> For `black_ip_list` and `white_ip_list`, the value could be IPv4/IPv6 address or CIDR.
  The IP address or IP address range must be unique. For each IP address and IP address range, the minimum length is 7,
  the maximum length is 128.
  The total number of `black_ip_list` and `white_ip_list` cannot exceed the maximum allowed, the default maximum is 200.
  At least one of `black_ip_list` or `white_ip_list` must be configured.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

## Import

The CNAD advanced policy black and white IP list can be imported using the `policy_id`, e.g.

```bash
$ terraform import huaweicloud_cnad_advanced_black_white_list.test <policy_id>
```
