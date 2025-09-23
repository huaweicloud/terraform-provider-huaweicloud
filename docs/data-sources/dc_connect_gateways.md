---
subcategory: "Direct Connect (DC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dc_connect_gateways"
description: |-
  Use this data source to get the list of DC connect gateways.
---

# huaweicloud_dc_connect_gateways

Use this data source to get the list of DC connect gateways.

## Example Usage

```hcl
data "huaweicloud_dc_global_gateways" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `connect_gateway_id` - (Optional, List) Specifies the IDs of the DC connect gateway.

* `name` - (Optional, List) Specifies the names of the DC connect gateway.

* `sort_key` - (Optional, String) Specifies the sorting field.
  Defaults to **id**.

* `sort_dir` - (Optional, List) Specifies the sorting order of returned results.
  Value options: **asc (default)** and **desc**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `connect_gateways` - Indicates the list of connect gateway.

  The [connect_gateways](#connect_gateways_struct) structure is documented below.

<a name="connect_gateways_struct"></a>
The `connect_gateways` block supports:

* `id` - Indicates the ID of the connect gateway.

* `name` - Indicates the name of the connect gateway.

* `description` - Indicates the description of the DC connect gateway.

* `address_family` - Indicates the address family.
  The value can be:
  **ipv4**: Only IPv4 is supported.
  **dual**: IPv4 and IPv6 are supported.

* `status` - Indicates the dc connect gateway status.
  **DOWN**: The DC connect gateway is not in use or the associated device goes down.
  **ACTIVE**: The DC connect gateway is normal.
  **ERROR**: The DC connect gateway is abnormal.

* `gateway_site` - Indicates the gateway location.

* `gcb_id` - Indicates the global connection bandwidth ID.

* `bgp_asn` - Indicates the BGP ASN.

* `current_geip_count` - Indicates the number of global EIPs bound to the connect gateway.

* `access_site` - Indicates the access site of the connect gateway.

* `created_time` - Indicates the time when the connect gateway was created.

* `updated_time` - Indicates the time when the connect gateway was updated.
