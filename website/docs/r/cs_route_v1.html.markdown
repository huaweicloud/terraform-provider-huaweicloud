---
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cs_route_v1"
sidebar_current: "docs-huaweicloud-resource-cs-route-v1"
description: |-
  Cloud Stream Service cluster peering connect route management
---

# huaweicloud\_cs\_route\_v1

Cloud Stream Service cluster peering connect route management

## Example Usage

### create a cluster peering connect route

```hcl
resource "huaweicloud_cs_cluster_v1" "cluster" {
  name = "terraform_cs_cluster_v1_test"
}

resource "huaweicloud_vpc_v1" "vpc" {
  name = "terraform_vpc_v1_test"
  cidr = "192.168.0.0/16"
}

resource "huaweicloud_vpc_subnet_v1" "subnet" {
  name = "terraform_vpc_subnet_v1_test"
  cidr = "192.168.0.0/16"
  gateway_ip = "192.168.0.1"
  vpc_id = "${huaweicloud_vpc_v1.vpc.id}"
}

resource "huaweicloud_cs_peering_connect_v1" "peering" {
  name = "terraform_cs_peering_connect_v1_test"
  target_vpc_info {
    vpc_id = "${huaweicloud_vpc_v1.vpc.id}"
  }
  cluster_id = "${huaweicloud_cs_cluster_v1.cluster.id}"
}

resource "huaweicloud_cs_route_v1" "route" {
  cluster_id = "${huaweicloud_cs_cluster_v1.cluster.id}"
  peering_id = "${huaweicloud_cs_peering_connect_v1.peering.id}"
  destination = "${huaweicloud_vpc_subnet_v1.subnet.cidr}"
}
```

## Argument Reference

The following arguments are supported:

* `cluster_id` -
  (Required)
  The id of cloud stream cluster. Changing this parameter will create a new resource.

* `destination` -
  (Required)
  Routing destination CIDR. Changing this parameter will create a new resource.

* `peering_id` -
  (Required)
  The peering connection id of cloud stream cluster. Changing this parameter will create a new resource.
