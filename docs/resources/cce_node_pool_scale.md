---
subcategory: "Cloud Container Engine (CCE)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cce_node_pool_scale"
description: |-
  Use this resource to scale the CCE node pool within HuaweiCloud.
---

# huaweicloud_cce_node_pool_scale

Use this resource to scale the CCE node pool within HuaweiCloud.

## Example Usage

### Basic Usage

```hcl
variable "cluster_id" {}
variable "nodepool_id" {}

resource "huaweicloud_cce_node_pool_scale" "test" {
  cluster_id         = var.cluster_id
  nodepool_id        = var.nodepool_id
  scale_groups       = ["default"]
  desired_node_count = 2
}
```

~> Deleting node pool scale is not supported, it will only be removed from the state.

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the node pool scale resource.
  If omitted, the provider-level region will be used. Changing this will create a new resource.

* `cluster_id` - (Required, String, NonUpdatable) Specifies the cluster ID.

* `nodepool_id` - (Required, String, NonUpdatable) Specifies the node pool ID.

* `desired_node_count` - (Required, Int, NonUpdatable) Specifies the number of desired nodes.

* `scale_groups` - (Required, List, NonUpdatable) Specifies the names of scale groups to scale.
  **default** indicates the default group.

* `scalable_checking` - (Optional, String, NonUpdatable) Specifies the scalable checking.
  The value can be **instant** and **async**, defaults to **instant**.

* `charging_mode` - (Optional, String, ForceNew) Specifies the charging mode of the nodes.
  Valid values are **prePaid** and **postPaid**, defaults to **postPaid**.
  Changing this parameter will create a new cluster resource.

* `period_unit` - (Optional, String, ForceNew) Specifies the charging period unit of the nodes.
  Valid values are **month** and **year**. This parameter is mandatory if `charging_mode` is set to **prePaid**.
  Changing this parameter will create a new cluster resource.

* `period` - (Optional, Int, ForceNew) Specifies the charging period of the nodes.
  If `period_unit` is set to **month**, the value ranges from 1 to 9.
  If `period_unit` is set to **year**, the value ranges from 1 to 3.
  This parameter is mandatory if `charging_mode` is set to **prePaid**.
  Changing this parameter will create a new cluster resource.

* `auto_renew` - (Optional, String, ForceNew) Specifies whether auto renew is enabled. Valid values are **true** and **false**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID, which equals to `nodepool_id`.
