---
subcategory: "Elastic IP (EIP)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_global_eip_tags"
description: |-
  Use this data source to get a list of global EIP tags.
---

# huaweicloud_global_eip_tags

Use this data source to get a list of global EIP tags.

## Example Usage

```hcl
data "huaweicloud_global_eip_tags" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `total_count` - Indicates the total count of tags.

* `tags` - Indicates the list of tags.

  The [tags](#tags_struct) structure is documented below.

<a name="tags_struct"></a>
The `tags` block supports:

* `key` - Indicates the tag key.

* `values` - Indicates the list of tag values.
