---
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cs_peering_connect_v1"
sidebar_current: "docs-huaweicloud-resource-cs-peering-connect-v1"
description: |-
  Cloud Stream Service cluster peering connect management
---

# huaweicloud\_cs\_peering\_connect\_v1

Cloud Stream Service cluster peering connect management

## Example Usage

### create a cluster peering connect

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
```

## Argument Reference

The following arguments are supported:

* `cluster_id` -
  (Required)
  The id of cloud stream cluster. Changing this parameter will create a new resource.

* `name` -
  (Required)
  The name of peering connection. Changing this parameter will create a new resource.

* `target_vpc_info` -
  (Optional)
  The information of target vpc. Structure is documented below. Changing this parameter will create a new resource.

The `target_vpc_info` block supports:

* `project_id` -
  (Optional)
  The project ID to which target vpc belongs. Changing this parameter will create a new resource.

* `vpc_id` -
  (Required)
  The VPC ID. Changing this parameter will create a new resource.

## Timeouts

This resource provides the following timeouts configuration options:
- `create` - Default is 30 minute.
