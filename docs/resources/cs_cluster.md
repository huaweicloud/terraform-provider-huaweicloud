---
subcategory: "Deprecated"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cs_cluster"
description: ""
---

# huaweicloud_cs_cluster

Cloud Stream Service cluster management.

!> **WARNING:** It has been deprecated, use `huaweicloud_dli_queue` instead.

## Example Usage

### create a cluster

```hcl
resource "huaweicloud_cs_cluster" "cluster" {
  name = "terraform_cs_cluster_test"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) The region in which to create the cloud stream service cluster resource. If
  omitted, the provider-level region will be used. Changing this creates a new cloud stream service cluster resource.

* `name` - (Required, String) Cluster name.

* `description` - (Optional, String) cluster description.

* `max_spu_num` - (Optional, Int) Cluster maximum SPU number.

* `subnet_cidr` - (Optional, String, ForceNew) Cluster sub segment. Changing this parameter will create a new resource.

* `subnet_gateway` - (Optional, String, ForceNew) Cluster subnet gateway. Changing this parameter will create a new
  resource.

* `vpc_cidr` - (Optional, String, ForceNew) Cluster VPC network segment. Changing this parameter will create a new
  resource.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - Specifies a resource ID in UUID format.

* `created_at` - Cluster creation time.

* `manager_node_spu_num` - Cluster management node SPU number.

* `used_spu_num` - The used SPU number of Cluster.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 30 minutes.
* `delete` - Default is 30 minutes.
