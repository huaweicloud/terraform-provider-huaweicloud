---
subcategory: "Distributed Message Service (DMS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dms_rocketmq_instance_nodes"
description: |-
  Use this data source to get the list of RocketMQ nodes within HuaweiCloud.
---

# huaweicloud_dms_rocketmq_instance_nodes

Use this data source to get the list of RocketMQ nodes within HuaweiCloud.

## Example Usage

```hcl
variable "instance_id" {}

data "huaweicloud_dms_rocketmq_instance_nodes" "test" {
  instance_id = var.instance_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.  
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the instance ID.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `nodes` - The list of nodes.  
  The [nodes](#nodes_struct) structure is documented below.

<a name="nodes_struct"></a>
The `nodes` block supports:

* `id` - The node ID.

* `broker_name` - The broker name.

* `broker_id` - The broker ID.

* `address` - The private address.

* `public_address` - The public address.
