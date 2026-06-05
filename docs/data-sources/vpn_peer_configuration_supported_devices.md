---
subcategory: "VPN"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_vpn_peer_configuration_supported_devices"
description: |-
  Use this data source to get the list of VPN peer configuration supported devices within HuaweiCloud.
---

# huaweicloud_vpn_peer_configuration_supported_devices

Use this data source to get the list of VPN peer configuration supported devices within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_vpn_peer_configuration_supported_devices" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `supported_devices` - The list of supported devices.
  The [supported_devices](#supported_devices_struct) structure is documented below.

<a name="supported_devices_struct"></a>
The `supported_devices` block supports:

* `vendor` - The vendor of the device.

* `type` - The series of the device.

* `model` - The model of the device.

* `version` - The version of the device.
