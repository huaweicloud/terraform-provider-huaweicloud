---
subcategory: "Resource Access Manager (RAM)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_ram_tags"
description: |-
  Use this data source to get the list of tags that have been used in Resource Access Manager.
---

# huaweicloud_ram_tags

Use this data source to get the list of tags that have been used in Resource Access Manager.

## Example Usage

```hcl
data "huaweicloud_ram_tags" "test" {}
```

## Argument Reference

There are no arguments available for this data source.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `tags` - The list of tags.

  The [tags](#tags_struct) structure is documented below.

<a name="tags_struct"></a>
The `tags` block supports:

* `key` - The key of tags.

* `values` - All values of the key in the tags.
