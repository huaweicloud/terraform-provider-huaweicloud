---
subcategory: "Anti-DDoS"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_antiddos_eip_defense_statuses"
description: |-
  Use this data source to query the list of defense statuses of EIPs within HuaweiCloud.
---

# huaweicloud_antiddos_eip_defense_statuses

Use this data source to query the list of defense statuses of EIPs within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_antiddos_eip_defense_statuses" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `status` - (Optional, String) Specifies the defense status. Valid values are **normal**, **configging**, **notConfig**,
  **packetcleaning**, and **packetdropping**. Query all by default.

* `ip` - (Optional, String) Specifies the IP address. Both IPv4 and IPv6 addresses are supported. For example, if you
  enter **?ip=192.168**, the defense status of EIPs corresponding to **192.168.111.1** and **10.192.168.8** is returned.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `ddos_status` - The list of Anti-DDos statuses.

  The [ddos_status](#ddos_status_struct) structure is documented below.

<a name="ddos_status_struct"></a>
The `ddos_status` block supports:

* `blackhole_endtime` - The end time of black hole.

* `protect_type` - The protect type.

* `traffic_threshold` - The traffic cleaning threshold in Mbps.

* `http_threshold` - The threshold of http traffic.

* `eip_id` - The ID of an EIP.

* `public_ip` - The public address of the EIP.

* `network_type` - The EIP type. Valid values are:
  + **EIP**: EIP bound or not bound to ECS.
  + **ELB**: EIP bound to ELB.

* `status` - The defense status.
