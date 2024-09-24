---
subcategory: "Anti-DDoS"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_antiddos_weekly_protection_statistics"
description: |-
  Use this data source to query Cloud Native Anti-DDos weekly protection statistics within HuaweiCloud.
---

# huaweicloud_antiddos_weekly_protection_statistics

Use this data source to query Cloud Native Anti-DDos weekly protection statistics within HuaweiCloud.

-> Only supports querying weekly statistical data within four weeks.

## Example Usage

```hcl
data "huaweicloud_antiddos_weekly_protection_statistics" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `period_start_date` - (Optional, String) Specifies the start date of the seven-day period, the value is a timestamp
  in milliseconds.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID in UUID format.

* `ddos_intercept_times` - The number of DDoS attacks blocked in a week.

* `weekdata` - The number of attacks in a week.

  The [weekdata](#weekdata_struct) structure is documented below.

* `top10` - Top `10` attacked IP addresses.

  The [top10](#top10_struct) structure is documented below.

<a name="weekdata_struct"></a>
The `weekdata` block supports:

* `period_start_date` - The start date of the seven-day period, the value is a timestamp in milliseconds.

* `ddos_intercept_times` - The number of DDoS attacks blocked.

* `ddos_blackhole_times` - The number of DDoS black holes.

* `max_attack_bps` - The maximum attack traffic.

* `max_attack_conns` - The maximum number of attack connections.

<a name="top10_struct"></a>
The `top10` block supports:

* `floating_ip_address` - The Elastic IP address.

* `times` - The number of DDoS attacks blocked, including scrubbing and black holes.
