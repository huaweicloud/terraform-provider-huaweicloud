---
subcategory: "GaussDB(DWS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dws_logical_cluster"
description: |-
  Manages a GaussDB(DWS) logical cluster resource within HuaweiCloud.
---

# huaweicloud_dws_logical_cluster

Manages a GaussDB(DWS) logical cluster resource within HuaweiCloud.

-> **NOTE:** The first DWS logical cluster can't be deleted. When performing a delete operation, the resource will only
be removed from the state, but it remains in the cloud.

## Example Usage

```hcl
variable "cluster_id" {}
variable "ring_hosts" {
  type = list(object({
    host_name = string
    back_ip   = string
    cpu_cores = number
    memory    = number
    disk_size = number
  }))
}

resource "huaweicloud_dws_logical_cluster" "test" {
  logical_cluster_name = "test_name"
  cluster_id           = var.cluster_id

  cluster_rings {
    dynamic "ring_hosts" {
      for_each = var.ring_hosts
      content {
        host_name = ring_hosts.value.host_name
        back_ip   = ring_hosts.value.back_ip
        cpu_cores = ring_hosts.value.cpu_cores
        memory    = ring_hosts.value.memory
        disk_size = ring_hosts.value.disk_size
      }
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `cluster_id` - (Required, String, ForceNew) Specifies the DWS cluster ID.

  Changing this parameter will create a new resource.

* `logical_cluster_name` - (Required, String, ForceNew) Specifies the logical cluster name. Changing this parameter will
  create a new resource. Only letters, digits, and underscores (_) are allowed. The maximum length is 63 characters.
  The name must be unique and cannot be the keywords `group_version1`, `group_version2`, `group_version3`,
  `installation`, `elastic_group`, `optimal`, and `query`.

* `cluster_rings` - (Required, List, ForceNew) Specifies the DWS logical cluster ring list information.
  Changing this parameter will create a new resource.
The [cluster_rings](#LogicalCluster_ClusterRings) structure is documented below.

<a name="LogicalCluster_ClusterRings"></a>
The `cluster_rings` block supports:

* `ring_hosts` - (Required, List, ForceNew) Specifies the cluster host ring information. All host information of a ring
  must be specified. Changing this parameter will create a new resource.
The [ring_hosts](#LogicalCluster_RingHosts) structure is documented below.

<a name="LogicalCluster_RingHosts"></a>
The `ring_hosts` block supports:

* `host_name` - (Required, String, ForceNew) Specifies the host name. Changing this parameter will create a new resource.

* `back_ip` - (Required, String, ForceNew) Specifies the backend IP address. Changing this parameter will create a new resource.

* `cpu_cores` - (Required, Int, ForceNew) Specifies the number of CPU cores. Changing this parameter will create a new resource.

* `memory` - (Required, Float, ForceNew) Specifies the host memory. Changing this parameter will create a new resource.

* `disk_size` - (Required, Float, ForceNew) Specifies the host disk size. Changing this parameter will create a new resource.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `status` - The DWS logical cluster status.

* `first_logical_cluster` - Whether it is the first logical cluster. The first logical cluster cannot be deleted.

* `edit_enable` - Whether editing is allowed.

* `restart_enable` - Whether to allow restart.

* `delete_enable` - Whether deletion is allowed.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 10 minutes.
* `delete` - Default is 10 minutes.

## Import

The DWS logical cluster resource can be imported using the `cluster_id` and `id`, separated by a slash, e.g.

```bash
$ terraform import huaweicloud_dws_logical_cluster.test <cluster_id>/<id>
```
