---
subcategory: "Cloud Search Service (CSS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_css_tags"
description: |-
  Use this data source to query the list of tags for a specified cluster.
---

# huaweicloud_css_tags

Use this data source to query the list of tags for a specified cluster.

## Example Usage

```hcl
variable "cluster_id" {}
variable "resource_type" {}

data "huaweicloud_css_tags" "test" {
  cluster_id    = var.cluster_id
  resource_type = var.resource_type
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `cluster_id` - (Required, String) Specifies the cluster ID.

* `resource_type` - (Required, String) Specifies the cluster resource type.
  The valid value only can be **css-cluster**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `tags` - The list of tags.
  The [tags](#tags_struct) structure is documented below.

<a name="tags_struct"></a>
The `tags` block supports:

* `key` - The tag key.

* `value` - The tag value.
