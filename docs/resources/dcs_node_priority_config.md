---
subcategory: "Distributed Cache Service (DCS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dcs_node_priority_config"
description: |-
  Manages a DCS node priority config resource within HuaweiCloud.
---

# huaweicloud_dcs_node_priority_config

Manages a DCS node priority config resource within HuaweiCloud.

## Example Usage

### Basic Usage

```hcl
variable "instance_id" {}
variable "group_id" {}
variable "node_id" {}

resource "huaweicloud_dcs_node_priority_config" "test" {
  instance_id           = var.instance_id
  group_id              = var.group_id
  node_id               = var.node_id
  slave_priority_weight = 50
}
```

### Update Slave Priority Weight

```hcl
variable "instance_id" {}
variable "group_id" {}
variable "node_id" {}

resource "huaweicloud_dcs_node_priority_config" "test" {
  instance_id           = var.instance_id
  group_id              = var.group_id
  node_id               = var.node_id
  slave_priority_weight = 100
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.
  Changing this parameter will create a new resource.

* `instance_id` - (Required, String, NonUpdatable) Specifies the ID of the DCS instance.

* `group_id` - (Required, String, NonUpdatable) Specifies the ID of the instance shard group.

* `node_id` - (Required, String, NonUpdatable) Specifies the ID of the instance node.

* `slave_priority_weight` - (Required, Int) Specifies the slave node priority weight.
  The value ranges from **0** to **100**.
    + **0**: Default value, indicates that the failover is prohibited.
    + **1-100**: The priority decreases gradually from 1 to 100, where **1** is the highest priority and **100** is the
      lowest priority.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `logical_node_id` - The logical node ID.

* `name` - The node name.

* `status` - The node status.

* `az_code` - The availability zone code.

* `node_role` - The node role. The value can be:
    + **redis-server**: Redis server node.
    + **redis-proxy**: Proxy node.

* `node_type` - The node type. The value can be:
    + **master**: Master node.
    + **slave**: Slave node.
    + **proxy**: Proxy instance node.

* `node_ip` - The node IP address.

* `node_port` - The node port.

* `priority_weight` - The node priority weight for master/slave switchover.

* `is_access` - Whether the node IP is directly accessible.

* `group_name` - The instance shard group name.

* `is_remove_ip` - Whether the IP is removed from the read-only domain name.

* `replication_id` - The instance replication ID.

## Import

The DCS node priority config can be imported using the `instance_id`, `group_id`, and `node_id` separated by slashes,
e.g.

```bash
$ terraform import huaweicloud_dcs_node_priority_config.test <instance_id>/<group_id>/<node_id>
```
