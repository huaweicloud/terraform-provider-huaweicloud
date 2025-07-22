---
subcategory: "Host Security Service (HSS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_hss_tags"
description: |-
  Use this data source to get the list of used tags.
---

# huaweicloud_hss_tags

Use this data source to get the list of used tags.

## Example Usage

```hcl
variable "resource_type" {}

data "huaweicloud_hss_tags" "test" {
  resource_type = var.resource_type
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `resource_type` - (Required, String) Specifies the resource type.
  The value only can be **hss**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `tags` - The list of tags.
  The [tags](#tags_struct) structure is documented below.

<a name="tags_struct"></a>
The `tags` block supports:

* `key` - The tags key

* `values` - The tags value list.
