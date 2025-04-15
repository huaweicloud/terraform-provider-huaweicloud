---
subcategory: "Tag Management Service (TMS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_tms_tags"
description: |-
  Use this data source to get the list of predefined tags.
---

# huaweicloud_tms_tags

Use this data source to get the list of predefined tags.

## Example Usage

```hcl
data "huaweicloud_tms_tags" "test" {}
```

## Argument Reference

The following arguments are supported:

* `key` - (Optional, String) Specifies the tag key.
  Fuzzy search is supported. Key is case-insensitive. If the key contains non-URL-safe characters, it must be URL encoded.

* `value` - (Optional, String) Specifies the tag value.
  Fuzzy search is supported. Value is case-insensitive. If the value contains non-URL-safe characters, it must be URL encoded.

* `order_field` - (Optional, String) Specifies the sorting field:
  The field is case-sensitive, value options:
  + **update_time**: keys and values are sorted in ascending order.
  + **key**: values of `update_time` are sorted in descending order and `value` in ascending order.
  + **value**: values of `update_time` are sorted in descending order and `key` in ascending order.
  
  Defaults to **update_time**.

* `order_method` - (Optional, String) Specifies the sorting method of the `order_field` field.
  The method is case-sensitive and can be:
  + **asc**: ascending order
  + **desc**: descending order
  
  Defaults to **desc**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `tags` - Indicates the list of tags.

  The [tags](#tags_struct) structure is documented below.

<a name="tags_struct"></a>
The `tags` block supports:

* `key` - Indicates the key of the tag.

* `value` - Indicates the value of the tag.

* `update_time` - Indicates the time when the tag is updated.
