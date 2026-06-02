---
subcategory: "Cloud Search Service (CSS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_css_cluster_volume_usage"
description: |-
  Use this data source to query the disk usage of a CSS cluster within HuaweiCloud.
---

# huaweicloud_css_cluster_volume_usage

Use this data source to query the disk usage of a CSS cluster within HuaweiCloud.

## Example Usage

### Basic Usage

```hcl
variable "cluster_id" {}

data "huaweicloud_css_cluster_volume_usage" "test" {
  cluster_id = var.cluster_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the cluster volume usage.
  If omitted, the provider-level region will be used.

* `cluster_id` - (Required, String) Specifies the ID of the CSS cluster.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `disk_info_list` - The disk usage of each node in the cluster.
  The [disk_info_list](#disk_info_list_struct) structure is documented below.

<a name="disk_info_list_struct"></a>
The `disk_info_list` block supports:

* `id` - The node ID.

* `name` - The node name.

* `group` - The node group.

* `role` - The node role.

* `disk_type` - The node disk type.

* `disk_capacity` - The total node storage capacity, in GB.

* `disk_used` - The used node storage capacity, in GB.

* `percentage` - The node storage usage, in percentage.
