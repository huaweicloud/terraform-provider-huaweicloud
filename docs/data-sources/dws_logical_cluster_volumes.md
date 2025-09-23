---
subcategory: "GaussDB(DWS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dws_logical_cluster_volumes"
description: |-
  Use this data source to get the list of disk volumes corresponding to the logical cluster within HuaweiCloud.
---

# huaweicloud_dws_logical_cluster_volumes

Use this data source to get the list of disk volumes corresponding to the logical cluster within HuaweiCloud.

## Example Usage

```hcl
variable "dws_cluster_id" {}

data "huaweicloud_dws_logical_cluster_volumes" "test" {
  cluster_id = var.dws_cluster_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `cluster_id` - (Required, String) Specified the ID of the cluster to which the logical cluster belongs.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `volumes` - The list of the disk volumes corresponding to the logical cluster.

  The [volumes](#logical_cluster_volumes_struct) structure is documented below.

<a name="logical_cluster_volumes_struct"></a>
The `volumes` block supports:

* `logical_cluster_name` - The name of the logical cluster.

* `percentage` - The percentage of disk space used.

* `usage` - The used capacity of the disk.

* `total` - The total capacity of the disk.
