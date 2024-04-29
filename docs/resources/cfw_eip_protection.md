---
subcategory: "Cloud Firewall (CFW)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cfw_eip_protection"
description: ""
---

# huaweicloud_cfw_eip_protection

Manages the protected EIPs under the protect object for CFW service within HuaweiCloud.

~> A protection object (`object_id`) can only create one `huaweicloud_cfw_eip_protection` resource for managing
protected EIPs.

## Example Usage

```hcl
variable "object_id" {}
variable "protected_eips" {
  type = list(object({
    id           = string
    ipv4_address = string
  }))
}

resource "huaweicloud_cfw_eip_protection" "test" {
  object_id = var.object_id

  dynamic "protected_eip" {
    for_each = var.protected_eips

    content {
      id          = protected_eip.value["id"]
      public_ipv4 = protected_eip.value["ipv4_address"]
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `object_id` - (Required, String, ForceNew) The protected object ID.
  Changing this parameter will create a new resource.

* `protected_eip` - (Required, List) The protected EIP configurations.
  The [object](#cfw_protected_eip) structure is documented below.

<a name="cfw_protected_eip"></a>
The `protected_eip` block supports:

* `id` - (Required, String) The ID of the protected EIP.

* `public_ipv4` - (Optional, String) The IPv4 address of the protected EIP.

* `public_ipv6` - (Optional, String) The IPv6 address of the protected EIP.

-> At least one of `public_ipv4` and `public_ipv6` must be set.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID (also protected object ID).

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 5 minutes.
* `update` - Default is 5 minutes.
* `delete` - Default is 5 minutes.

## Import

The protection resource can be imported using their `object_id` or `id`, e.g.

```bash
$ terraform import huaweicloud_cfw_eip_protection.test <id>
```
