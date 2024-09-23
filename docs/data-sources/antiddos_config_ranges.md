---
subcategory: "Anti-DDoS"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_antiddos_config_ranges"
description: |-
  Use this data source to query Cloud Native Anti-DDos config ranges within HuaweiCloud.
---

# huaweicloud_antiddos_config_ranges

Use this data source to query Cloud Native Anti-DDos config ranges within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_antiddos_config_ranges" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID in UUID format.

* `traffic_limited_list` - The list of traffic limits.

  The [traffic_limited_list](#traffic_limited_list_struct) structure is documented below.

* `http_limited_list` - The list of HTTP limits.

  The [http_limited_list](#http_limited_list_struct) structure is documented below.

* `connection_limited_list` - The list of limits of numbers of connections.

  The [connection_limited_list](#connection_limited_list_struct) structure is documented below.

* `extend_ddos_config` - The list of extend ddos limits.

  The [extend_ddos_config](#extend_ddos_config_struct) structure is documented below.

<a name="traffic_limited_list_struct"></a>
The `traffic_limited_list` block supports:

* `traffic_pos_id` - The position ID of traffic.

* `traffic_per_second` - The threshold of traffic per second (Mbit/s).

* `packet_per_second` - The threshold of number of packets per second.

<a name="http_limited_list_struct"></a>
The `http_limited_list` block supports:

* `http_request_pos_id` - The position ID of number of HTTP requests.

* `http_packet_per_second` - The threshold of number of HTTP requests per second.

<a name="connection_limited_list_struct"></a>
The `connection_limited_list` block supports:

* `new_connection_limited` - The number of new connections of a source IP address.

* `total_connection_limited` - The total number of connections of a source IP address.

* `cleaning_access_pos_id` - The position ID of access limit during cleaning.

<a name="extend_ddos_config_struct"></a>
The `extend_ddos_config` block supports:

* `traffic_per_second` - The threshold of traffic per second (Mbit/s).

* `packet_per_second` - The threshold of number of packets per second.

* `set_id` - The position ID of config.

* `new_connection_limited` - The number of new connections of a source IP address.

* `total_connection_limited` - The total number of connections of a source IP address.

* `http_packet_per_second` - The threshold of number of HTTP requests per second.
