---
subcategory: "Cloud Connect (CC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cc_central_network_connections"
description: ""
---

# huaweicloud_cc_central_network_connections

Use this data source to get the list of CC central network connections.

## Example Usage

```hcl
variable "central_network_id" {}
variable "central_network_connection_id" {}

data "huaweicloud_cc_central_network_connections" "test" {
  central_network_id = var.central_network_id
  connection_id      = var.central_network_connection_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `central_network_id` - (Required, String) Specifies the central network ID.

* `connection_id` - (Optional, String) Specifies the central network connection ID.

* `status` - (Optional, String) Specifies the central network connection status.

* `global_connection_bandwidth_id` - (Optional, String) Specifies the bandwidth package ID of the central network connection.

* `bandwidth_type` - (Optional, String) Specifies the bandwidth type of the central network connection.
  The bandwidth types are as follows:
  + **BandwidthPackage**: A global private bandwidth billed by fixed bandwidth is required, and cross-site connection
    bandwidths are assigned from the global private bandwidth.
  + **TestBandwidth**: Only the minimum bandwidth is used for testing cross-region connectivity.

* `type` - (Optional, String) Specifies the central network connection type.

* `is_cross_region` - (Optional, String) Specifies whether there are different regions.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `central_network_connections` - The list of the central network connections.

  The [central_network_connections](#central_network_connections_struct) structure is documented below.

<a name="central_network_connections_struct"></a>
The `central_network_connections` block supports:

* `id` - The central network connection ID.

* `name` - The central network connection name.

* `description` - The central network connection description.

* `enterprise_project_id` - The ID of the enterprise project that the virtual gateway belongs to.

* `central_network_id` - The the central network ID.

* `central_network_plane_id` - The plane ID of the enterprise router connection on the central network.

* `global_connection_bandwidth_id` - The bandwidth ID of the enterprise router connection on the central network.

* `bandwidth_type` - The bandwidth type of the enterprise router connection on the central network.

* `bandwidth_size` - The bandwidth size of the enterprise router connection on the central network.

* `status` - The central network connection status.

* `is_frozen` - Whether the central network connection is frozen.

* `type` - The type of the enterprise router connection on the central network.

* `connection_point_pair` - The both ends of a central network connection. The length is fixed to an array of 2.

  The [connection_point_pair](#central_network_connections_connection_point_pair_struct) structure is documented below.

* `created_at` - The creation time. The time is in the **yyyy-MM-ddTHH:mm:ss** format.

* `updated_at` - The update time. The time must be in the **yyyy-MM-ddTHH:mm:ss** format.

<a name="central_network_connections_connection_point_pair_struct"></a>
The `connection_point_pair` block supports:

* `id` - The point ID of a central network connection.

* `project_id` - The point project ID of a central network connection.

* `region_id` - The point region ID of central network connection.

* `site_code` - The site code of the point of central network connection.

* `instance_id` - The point instance ID of central network connection.

* `type` - The point type of central network connection.
