---
subcategory: "Direct Connect (DC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dc_connect_gateway"
description: |-
  Manages a DC connect gateway resource within HuaweiCloud.
---

# huaweicloud_dc_connect_gateway

Manages a DC connect gateway resource within HuaweiCloud.

## Example Usage

```hcl
variable "name" {}

resource "huaweicloud_dc_connect_gateway" "test" {
  name           = "test_connect_gateway_name"
  description    = "test description"
  address_family = "ipv4"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this will create new resource.

* `name` - (Required, String) Specifies the name of the DC connect gateway.

* `description` - (Optional, String) Specifies the description of the DC connect gateway.

* `address_family` - (Optional, String) Specifies the IP address family of the DC connect gateway. Value options:
  + **ipv4**: Only IPv4 is supported.
  + **dual**: Both IPv4 and IPv6 are supported.

  Defaults to **ipv4**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `status` - Indicates the status of the DC connect gateway.

* `access_site` - Indicates the access site of the connect gateway.

* `bgp_asn` - Indicates the BGP ASN.

* `current_geip_count` - Indicates the number of global EIPs bound to the connect gateway.

* `created_time` - Indicates the time when the connect gateway was created.

* `updated_time` - Indicates the time when the connect gateway was updated.

* `gcb_id` - Indicates the global connection bandwidth ID.

* `gateway_site` - Indicates the gateway location.

## Import

The DC connect gateway resource can be imported using the `id`, e.g.

```bash
$ terraform import huaweicloud_dc_connect_gateway.test <id>
```
