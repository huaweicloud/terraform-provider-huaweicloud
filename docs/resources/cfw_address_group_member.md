---
subcategory: "Cloud Firewall (CFW)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cfw_address_group_member"
description: |-
  Manages a CFW IP address group member resource within HuaweiCloud.
---

# huaweicloud_cfw_address_group_member

Manages a CFW IP address group member resource within HuaweiCloud.

## Example Usage

```hcl
variable "group_id" {}
variable "address" {}

resource "huaweicloud_cfw_address_group_member" "test" {
  group_id = var.group_id
  address  = var.address
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `group_id` - (Required, String, ForceNew) Specifies the ID of the IP address group.

  Changing this parameter will create a new resource.

* `address` - (Required, String, ForceNew) Specifies the IP address.

  Changing this parameter will create a new resource.

* `address_type` - (Optional, Int, ForceNew) Specifies the address type.
  The value can be **0** (IPv4) or **1** (IPv6).

  Changing this parameter will create a new resource.

* `description` - (Optional, String, ForceNew) Specifies address description.

  Changing this parameter will create a new resource.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `address_type` - The address type. The value can be **0** (IPv4) or **1** (IPv6).

## Import

The CFW IP address group member can be imported using `group_id`, `id`, separated by a slash, e.g.

```bash
$ terraform import huaweicloud_cfw_address_group_member.test <group_id>/<id>
```
