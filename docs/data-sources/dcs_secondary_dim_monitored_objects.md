---
subcategory: "Distributed Cache Service (DCS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dcs_secondary_dim_monitored_objects"
description: |-
  Use this data source to query the list of monitored objects for a specific secondary dimension under DCS within HuaweiCloud.
---

# huaweicloud_dcs_secondary_dim_monitored_objects

Use this data source to query the list of monitored objects for a specific secondary dimension under DCS within
HuaweiCloud.

## Example Usage

### Basic Usage

```hcl
variable "instance_id" {}

data "huaweicloud_dcs_secondary_dim_monitored_objects" "test" {
  instance_id = var.instance_id
  dim_name    = "dcs_instance_id"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to query the monitored objects. If omitted, the
  provider-level region will be used.

* `instance_id` - (Required, String, NonUpdatable) Specifies the ID of the secondary dimension object (DCS instance ID).

* `dim_name` - (Required, String) Specifies the secondary dimension ID. Currently, only **dcs_instance_id** is
  supported.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID, formatted as `<instance_id>-<dim_name>`.

* `router` - The current query dimension route.

* `total` - The total number of monitoring objects for the secondary dimension.

* `children` - The list of sub-dimension objects for the current query dimension.
  The [children](#children_struct) structure is documented below.

* `instances` - The list of monitoring objects for the current query dimension.
  The [instances](#instances_struct) structure is documented below.

* `dcs_cluster_redis_node` - The list of monitoring objects for the data node dimension. Available for Proxy or Cluster
  instances.
  The [redis_nodes](#redis_nodes_struct) structure is documented below.

* `dcs_cluster_proxy_node` - The list of monitoring objects for the Proxy node dimension. Available for Redis 3.0 Proxy
  instances.
  The [proxy_nodes](#proxy_nodes_struct) structure is documented below.

* `dcs_cluster_proxy2_node` - The list of monitoring objects for the Proxy node dimension. Available for Redis 4.0 and
  5.0 Proxy instances.
  The [proxy2_nodes](#proxy2_nodes_struct) structure is documented below.

<a name="children_struct"></a>
The `children` block supports:

* `dim_name` - The dimension name.

* `dim_route` - The route of the dimension.

<a name="instances_struct"></a>
The `instances` block supports:

* `dcs_instance_id` - The measurement object ID, which is the instance ID.

* `name` - The measurement object name, which is the instance name.

* `status` - The measurement object status, which is the instance status.

<a name="redis_nodes_struct"></a>
The `redis_nodes` block supports:

* `dcs_instance_id` - The measurement object ID, which is the node ID.

* `name` - The measurement object name, which is the node IP.

* `dcs_cluster_redis_node` - The ID of the measurement object for the dimension **dcs_cluster_redis_node**.

* `status` - The measurement object status, which is the node status.

<a name="proxy_nodes_struct"></a>
The `proxy_nodes` block supports:

* `dcs_instance_id` - The measurement object ID, which is the node ID.

* `name` - The measurement object name, which is the node IP.

* `dcs_cluster_proxy_node` - The ID of the measurement object for the dimension **dcs_cluster_proxy_node**.

* `status` - The measurement object status, which is the node status.

<a name="proxy2_nodes_struct"></a>
The `proxy2_nodes` block supports:

* `dcs_instance_id` - The measurement object ID, which is the node ID.

* `name` - The measurement object name, which is the node IP.

* `dcs_cluster_proxy2_node` - The ID of the measurement object for the dimension **dcs_cluster_proxy_node**.

* `status` - The measurement object status, which is the node status.
