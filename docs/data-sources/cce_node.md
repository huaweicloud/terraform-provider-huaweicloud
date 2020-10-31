---
subcategory: "Cloud Container Engine (CCE)"
---

# huaweicloud\_cce\_node

To get the specified node in a cluster.
This is an alternative to `huaweicloud_cce_node_v3`

## Example Usage

```hcl
variable "cluster_id" { }
variable "node_name" { }

data "huaweicloud_cce_node" "node" {
  cluster_id = var.cluster_id
  name       = var.node_name
}
```
## Argument Reference

The following arguments are supported:

* `region` - (Optional) The region in which to obtain the cce nodes. If omitted, the provider-level region will work as default.
 
* `Cluster_id` - (Required) The id of container cluster.

* `name` - (Optional) - Name of the node.

* `node_id` - (Optional) - The id of the node.

* `status` - (Optional) - The state of the node.


## Attributes Reference

All above argument parameters can be exported as attribute parameters along with attribute reference:

* `flavor_id` - The flavor id to be used. 

* `availability_zone` - Available partitions where the node is located. 

* `key_pair` - Key pair name when logging in to select the key pair mode.

* `billing_mode` - Node's billing mode: The value is 0 (on demand).

* `charge_mode` - Bandwidth billing type.

* `bandwidth_size` - Bandwidth (Mbit/s), in the range of [1, 2000].

* `extendparam` - 	Extended parameters. 
    
* `node_count` - The number of nodes in batch creation.

* `eip_ids` - List of existing elastic IP IDs.
 
* `server_id` - The node's virtual machine ID in ECS.

* `public_ip` - Elastic IP parameters of the node.

* `private_ip` - Private IP of the node

* `ip_type` - Elastic IP address type.

* `share_type` - The bandwidth sharing type.

NOTE:
This parameter is mandatory when share_type is set to PER and is optional when share_type is set to WHOLE with an ID specified.

Enumerated values: PER (indicates exclusive bandwidth) and WHOLE (indicates sharing)


**root_volumes**

* `disk_size` - Disk size in GB.

* `volumetype` - Disk type.

**data_volumes**

* `disk_size` - Disk size in GB.

* `volumetype` - Disk type.
