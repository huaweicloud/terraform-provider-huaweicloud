---
subcategory: "Cloud Search Service (CSS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_css_cluster_resource_tags"
description: |-
  Use this data source to get the all tags of a specific cluster.
---

# huaweicloud_css_cluster_resource_tags

Use this data source to get the all tags of a specific cluster.

## Example Usage

```hcl
varivale "cluster_id" {}

data "huaweicloud_css_cluster_resource_tags" "test" {
  resource_type = "css-cluster"
  cluster_id    = var.cluster_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `resource_type` - (Required, String) Specifies the resource type.
  Currently, its value can only be **css-cluster**.

* `cluster_id` - (Required, String) Specifies the ID of the cluster.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `tags` - The list of cluster tags.
