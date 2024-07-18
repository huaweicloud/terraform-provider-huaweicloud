---
subcategory: "GaussDB(for MySQL)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_gaussdb_mysql_proxy"
description: |-
  Manages GaussDB mysql proxy resource within HuaweiCloud.
---

# huaweicloud_gaussdb_mysql_proxy

Manages GaussDB mysql proxy resource within HuaweiCloud.

## Example Usage

```hcl
variable "instance_id" {}

resource "huaweicloud_gaussdb_mysql_proxy" "test" {
  instance_id = var.instance_id
  flavor      = "gaussdb.proxy.xlarge.x86.2"
  node_num    = 3
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) The region in which to create the GaussDB mysql proxy resource. If omitted,
  the provider-level region will be used. Changing this creates a new resource.

* `instance_id` - (Required, String, ForceNew) Specifies the ID of the GaussDB MySQL instance. Changing this parameter
  will create a new resource.

* `flavor` - (Required, String, ForceNew) Specifies the flavor of the proxy. Changing this parameter will create a new
  resource.

* `node_num` - (Required, Int) Specifies the node count of the proxy.

* `proxy_name` - (Optional, String) Specifies the name of the proxy. The name consists of **4** to **64** characters and
  starts with a letter. It is case-sensitive and can contain only letters, digits, hyphens (-), and underscores (_).

* `proxy_mode` - (Optional, String, ForceNew) Specifies the type of the proxy. Changing this creates a new resource.
  Value options:
  + **readwrite**: read and write.
  + **readonly**: read-only.

  Defaults to **readwrite**.

* `route_mode` - (Optional, Int, ForceNew) Specifies the routing policy of the proxy. Changing this creates a new
  resource. Value options:
  + **0**: weighted load balancing.
  + **1**: load balancing (The primary node does not process read requests).
  + **2**: load balancing (The primary node processes read requests).

* `subnet_id` - (Optional, String, ForceNew) Specifies the network ID of a subnet. Changing this creates a new resource.

* `master_node_weight` - (Optional, List) Specifies the read weight of the master node.
  The [master_node_weight](#node_weight_struct) structure is documented below.

* `readonly_nodes_weight` - (Optional, List) Specifies the read weight of the read-only node.
  The [readonly_nodes_weight](#node_weight_struct) structure is documented below.

* `new_node_auto_add_status` - (Optional, String) Specifies whether new nodes are automatically associate with proxy.
  Value options:
  + **ON**: New nodes are automatically associate with proxy.
  + **OFF**: New nodes are not automatically associate with proxy.

  -> **NOTE:** To configure this parameter, contact customer service.

* `new_node_weight` - (Optional, Int) Specifies the read weight of the new node.
  + If `route_mode` is **0** and `new_node_auto_add_status` is **ON**, the value of this parameter ranges from **0** to **1000**.
  + If `route_mode` is not **0** and `new_node_auto_add_status` is **OFF**, this parameter is unavailable.

* `port` - (Optional, Int) Specifies the port of the proxy.

<a name="node_weight_struct"></a>
The `master_node_weight` and `readonly_nodes_weight` block supports:

* `id` - (Required, String) Specifies the ID of the node.

* `weight` - (Required, Int) Specifies the weight assigned to the node.
  + If `route_mode` is **0**, the value is **0** to **1000**.
  + If `route_mode` is **1**, the value for the primary node is **0** and the value for read replicas is **0** or **1**.
  + If `route_mode` is **2**, the value for the primary node is **1** and the value for read replicas is **0** or **1**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - Indicates the resource ID.

* `address` - Indicates the address of the proxy.

* `nodes` - Indicates the node information of the proxy.
  The [nodes](#nodes_struct) structure is documented below.

<a name="nodes_struct"></a>
The `nodes` block supports:

* `id` - Indicates the proxy node ID.

* `status` - Indicates the proxy node status. The values can be:
  + **ACTIVE**: The node is available.
  + **ABNORMAL**: The node is abnormal.
  + **FAILED**: The node fails.
  + **DELETED**: The node has been deleted.

* `name` - Indicates the proxy node name.

* `role` - Indicates the proxy node role. The values can be:
  + **master**: primary node.
  + **slave**: read replica.

* `az_code` - Indicates the proxy node AZ.

* `frozen_flag` - Indicates whether the proxy node is frozen. The values can be:
  + **0**: unfrozen.
  + **1**: frozen.
  + **2**: deleted after being frozen.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 30 minutes.
* `update` - Default is 30 minutes.
* `delete` - Default is 10 minutes.

## Import

The GaussDB MySQL proxy can be imported using the `instance_id` and `id` separated by a slash, e.g.

```bash
$ terraform import huaweicloud_gaussdb_mysql_proxy.test <instance_id>/<id>
```

Note that the imported state may not be identical to your resource definition, due to the attribute missing from the
API response. The missing attribute is: `new_node_weight`, `proxy_mode`, `readonly_nodes_weight`. It is generally
recommended running `terraform plan` after importing a GaussDB MySQL proxy. You can then decide if changes should be
applied to the GaussDB MySQL proxy, or the resource definition should be updated to align with the GaussDB MySQL proxy.
Also you can ignore changes as below.

```hcl
resource "huaweicloud_gaussdb_mysql_proxy" "test" {
  ...

  lifecycle {
    ignore_changes = [
      new_node_weight, proxy_mode, readonly_nodes_weight,
    ]
  }
}
```
