---
subcategory: "Cloud Firewall (CFW)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cfw_address_group"
description: ""
---

# huaweicloud_cfw_address_group

Manages a CFW IP address group resource within HuaweiCloud.

## Example Usage

```hcl
variable "name" {}
variable "description" {}

data "huaweicloud_cfw_firewalls" "test" {}

resource "huaweicloud_cfw_address_group" "test" {
  object_id   = data.huaweicloud_cfw_firewalls.test.records[0].protect_objects[0].object_id
  name        = var.name
  description = var.description
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `object_id` - (Required, String, ForceNew) Specifies the protected object ID.

  Changing this parameter will create a new resource.

* `name` - (Required, String) Specifies the IP address group name.

* `description` - (Optional, String) Specifies the Address group description.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `address_type` - The address type. The value can be **0** (IPv4) or **1** (IPv6).

## Import

The ipaddressgroup can be imported using the `id`, e.g.

```bash
$ terraform import huaweicloud_cfw_address_group.test 0ce123456a00f2591fabc00385ff1234
```
