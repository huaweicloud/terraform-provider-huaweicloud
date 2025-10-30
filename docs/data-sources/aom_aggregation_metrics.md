---
subcategory: "Application Operations Management (AOM 2.0)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_aom_aggregation_metrics"
description: |-
  Use this data source to get the list of AOM aggregation metrics.
---

# huaweicloud_aom_aggregation_metrics

Use this data source to get the list of AOM aggregation metrics.

## Example Usage

```hcl
data "huaweicloud_aom_aggregation_metrics" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `service_metrics` - Indicates the aggregation metrics list.
  The [service_metrics](#service_metrics_struct) structure is documented below.

<a name="service_metrics_struct"></a>
The `service_metrics` block supports:

* `service` - Indicates the service name.

* `metrics` - Indicates the metrics list.
