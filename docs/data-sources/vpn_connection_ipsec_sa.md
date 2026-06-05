---
subcategory: "VPN"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_vpn_connection_ipsec_sa"
description: |-
  Use this data source to get the list of VPN connection IPSec SA within HuaweiCloud.
---

# huaweicloud_vpn_connection_ipsec_sa

Use this data source to get the list of VPN connection IPSec SA within HuaweiCloud.

## Example Usage

```hcl
variable "vpn_connection_id" {}

data "huaweicloud_vpn_connection_ipsec_sa" "test" {
  vpn_connection_id = var.vpn_connection_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `vpn_connection_id` - (Required, String) Specifies the VPN connection ID.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `sa_infos` - The list of IPSec SA information.
  The [sa_infos](#ipsec_sa_infos) structure is documented below.

<a name="ipsec_sa_infos"></a>
The `sa_infos` block supports:

* `id` - The network negotiation ID.

* `source_ip_cidr` - The source IP CIDR.

* `dest_ip_cidr` - The destination IP CIDR.

* `packets_sent` - The number of packets sent.

* `packets_recv` - The number of packets received.

* `traffic_sent` - The traffic sent in bytes.

* `traffic_recv` - The traffic received in bytes.

* `collected_at` - The data collection time in UTC format.
