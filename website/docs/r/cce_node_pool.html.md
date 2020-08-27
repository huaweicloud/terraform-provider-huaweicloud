---
layout: "huaweicloud"
page_title: "huaweicloud: huaweicloud_cce_node_pool"
sidebar_current: "docs-huaweicloud-resource-cce-node-pool"
description: |-
  Add a node pool to a container cluster. 
---

# huaweicloud\_cce\_node\_pool
Add a node pool to a container cluster. 
This is an alternative to `huaweicloud_cce_node_pool`

## Example Usage

```hcl
variable "cluster_id" { }
variable "key_pair" { }
variable "availability_zone" { }

resource "huaweicloud_cce_node_pool" "node_pool" {
  cluster_id        = var.cluster_id
  name              = "testpool"
  os                = "EulerOS"
  initial_node_count = 2
  flavor_id         = "s3.large.4"
  availability_zone = var.availability_zone
  key_pair          = var.keypair
  region           =  "cn-east-3"
  scall_enable      = true
  min_node_count    = 1
  max_node_count    = 10
  scale_down_cooldown_time = 100
  priority         = 1

  root_volume {
    size       = 40
    volumetype = "SAS"
  }
  data_volumes {
    size       = 100
    volumetype = "SAS"
  }
}
``` 

## Argument Reference
The following arguments are supported:

* `cluster_id` - (Required) ID of the cluster. Changing this parameter will create a new resource.

* `name` - (Required) Node Pool Name.

* `initial_node_count` - (Required) Initial number of expected nodes in the node pool.

* `flavor_id` - (Required) Specifies the flavor id. Changing this parameter will create a new resource.
 
* `availability_zone` - (Required) specify the name of the available partition (AZ). Changing this parameter will create a new resource.

* `os` - (Optional) Operating System of the node. Changing this parameter will create a new resource.
    - For VM nodes, clusters of v1.13 and later support *EulerOS 2.5* and *CentOS 7.6*.
    - For BMS nodes purchased in the yearly/monthly billing mode, only *EulerOS 2.3* is supported.

* `key_pair` - (Optional) Key pair name when logging in to select the key pair mode. This parameter and `password` are alternative.
    Changing this parameter will create a new resource.

* `password` - (Optional) root password when logging in to select the password mode. This parameter must be salted and alternative to `key_pair`.
    Changing this parameter will create a new resource.

* `subnet_id` - (Optional) The ID of the subnet to which the NIC belongs. Changing this parameter will create a new resource.

* `bandwidth_charge_mode` - (Optional) Bandwidth billing type. Changing this parameter will create a new resource.

* `extend_param_charging_mode` - (Optional) Billing mode of a node. Value 0 indicates pay-per-use.
    Changing this parameter will create a new resource.

* `preinstall` - (Optional) Script required before installation. The input value can be a Base64 encoded string or not.
    Changing this parameter will create a new resource.

* `postinstall` - (Optional) Script required after the installation. The input value can be a Base64 encoded string or not.
    Changing this parameter will create a new resource.

* `scall_enable` - (Optional) Whether to enable auto scaling.

* `min_node_count` - (Optional) Minimum number of nodes allowed if auto scaling is enabled.

* `max_node_count` - (Optional) Maximum number of nodes allowed if auto scaling is enabled.

* `scale_down_cooldown_time` - (Optional) Interval between two scaling operations, in minutes.

* `priority` - (Optional) Weight of a node pool. A node pool with a higher weight has a higher priority during scaling.

**root_volume** **- (Required)** It corresponds to the system disk related configuration. Changing this parameter will create a new resource.

* `size` - (Required) Disk size in GB.
    
* `volumetype` - (Required) Disk type.
    
* `extend_param` - (Optional) Disk expansion parameters. 

**data_volumes** **- (Required)** Represents the data disk to be created. Changing this parameter will create a new resource.
    
* `size` - (Required) Disk size in GB.
    
* `volumetype` - (Required) Disk type.
    
* `extend_param` - (Optional) Disk expansion parameters. 

**taints** **- (Optional)** You can add taints to created nodes to configure anti-affinity. Each taint contains the following parameters:
    
* `key` - (Required) A key must contain 1 to 63 characters starting with a letter or digit. Only letters, digits, hyphens (-), 
  underscores (_), and periods (.) are allowed. A DNS subdomain name can be used as the prefix of a key.
    
* `value` - (Required) A value must start with a letter or digit and can contain a maximum of 63 characters, including letters, 
  digits, hyphens (-), underscores (_), and periods (.).
    
* `effect` - (Required) Available options are NoSchedule, PreferNoSchedule, and NoExecute. 
    
## Attributes Reference

All above argument parameters can be exported as attribute parameters along with attribute reference.

 * `status` -  Node status information.

 * `billing_mode` -  Billing mode of a node.