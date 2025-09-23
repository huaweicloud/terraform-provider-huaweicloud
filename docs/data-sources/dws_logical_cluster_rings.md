---
subcategory: "GaussDB(DWS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dws_logical_cluster_rings"
description: |-
  Use this data source to get the list of DWS logical cluster rings.
---

# huaweicloud_dws_logical_cluster_rings

Use this data source to get the list of DWS logical cluster rings.

## Example Usage

```hcl
variable "cluster_id" {}

data "huaweicloud_dws_logical_cluster_rings" "test" {
  cluster_id = var.cluster_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `cluster_id` - (Required, String) Specifies the DWS cluster ID.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `cluster_rings` - Indicates the cluster ring list information.
  The [cluster_rings](#LogicalClusterRings_ClusterRings) structure is documented below.

<a name="LogicalClusterRings_ClusterRings"></a>
The `cluster_rings` block supports:

* `is_available` - Indicates whether the cluster host ring is available. Only host rings with this field set to **true**
  can be used to create logical clusters.

* `ring_hosts` - Indicates the cluster host ring list information.
  The [ring_hosts](#LogicalClusterRings_ClusterRingsRingHosts) structure is documented below.

<a name="LogicalClusterRings_ClusterRingsRingHosts"></a>
The `ring_hosts` block supports:

* `host_name` - Indicates the host name.

* `back_ip` - Indicates the backend IP address.

* `cpu_cores` - Indicates the number of CPU cores.

* `memory` - Indicates the host memory.

* `disk_size` - Indicates the host disk size. The unit is GB.
