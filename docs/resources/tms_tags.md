---
subcategory: "Tag Management Service (TMS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_tms_tags"
description: |-
  Manages TMS tags resource within HuaweiCloud.
---

# huaweicloud_tms_tags

Manages TMS tags resource within HuaweiCloud.

## Example Usage

```hcl
resource "huaweicloud_tms_tags" "test" {
  tags {
    key   = "foo"
    value = "bar"
  }
}
```

## Argument Reference

The following arguments are supported:

* `tags` - (Required, List) Specifies an array of one or more predefined tags.
  The [tags](#tags_struct) structure is documented below.

<a name="tags_struct"></a>
The `tags` block supports:

* `key` - (Required, String) Specifies the tag key. The value can contain up to `36` characters. Only English letters,
  Chinese characters, digits, hyphens (-) and underscores (_) are allowed.

* `value` - (Required, String) Specifies the tag value. The value can contain up to `43` characters. Only English letters,
  Chinese characters, digits, periods (.), hyphens (-) and underscores (_) are allowed.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.
