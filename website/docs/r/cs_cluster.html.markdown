---
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cs_cluster"
sidebar_current: "docs-huaweicloud-resource-cs-cluster"
description: |-
  Cloud Stream Service cluster management
---

# huaweicloud\_cs\_cluster

Cloud Stream Service cluster management
This is an alternative to `huaweicloud_cs_cluster_v1`

## Example Usage

### create a cluster

```hcl
resource "huaweicloud_cs_cluster" "cluster" {
  name = "terraform_cs_cluster_test"
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
