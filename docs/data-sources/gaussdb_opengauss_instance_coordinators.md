---
subcategory: "GaussDB"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_gaussdb_opengauss_instance_coordinators"
description: |-
  Use this data source to get the coordinate nodes of the GaussDB OpenGauss instance.
---

# huaweicloud_gaussdb_opengauss_instance_coordinators

Use this data source to get the coordinate nodes of the GaussDB OpenGauss instance.

## Example Usage

```hcl
variable "instance_id" {}

data "huaweicloud_gaussdb_opengauss_instance_coordinators" "test" {
  instance_id = var.instance_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the ID of the GaussDB OpenGauss instance.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `max_reduction_num` - Indicates the maximum number of nodes that can be deleted at a time.

* `nodes` - Indicates the node information list.

  The [nodes](#nodes_struct) structure is documented below.

<a name="nodes_struct"></a>
The `nodes` block supports:

* `id` - Indicates the node ID.

* `name` - Indicates the node name.

* `status` - Indicates the node status.
  The value can be:
  + **normal**: The node is normal.
  + **abnormal**: The node is abnormal.
  + **creating**: The node is being created.
  + **createfail**: The node fails to be created.

* `availability_zone` - Indicates the availability zone.

* `support_reduce` - Indicates whether the node can be deleted.
