---
subcategory: "VPN"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_vpn_metrics"
description: |-
  Use this data source to get the list of VPN metrics within HuaweiCloud.
---

# huaweicloud_vpn_metrics

Use this data source to get the list of VPN metrics within HuaweiCloud.

## Example Usage

### Basic Usage

```hcl
data "huaweicloud_vpn_metrics" "test" {
  namespace = "SYS.VPN"
}
```

### Filter by metric name

```hcl
data "huaweicloud_vpn_metrics" "test" {
  namespace   = "SYS.VPN"
  metric_name = "gateway_connection_num"
}
```

### Filter by dimension

```hcl
data "huaweicloud_vpn_metrics" "test" {
  namespace   = "SYS.VPN"
  metric_name = "gateway_connection_num"
  dim         = ["p2c_vpn_gateway_id,4ae00000-0000-0000-0000-0000b1428ef7"]
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the metrics.
  If omitted, the provider-level region will be used.

* `namespace` - (Optional, String) Specifies the service metric namespace.
  The namespace for enterprise VPN is **SYS.VPN**.

* `metric_name` - (Optional, String) Specifies the metric ID.

* `dim` - (Optional, List) Specifies the metric dimensions.
  The format is `key,value`. For example: `p2c_vpn_gateway_id,4ae00000-0000-0000-0000-0000b1428ef7`.
  Multiple dimensions are supported, with a maximum of 3 levels.

* `order` - (Optional, String) Specifies the sorting method of the results.
  The valid values are as follows:
  + **asc** - Ascending order
  + **desc** - Descending order

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `metrics` - The list of VPN metrics.
  The [metrics](#vpn_metrics) structure is documented below.

<a name="vpn_metrics"></a>
The `metrics` block supports:

* `namespace` - The namespace of the metric.

* `metric_name` - The metric name.

* `unit` - The metric unit.

* `dimensions` - The list of metric dimensions.
  The [dimensions](#vpn_metrics_dimensions) structure is documented below.

<a name="vpn_metrics_dimensions"></a>
The `dimensions` block supports:

* `name` - The dimension name.

* `value` - The dimension value.
