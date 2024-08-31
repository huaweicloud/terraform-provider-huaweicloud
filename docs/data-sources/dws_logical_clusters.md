---
subcategory: "GaussDB(DWS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dws_logical_clusters"
description: |-
  Use this data source to get the list of DWS logical clusters within HuaweiCloud.
---

# huaweicloud_dws_logical_clusters

Use this data source to get the list of DWS logical clusters within HuaweiCloud.

## Example Usage

```hcl
variable "dws_cluster_id" {}

data "huaweicloud_dws_logical_clusters" "test" {
  cluster_id = var.dws_cluster_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `cluster_id` - (Required, String) Specified the ID of the cluster to which the logical clusters belong.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `add_enable` - Whether the logical cluster can be added.

* `logical_clusters` - All logical clusters that match the filter parameters.

  The [logical_clusters](#logical_clusters_struct) structure is documented below.

<a name="logical_clusters_struct"></a>
The `logical_clusters` block supports:

* `id` - The ID of the logical cluster.

* `name` - The name of the logical cluster.

* `first_logical_cluster` - Whether it is the first logical cluster.

* `cluster_rings` - The list of logical cluster rings.

  The [cluster_rings](#logical_clusters_cluster_rings_struct) structure is documented below.

* `status` - The current status of the logical cluster.
  + **Failed**
  + **Normal**
  + **Unavailable**
  + **Redistribute**
  + **Redistribute_failed**
  + **Unbalanced**
  + **Stopped**

* `edit_enable` - Whether the logical cluster is allowed to be edited.

* `restart_enable` - Whether the logical cluster is allowed to be restarted.

* `delete_enable` - Whether the logical cluster is allowed to be deleted.

<a name="logical_clusters_cluster_rings_struct"></a>
The `cluster_rings` block supports:

* `ring_hosts` - The list of the cluster hosts.

  The [ring_hosts](#cluster_rings_ring_hosts_struct) structure is documented below.

<a name="cluster_rings_ring_hosts_struct"></a>
The `ring_hosts` block supports:

* `host_name` - The host name.

* `back_ip` - The backend IP address.

* `cpu_cores` - The number of CPU cores.

* `memory` - The host memory, in GB.

* `disk_size` - The host disk size, in GB.
