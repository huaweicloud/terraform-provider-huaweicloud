---
layout: "huaweicloud"
page_title: "huaweicloud: huaweicloud_cce_nodes_v3"
sidebar_current: "docs-huaweicloud-resource-cce-nodes-v3"
description: |-
  Add a node to a container cluster. 
---


# huaweicloud_cce_nodes_v3
Add a node to a container cluster. 

## Example Usage

 ```hcl
   variable "cluster_id" { }
   variable "ssh_key" { }

   resource "huaweicloud_cce_node_v3" "node_1" {
    cluster_id="${var.cluster_id}"
    name = "node1"
    flavor="s1.medium"
    iptype="5_bgp"
    az= "cn-east-2a"
    sshkey="${var.ssh_key}"
    root_volume = {
     size= 40,
     volumetype= "SATA"
    }
    sharetype= "PER"
    bandwidth_size= 100,
    data_volumes = [
     {
      size= 100,
      volumetype= "SATA"
     },
    ]
  }
 ```    


## Argument Reference
The following arguments are supported:

* `cluster_id` - (Required) ID of the cluster. Changing this parameter will create a new resource.

* `billing_mode` - (Optional) Node's billing mode: The value is 0 (on demand). Changing this parameter will create a new resource.

* `node_count` - (Required) The number of nodes in batch creation.

* `name` - (Optional) Node Name.

* `labels` - (Optional) Node tag, key/value pair format. Changing this parameter will create a new resource.

* `annotations` - (Optional) Node annotation, key/value pair format. Changing this parameter will create a new resource.
    
* `flavor` - (Required) Specifies the flavor id. Changing this parameter will create a new resource.
    
* `az` - (Required) specify the name of the available partition (AZ). Changing this parameter will create a new resource.

* `sshkey` - (Required) Key pair name when logging in to select the key pair mode. Changing this parameter will create a new resource.

* `eip_ids` - (Optional) List of existing elastic IP IDs. Changing this parameter will create a new resource.

**Note:**
If the eip_ids parameter is configured, you do not need to configure the eip_count and bandwidth parameters: iptype, charge_mode, bandwidth_size and share_type.

* `eip_count` - (Optional) Number of elastic IPs to be dynamically created. Changing this parameter will create a new resource.

* `iptype` - (Required) Elastic IP type. 

* `bandwidth_charge_mode` - (Optional) Bandwidth billing type. Changing this parameter will create a new resource.

* `sharetype` - (Required) Bandwidth sharing type. Changing this parameter will create a new resource.

* `bandwidth_size` - (Required) Bandwidth size. Changing this parameter will create a new resource.

* `extend_param_charging_mode` - (Optional) Node charging mode, 0 is on-demand charging. Changing this parameter will create a new cluster resource.

* `ecs_performance_type` - (Optional) Classification of cloud server specifications. Changing this parameter will create a new cluster resource.

* `order_id` - (Optional) Order ID, mandatory when the node payment type is the automatic payment package period type. Changing this parameter will create a new cluster resource.

* `product_id` - (Optional) The Product ID. Changing this parameter will create a new cluster resource.

* `max_pods` - (Optional) The maximum number of instances a node is allowed to create. Changing this parameter will create a new cluster resource.

* `public_key` - (Optional) The Public key. Changing this parameter will create a new cluster resource.

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
  
