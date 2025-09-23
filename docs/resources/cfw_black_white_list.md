---
subcategory: "Cloud Firewall (CFW)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cfw_black_white_list"
description: ""
---

# huaweicloud_cfw_black_white_list

Manages a CFW black white list resource within HuaweiCloud.

## Example Usage

```hcl
variable "list_type" {}
variable "direction" {}
variable "address_type" {}
variable "address" {}
variable "protocol" {}
variable "port" {}

data "huaweicloud_cfw_firewalls" "test" {}

resource "huaweicloud_cfw_black_white_list" "test" {
  object_id    = data.huaweicloud_cfw_firewalls.test.records[0].protect_objects[0].object_id
  list_type    = var.list_type
  direction    = var.direction
  address_type = var.address_type
  address      = var.address
  protocol     = var.protocol
  port         = var.port
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `object_id` - (Required, String, ForceNew) Specifies the protected object ID.

  Changing this parameter will create a new resource.

* `list_type` - (Required, Int, ForceNew) Specifies the list type.
  The options are `4` (blacklist) and `5` (whitelist).

  Changing this parameter will create a new resource.

* `direction` - (Required, Int) Specifies the address direction.
  The options are `0` (source address) and `1` (destination address).

* `protocol` - (Required, Int) Specifies the protocol type. The value can be:
  + **6**: indicates TCP;
  + **17**: indicates UDP;
  + **1**: indicates ICMP;
  + **58**: indicates ICMPv6;
  + **-1**: indicates any protocol;

* `address_type` - (Required, Int) Specifies the IP address type.
  The options are `0` (ipv4), `1` (ipv6) and `2` (domain).

* `address` - (Required, String) Specifies the address.

* `port` - (Optional, String) Specifies the destination port.
  Required and only available if protocol is **TCP** or **UDP**.

* `description` - (Optional, String) Specifies the description of the list.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

## Import

The black whitelist can be imported using `object_id`, `list_type`, `address`, separated by slashes, e.g.

```bash
$ terraform import huaweicloud_cfw_black_white_list.test <object_id>/<list_type>/<address>
```
