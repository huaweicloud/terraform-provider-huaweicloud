---
subcategory: "Advanced Anti-DDoS"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_aad_ip_ddos_statistics"
description: |-
  Use this data source to get the Advanced Anti-DDos IP ddos statistics within HuaweiCloud.
---

# huaweicloud_aad_ip_ddos_statistics

Use this data source to get the Advanced Anti-DDos IP ddos statistics within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_aad_ip_ddos_statistics" "test" {}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String) Specifies the instance ID.

* `ip` - (Required, String) Specifies the IP address.

* `start_time` - (Required, String) Specifies the start time, millisecond timestamp.

* `end_time` - (Required, String) Specifies the end time, millisecond timestamp.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID in UUID format.

* `attack_kbps_peak` - The attack peak.

* `in_kbps_peak` - The traffic peak.

* `ddos_count` - The number of attacks.

* `timestamp` - The attack peak timestamp.

* `vip` - The Anti-DDoS IP.
