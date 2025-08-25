---
subcategory: "rms"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_rms_tags"
description: |-
  Use this data source to get a list of RMS tags.
---

# huaweicloud_rms_tags

Use this data source to get a list of RMS tags.

## Example Usage

```hcl
data "huaweicloud_rms_tags" "test" {}
```

## Argument Reference

The following arguments are supported:

* `key` - (Optional, String) Tag key name.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `tags` - Tag list.
  The [tags](#tags_struct) structure is documented below.

<a name="tags_struct"></a>The `tags` block supports:

* `key` - Tag key.

* `value` - Tag value list.
