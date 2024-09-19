---
subcategory: "GaussDB(DWS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dws_cluster_nodes"
description: |-
  Use this data source to get the list of nodes under specified DWS cluster within HuaweiCloud.
---

# huaweicloud_dws_cluster_nodes

Use this data source to get the list of nodes under specified DWS cluster within HuaweiCloud.

## Example Usage

### Use node ID to query a ring to which the node belongs

```hcl
variable "dws_cluster_id" {}
variable "node_id" {}

data "huaweicloud_dws_cluster_nodes" "test" {
  cluster_id = var.dws_cluster_id
  node_id    = var.node_id
}
```

### Use resource status to query all idel nodes under the specified DWS cluster

```hcl
variable "dws_cluster_id" {}

data "huaweicloud_dws_cluster_nodes" "test" {
  cluster_id = var.dws_cluster_id
  filter_by  = "instCreateType"
  filter     = "NODE"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `cluster_id` - (Required, String) Specifies the DWS cluster ID.

* `node_id` - (Optional, String) Specifies the ID of the node.
  If you specify this parameter, the query result is a ring containing the node.

* `filter_by` - (Optional, String) Specifies the query filter criteria.  
  The valid values are as follows:
  + **status**
  + **instCreateType**

* `filter` - (Optional, String) Specifies the type corresponding to the `filter_by` parameter.
  + If the `filter_by` is set to `status`, the valid values are **ALL**, **FREE**, **ACTIVE**, **FAILED**, **UNKNOWN**,
  **CREATE_FAILED**, **DELETE_FAILED** and **STOPPED**.
  + If the `filter_by` is set to `instCreatetype`, the valid values are **ALL**, **INST** (idle) and **NODE** (used).

  -> 1. If `filter` is set to **All**, it means to query all nodes, including deleted historical nodes.
  <br/>2. If not specified, it means to query existing all nodes.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `nodes` - All nodes that match the filter parameters.

  The [nodes](#nodes_struct) structure is documented below.

<a name="nodes_struct"></a>
The `nodes` block supports:

* `id` - The ID of the node.

* `name` - The name of the node.

* `status` - The current status of the node.

* `sub_status` - The sub-status of the node.
  + **READ**: The ECS on the node ready.
  + **PREPAPED**: The node software has been installed.
  + **INITED**: The cluster has been created.
  + **CREATED**: The node has been created.

* `spec` - The specification of the node.

* `inst_create_type` - The occupancy status of nodes by the cluster.
  If the value is **NODE**, it indicates that the node is idle.
  If the value is empty, it indicates that the node has been used.

* `alias_name` - The alias of the node.

* `availability_zone` - The availability zone of the node.
