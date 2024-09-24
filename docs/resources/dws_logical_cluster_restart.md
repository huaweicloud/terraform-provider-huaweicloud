---
subcategory: "GaussDB(DWS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dws_logical_cluster_restart"
description: |-
  Use this resource to restart the logical cluster within HuaweiCloud.
---

# huaweicloud_dws_logical_cluster_restart

Use this resource to restart the logical cluster within HuaweiCloud.

-> This resource is only a one-time action resource for restarting the logical cluster. Deleting this resource will
   not clear the corresponding request record, but will only remove the resource information from the tfstate file.

## Example Usage

```hcl
variable "dws_cluster_id" {}
variable "logical_cluster_id" {}

resource "huaweicloud_dws_logical_cluster_restart" "test" {
  cluster_id         = var.dws_cluster_id
  logical_cluster_id = var.logical_cluster_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.
  Changing this creates a new resource.

* `cluster_id` - (Required, String, ForceNew) Specifies the DWS cluster ID to which the logical cluster belongs.
  Changing this creates a new resource.

* `logical_cluster_id` - (Required, String, ForceNew) Specifies the ID of the logical cluster to be restarted.
  Changing this creates a new resource.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 10 minutes.
