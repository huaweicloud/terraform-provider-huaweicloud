---
subcategory: "Direct Connect (DC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dc_global_gateway"
description: |-
  Manages a DC global gateway resource within HuaweiCloud.
---

# huaweicloud_dc_global_gateway

Manages a DC global gateway resource within HuaweiCloud.

## Example Usage

```hcl
variable "name" {}

resource "huaweicloud_dc_global_gateway" "test" {
  name           = var.name
  description    = "test description"
  bgp_asn        = 10
  address_family = "ipv4"

  tags = {
    foo = "bar"
    key = "value"
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this will create new resource.

* `name` - (Required, String) Specifies the name of the global DC gateway.

* `description` - (Optional, String) Specifies the description of the global DC gateway.

* `address_family` - (Optional, String) Specifies the IP address family of the global DC gateway. Valid values are:
  + **ipv4**: Only IPv4 is supported.
  + **dual**: Both IPv4 and IPv6 are supported.

  Defaults to **ipv4**.

* `bgp_asn` - (Optional, Int, NonUpdatable) Specifies the BGP ASN of the global DC gateway. Valid value is limited from `1`
  to `4,294,967,295`. Defaults to `64,512`.

* `enterprise_project_id` - (Optional, String, NonUpdatable) Specifies the enterprise project ID that the global DC gateway
  belongs to. For enterprise users, if omitted, default enterprise project will be used.

* `tags` - (Optional, Map, NonUpdatable) Specifies the key/value pairs to associate with the DC global gateway.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `reason` - The cause of the failure to create the global DC gateway.

* `global_center_network_id` - The ID of the central network that the global DC gateway is added to.

* `location_name` - The location where the underlying device of the global DC gateway is deployed.

* `locales` - The locale address description information.
  The [locales](#global_gateway_locales) structure is documented below.

* `current_peer_link_count` - The number of peer links allowed on a global DC gateway, indicating the number of
  enterprise routers that the global DC gateway can be attached to.

* `available_peer_link_count` - The number of peer links that can be created for a global DC gateway.

* `status` - The status of the global DC gateway.

* `created_time` - The time when the global DC gateway was created.

* `updated_time` - The time when the global DC gateway was updated.

* `all_tags` - The all key/value pairs to associate with the DC global gateway.

<a name="global_gateway_locales"></a>
The `locales` block supports:

* `en_us` - The region name in English.

* `zh_cn` - The region name in Chinese.

## Import

The DC global gateway resource can be imported using the `id`, e.g.

```bash
$ terraform import huaweicloud_dc_global_gateway.test <id>
```
