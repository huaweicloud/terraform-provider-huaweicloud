---
layout: "huaweicloud"
page_title: "Huaweicloud: huaweicloud_cce_cluster_v3 "
sidebar_current: "docs-huaweicloud-resource-cce-cluster-v3"
description: |-
  Provides Cloud Container Engine(CCE) resource.
---

# huaweicloud_cce_cluster_v3

Provides a cluster resource (CCE).

## Example Usage

 ```hcl
    variable "flavor_id" { }
    variable "vpc_id" { }
    variable "subnet_id" { }
	
    resource "huaweicloud_cce_cluster_v3" "cluster_1" {
     name = "cluster"
     cluster_type= "VirtualMachine"
     flavor_id= "${var.flavor_id}"
     vpc_id= "${var.vpc_id}"
     subnet_id= "${var.subnet_id}"
     container_network_type= "overlay_l2"
     description= "Create cluster"
    }
```

## Argument Reference

The following arguments are supported:


* `name` - (Required) Cluster name. Changing this parameter will create a new cluster resource.

* `labels` - (Optional) Cluster tag, key/value pair format. Changing this parameter will create a new cluster resource.

* `annotations` - (Optional) Cluster annotation, key/value pair format. Changing this parameter will create a new cluster resource.

* `flavor_id` - (Required) Cluster specifications. Changing this parameter will create a new cluster resource. Possible values:

	* `cce.s1.small` - small-scale single cluster (up to 50 nodes).
	* `cce.s1.medium` - medium-scale single cluster (up to 200 nodes).
	* `cce.s1.large` - large-scale single cluster (up to 1000 nodes).
	* `cce.s2.small` - small-scale HA cluster (up to 50 nodes).
	* `cce.s2.medium` - medium-scale HA cluster (up to 200 nodes).
	* `cce.s2.large` - large-scale HA cluster (up to 1000 nodes).
	* `cce.t1.small` - small-scale single physical machine cluster (up to 10 nodes).
	* `cce.t1.medium` - medium-scale single physical machine cluster (up to 100 nodes).
	* `cce.t1.large` - large-scale single physical machine cluster (up to 500 nodes).
	* `cce.t2.small` - small-scale HA physical machine cluster (up to 10 nodes).
	* `cce.t2.medium` - medium-scale HA physical machine cluster (up to 100 nodes).
	* `cce.t2.large` - large-scale HA physical machine cluster (up to 500 nodes).

* `cluster_version` - (Optional) For the cluster version, possible values are v1.7.3-r10 or v1.9.2-r1.

* `cluster_type` - (Required) Cluster Type, possible values are VirtualMachine and BareMetal. Changing this parameter will create a new cluster resource.

* `description` - (Optional) Cluster description.

* `billing_mode` - (Optional) Charging mode of the cluster, which is 0 (on demand). Changing this parameter will create a new cluster resource.

* `extend_param` - (Optional) Extended parameter. Changing this parameter will create a new cluster resource.

* `vpc_id` - (Required) The ID of the VPC used to create the node. Changing this parameter will create a new cluster resource.

* `subnet_id` - (Required) The ID of the subnet used to create the node. Changing this parameter will create a new cluster resource.

* `highway_subnet_id` - (Optional) The ID of the high speed network used to create bare metal nodes. Changing this parameter will create a new cluster resource.

* `container_network_type` - (Required) Container network parameters.

	* `overlay_l2` - An overlay_l2 network built for containers by using Open vSwitch(OVS)
	* `underlay_ipvlan` - An underlay_ipvlan network built for bare metal servers by using ipvlan.
	* `vpc-router` - An vpc-router network built for containers by using ipvlan and custom VPC routes.

* `container_network_cidr` - (Optional) Container network segment. Changing this parameter will create a new cluster resource.

## Attributes Reference

All above argument parameters can be exported as attribute parameters along with attribute reference.

  * `id` -  Id of the cluster resource.

  * `status` -  Cluster status information.

## Import

 Cluster can be imported using the cluster id, e.g.
 ```
 $ terraform import huaweicloud_cce_cluster_v3.cluster_1 4779ab1c-7c1a-44b1-a02e-93dfc361b32d  
```

