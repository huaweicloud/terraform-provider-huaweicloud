---
subcategory: "Data Security Center (DSC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dsc_tags"
description: |-
  Use this data source to get the metadata tag list within HuaweiCloud.
---

# huaweicloud_dsc_tags

Use this data source to get the metadata tag list within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_dsc_tags" "test" {
  category = "classification"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `category` - (Required, String) Specifies the tag category used to filter tags.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `tags` - The metadata tag list.

  The [tags](#tags_struct) structure is documented below.

<a name="tags_struct"></a>
The `tags` block supports:

* `id` - The tag ID.

* `name` - The tag name.

* `tag_gen_type` - The tag generation type.
