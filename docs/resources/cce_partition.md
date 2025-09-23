---
subcategory: "Cloud Container Engine (CCE)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cce_partition"
description: |-
  Manages a CCE partition resource within HuaweiCloud.
---

# huaweicloud_cce_partition

Manages a CCE partition resource within HuaweiCloud.

## Example Usage

```hcl
variable "cluster_id" {}

resource "huaweicloud_cce_partition" "test" {
  cluster_id = var.cluster_id
  category   = "IES"
  name       = "test_name"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the CCE add-on resource.
  If omitted, the provider-level region will be used. Changing this creates a new CCE add-on resource.

* `cluster_id` - (Optional, String, ForceNew) Specifies the cluster ID. Changing this parameter will create a new resource.

* `name` - (Optional, String, ForceNew) Specifies the name of the partition. Changing this parameter will create a new
  resource.

* `public_border_group` - (Optional, String, ForceNew) Specifies the group of the partition. Changing this parameter will
  create a new resource.

* `category` - (Optional, String, ForceNew) Specifies the category of the partition. Changing this parameter will create
  a new resource.

* `partition_subnet_id` - (Optional, String, ForceNew) Specifies the subnet ID of the partition. Changing this parameter
  will create a new resource.

* `container_subnet_ids` - (Optional, List) Specifies the container subnet IDs in the partition.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID. The value is the partition name.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 5 minutes.

## Import

The CCE partition can be imported using the `cluster_id` and `name` separated by a slash, e.g.:

```bash
$ terraform import huaweicloud_cce_partition.atest <cluster_id>/<name>
```
