---
subcategory: "GaussDB(DWS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dws_cluster_topo_rings"
description: |-
  Use this data source to get the list of topology rings under specified DWS cluster within HuaweiCloud.
---

# huaweicloud_dws_cluster_topo_rings

Use this data source to get the list of topology rings under specified DWS cluster within HuaweiCloud.

## Example Usage

```hcl
variable "dws_cluster_id" {}

data "huaweicloud_dws_cluster_topo_rings" "test" {
  cluster_id = var.dws_cluster_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `cluster_id` - (Required, String) Specifies the DWS cluster ID to which the topology rings belong.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `rings` - The list of the topology rings under DWS cluster.

  The [rings](#rings_struct) structure is documented below.

<a name="rings_struct"></a>
The `rings` block supports:

* `instances` - The list of the cluster instances.

  The [instances](#rings_instances_struct) structure is documented below.

<a name="rings_instances_struct"></a>
The `instances` block supports:

* `id` - The ID of the instance.

* `name` - The name of the instance.

* `eip_address` - The EIP address corresponding to the instance.

* `elb_address` - The ELB address corresponding to the instance

* `status` - The current status of the instance.
  + **200**: Available.
  + **300**: Unavailable.
  + **302**: Deletion failed.
  + **303**: Creation failed.
  + **400**: Deleted.
  + **800**: Frozen.
  + **900**: Stopped.

* `manage_ip` - The management IP address of the instance.

* `internal_ip` - The internal communication IP address of the instance.

* `internal_mgnt_ip` - The internal management IP address of the instance.

* `traffic_ip` - The server IP address of the instance.

* `availability_zone` - The availability zone of the instance.
