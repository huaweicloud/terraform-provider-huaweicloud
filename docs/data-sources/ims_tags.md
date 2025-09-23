---
subcategory: "Image Management Service (IMS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_ims_tags"
description: |-
  Use this data source to get the list of IMS tags within HuaweiCloud.
---

# huaweicloud_ims_tags

Use this data source to get the list of IMS tags within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_ims_tags" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID in UUID format.

* `tags` - The tag list.

  The [tags](#tags_struct) structure is documented below.

<a name="tags_struct"></a>
The `tags` block supports:

* `key` - The key of the tag.

* `values` - The value list of the tag. If the tag has only a key, it appears as an empty string in the value list.
