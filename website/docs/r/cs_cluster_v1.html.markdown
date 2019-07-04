---
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cs_cluster_v1"
sidebar_current: "docs-huaweicloud-resource-cs-cluster-v1"
description: |-
  Cloud Stream Service Cluster Management
---

# huaweicloud\_cs\_cluster\_v1

Cloud Stream Service Cluster Management

## Example Usage

### create a cluster

```hcl
resource "huaweicloud_cs_cluster_v1" "cluster" {
  name = "terraform_cs_cluster_v1_test"
}
```

## Argument Reference

The following arguments are supported:

* `name` -
  (Required)
  Cluster name.

* `description` -
  (Optional)
  cluster description.

* `max_spu_num` -
  (Optional)
  Cluster maximum SPU number.

* `subnet_cidr` -
  (Optional)
  Cluster sub segment. Changing this parameter will create a new resource.

* `subnet_gateway` -
  (Optional)
  Cluster subnet gateway. Changing this parameter will create a new resource.

* `vpc_cidr` -
  (Optional)
  Cluster VPC network segment. Changing this parameter will create a new resource.

## Attributes Reference

In addition to the arguments listed above, the following computed attributes are exported:

* `created_at` -
  Cluster creation time.

* `manager_node_spu_num` -
  Cluster management node SPU number.

* `used_spu_num` -
  The used SPU number of Cluster.

## Timeouts

This resource provides the following timeouts configuration options:
- `create` - Default is 30 minute.
- `delete` - Default is 30 minute.
