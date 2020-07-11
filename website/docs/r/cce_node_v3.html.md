---
layout: "huaweicloud"
page_title: "huaweicloud: huaweicloud_cce_node_v3"
sidebar_current: "docs-huaweicloud-resource-cce-node-v3"
description: |-
  Add a node to a container cluster. 
---


# huaweicloud_cce_node_v3
Add a node to a container cluster. 

## Example Usage

```hcl
variable "cluster_id" { }
variable "ssh_key" { }
variable "availability_zone" { }

resource "huaweicloud_cce_node_v3" "node_1" {
  cluster_id        = var.cluster_id
  availability_zone = var.availability_zone
  name              = "test"
  flavor_id         = "s1.medium"
  key_pair          = var.ssh_key

  root_volume {
    size       = 40
    volumetype = "SATA"
  }
  data_volumes {
    size       = 100
    volumetype = "SATA"
  }

  iptype         = "5_bgp"
  sharetype      = "PER"
  bandwidth_size = 100
}
``` 

## Argument Reference
The following arguments are supported:

* `cluster_id` - (Required) ID of the cluster. Changing this parameter will create a new resource.

* `name` - (Optional) Node Name.

* `flavor_id` - (Required) Specifies the flavor id. Changing this parameter will create a new resource.
 
* `availability_zone` - (Required) specify the name of the available partition (AZ). Changing this parameter will create a new resource.

* `os` - (Optional) Operating System of the node, possible values are EulerOS 2.2 and CentOS 7.1. Defaults to EulerOS 2.2.
    Changing this parameter will create a new resource.

* `key_pair` - (Optional) Key pair name when logging in to select the key pair mode. This parameter and `password` are alternative.
    Changing this parameter will create a new resource.

* `password` - (Optional) root password when logging in to select the password mode. This parameter must be salted and alternative to `key_pair`.
    Changing this parameter will create a new resource.

* `subnet_id` - (Optional) The ID of the subnet to which the NIC belongs. Changing this parameter will create a new resource.

* `eip_ids` - (Optional) List of existing elastic IP IDs. Changing this parameter will create a new resource.

**Note:**
If the eip_ids parameter is configured, you do not need to configure the bandwidth parameters: iptype, bandwidth_charge_mode, bandwidth_size and share_type.

* `iptype` - (Optional) Elastic IP type. Default is 5_bgp. Changing this parameter will create a new resource.

* `bandwidth_charge_mode` - (Optional) Bandwidth billing type. Changing this parameter will create a new resource.

* `sharetype` - (Optional) Bandwidth sharing type. Changing this parameter will create a new resource.

* `bandwidth_size` - (Optional) Bandwidth size. Changing this parameter will create a new resource.


* `max_pods` - (Optional) The maximum number of instances a node is allowed to create. Changing this parameter will create a new cluster resource.

* `preinstall` - (Optional) Script required before installation. The input value can be a Base64 encoded string or not.
    Changing this parameter will create a new resource.

* `postinstall` - (Optional) Script required after installation. The input value can be a Base64 encoded string or not.
   Changing this parameter will create a new resource.

**root_volume** **- (Required)** It corresponds to the system disk related configuration. Changing this parameter will create a new resource.

* `size` - (Required) Disk size in GB.
    
* `volumetype` - (Required) Disk type.
    
* `extend_param` - (Optional) Disk expansion parameters. 

**data_volumes** **- (Required)** Represents the data disk to be created. Changing this parameter will create a new resource.
    
* `size` - (Required) Disk size in GB.
    
* `volumetype` - (Required) Disk type.
    
* `extend_param` - (Optional) Disk expansion parameters. 
    
## Attributes Reference

All above argument parameters can be exported as attribute parameters along with attribute reference.

 * `status` -  Node status information.

 * `server_id` - ID of the ECS instance associated with the node.

 * `private_ip` - Private IP of the CCE node.

 * `public_ip` - Public IP of the CCE node.
