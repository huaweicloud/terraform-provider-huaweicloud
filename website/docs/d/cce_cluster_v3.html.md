---
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cce_cluster_v3  "
sidebar_current: "docs-huaweicloud-datasource-cce-cluster-v3"
description: |-
  Get information on Cloud Container Engine Cluster (CCE).
---

# huaweicloud_cce_cluster_v3

   Provides details about all clusters and obtains certificate for accessing cluster information.

## Example Usage

 ```hcl
  variable "cluster_name" { }
  variable "cluster_id" { }
  variable "vpc_id" { }

  data "huaweicloud_cce_cluster_v3" "cluster" {
   name = "${var.cluster_name}"
   id= "${var.cluster_id}"
   status= "Available"
  }
```

## Argument Reference

The following arguments are supported:

* `name` -  (Optional)The Name of the cluster resource.
 
* `id` - (Optional) The ID of container cluster.

* `status` - (Optional) The state of the cluster.

* `cluster_type` - (Optional) Type of the cluster. Possible values: VirtualMachine, BareMetal or Windows

## Attributes Reference

All above argument parameters can be exported as attribute parameters along with attribute reference:

* `billingMode` - Charging mode of the cluster.

* `description` - Cluster description.

* `name` - The name of the cluster in string format.

* `id` - The ID of the cluster.
  
* `flavor` - The cluster specification in string format.

* `cluster_version` - The version of cluster in string format.

* `container_network_cidr` - The container network segment.

* `container_network_type` - The container network type: overlay_l2 , underlay_ipvlan or vpc-router.
  
* `subnet_id` - The ID of the subnet used to create the node.

* `highway_subnet_id` - The ID of the high speed network used to create bare metal nodes.

**endpoints**

* `internal` - The address accessed within the user's subnet.

* `external` - Public network access address.






 


