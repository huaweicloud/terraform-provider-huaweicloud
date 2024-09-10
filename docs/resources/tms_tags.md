---
subcategory: "Tag Management Service (TMS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_tms_tags"
description: ""
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

* `tags` - (Required, List, ForceNew) Specifies an array of one or more predefined tags. The tags object
  structure is documented below. Changing this will create a new resource.

The `tags` block supports:

* `key` - (Required, String, ForceNew) Specifies the tag key. The value can contain up to `36` characters.
  Only English letters, Chinese characters, digits, hyphens (-) and underscores (_) are allowed.
  Changing this will create a new resource.

* `value` - (Required, String, ForceNew) Specifies the tag value. The value can contain up to `43` characters.
  Only English letters, Chinese characters, digits, periods (.), hyphens (-) and underscores (_) are allowed.
  Changing this will create a new resource.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 3 minutes.
* `delete` - Default is 3 minutes.
