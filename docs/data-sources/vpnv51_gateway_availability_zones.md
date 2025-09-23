---
subcategory: "Virtual Private Network (VPN)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_vpnv51_gateway_availability_zones"
description: |-
  Use this data source to get the list of VPN gateway availability zones.
---

# huaweicloud_vpnv51_gateway_availability_zones

Use this data source to get the list of VPN gateway availability zones.

## Example Usage

```hcl
data "huaweicloud_vpnv51_gateway_availability_zones" "test"{}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `availability_zones` - Indicates the availability zones.

  The [availability_zones](#availability_zones_struct) structure is documented below.

<a name="availability_zones_struct"></a>
The `availability_zones` block supports:

* `name` - Indicates the AZ name.

* `public_border_group` - Indicates the common boundary group.

* `available_specs` - Indicates the available specs.

  The [available_specs](#availability_zones_available_specs_struct) structure is documented below.

<a name="availability_zones_available_specs_struct"></a>
The `available_specs` block supports:

* `flavor` - Indicates the VPN gateway flavor.

* `attachment_type` - Indicates the attachment type.

* `ip_version` - Indicates the IP version of VPN gateway.
