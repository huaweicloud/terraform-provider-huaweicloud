---
subcategory: "Resource Center"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_rc_resource_tags"
description: |-
  Use this data source to get the resource tags within HuaweiCloud Resource Center.
---

# huaweicloud_rc_resource_tags

Use this data source to get the resource tags within HuaweiCloud Resource Center.

## Example Usage

```hcl
data "huaweicloud_rc_resource_tags" "test" {}
```

## Arguments Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `key` - (Optional, String) Specifies the tag key name.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID in UUID format.

* `tags` - The tag list.
  The [tags](#tags_struct) structure is documented below.

<a name="tags_struct"></a>
The `tags` block supports:

* `key` - The tag key.

* `value` - The tag value list.
