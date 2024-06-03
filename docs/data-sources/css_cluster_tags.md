---
subcategory: "Cloud Search Service (CSS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_css_cluster_tags"
description: |-
  Use this data source to get the list of CSS cluster tags.
---

# huaweicloud_css_cluster_tags

Use this data source to get the list of CSS cluster tags.

## Example Usage

```hcl
data "huaweicloud_css_cluster_tags" "test" {
  resource_type = "css-cluster"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `resource_type` - (Required, String) Specifies the resource type.
  Currently, its value can only be **css-cluster**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `tags` - The list of cluster tags.
