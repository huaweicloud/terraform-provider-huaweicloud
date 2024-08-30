---
subcategory: "Cloud Container Engine (CCE)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cce_node"
description: ""
---

# huaweicloud_cce_node

To get the specified node in a cluster.

## Example Usage

```hcl
variable "cluster_id" {}
variable "node_name" {}

data "huaweicloud_cce_node" "node" {
  cluster_id = var.cluster_id
  name       = var.node_name
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to obtain the CCE nodes.
  If omitted, the provider-level region will be used.

* `cluster_id` - (Required, String) Specifies the ID of container cluster.

* `name` - (Optional, String) Specifies the name of the node.

* `node_id` - (Optional, String) Specifies the ID of the node.

* `status` - (Optional, String) Specifies the state of the node.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The node ID.

* `flavor_id` - The flavor ID to be used.

* `availability_zone` - Available partitions where the node is located.

* `os` - Operating System of the node.

* `subnet_id` - The ID of the subnet to which the NIC belongs.

* `ecs_group_id` - The ID of ECS group to which the node belongs.

* `tags` - Tags of a VM node, key/value pair format.

* `key_pair` - Key pair name when logging in to select the key pair mode.

* `billing_mode` - Node's billing mode: The value is 0 (on demand).

* `server_id` - The node's virtual machine ID in ECS.

* `public_ip` - Elastic IP parameters of the node.

* `private_ip` - Private IP of the node.

* `root_volume` - It corresponds to the system disk related configuration. Structure is documented below.

* `data_volumes` - Represents the data disk to be created. Structure is documented below.

* `enterprise_project_id` - The enterprise project ID of the node.

* `hostname_config` - The hostname config of the kubernetes node.
  The [object](#hostname_config) structure is documented below.

The `root_volume` and `data_volumes` blocks support:

* `size` - Disk size in GB.

* `volumetype` - Disk type.

* `extend_params` - Disk expansion parameters.

<a name="hostname_config"></a>
The `hostname_config` block supports:

* `type` - The hostname type of the kubernetes node.
