---
subcategory: "VPN"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_vpn_connection_peer_configuration_export"
description: |-
  Use this data source to export VPN connection peer configuration within HuaweiCloud.
---

# huaweicloud_vpn_connection_peer_configuration_export

Use this data source to export VPN connection peer configuration within HuaweiCloud.

## Example Usage

```hcl
variable "vpn_connection_id" {}

data "huaweicloud_vpn_connection_peer_configuration_export" "test" {
  vpn_connection_id = var.vpn_connection_id
  vendor            = "Huawei"
  type              = "USG"
  model             = "USG6655F"
  version           = "default"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `vpn_connection_id` - (Required, String) Specifies the VPN connection ID.

* `vendor` - (Required, String) Specifies the vendor of the peer device.

* `type` - (Required, String) Specifies the series of the peer device.

* `model` - (Required, String) Specifies the model of the peer device.

* `version` - (Required, String) Specifies the version of the peer device.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `peer_config` - The peer device configuration information.
