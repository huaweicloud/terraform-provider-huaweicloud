---
subcategory: "Cloud Connect (CC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cc_cloud_connection_capabilities"
description: |-
  Use this data source to get the list of CC cloud connection capabilities.
---

# huaweicloud_cc_cloud_connection_capabilities

Use this data source to get the list of CC cloud connection capabilities.

## Example Usage

```hcl
data "huaweicloud_cc_cloud_connection_capabilities" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `resource_type` - (Optional, String) Specifies the resource type.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `capabilities` - Indicates the cloud connection capabilities list.

  The [capabilities](#capabilities_struct) structure is documented below.

<a name="capabilities_struct"></a>
The `capabilities` block supports:

* `id` - Indicates the instance ID.

* `description` - Indicates the resource description.

* `created_at` - Indicates the time when the resource was created.
  The UTC time is in the **yyyy-MM-ddTHH:mm:ss** format.

* `updated_at` - Indicates the time when the resource was updated.
  The UTC time must be in the **yyyy-MM-ddTHH:mm:ss** format.

* `resource_type` - Indicates the resource type.
  The value can be:
  **v2**: V2 APIs
  **v3**: V3 APIs
  **billing_mode_period_reduce**: specification downgrade (yearly/monthly subscription
  **billing_mode_demand**: pay-per-use billing
  **bwp95**: pay-per-use billing (95th percentile billing)
  **bwp95Avg**: pay-per-use billing (daily 95th percentile peak billing)
  **network-quality**: monitoring of packet loss rate and network latency
  **er**: support for enterprise routers or not
  **domain_bandwidth**: tenant bandwidth
  **ipv6**: support for IPv6 or not
  **ipv6_support_regions**: regions where IPv6 is supported

* `bandwidth` - Indicates the bandwidth.

  The [bandwidth](#capabilities_bandwidth_struct) structure is documented below.

* `support_regions` - Indicates the list of regions available to a tenant.

<a name="capabilities_bandwidth_struct"></a>
The `bandwidth` block supports:

* `min` - Indicates the minimum bandwidth.

* `max` - Indicates the maximum bandwidth.
