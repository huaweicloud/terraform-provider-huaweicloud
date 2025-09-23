---
subcategory: "Global Accelerator (GA)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_ga_tags"
description: |-
  Use this data source to get the list of tags.
---

# huaweicloud_ga_tags

Use this data source to get the list of tags.

## Example Usage

```hcl
variable "resource_type" {}

data "huaweicloud_ga_tags" "test" {
  resource_type = var.resource_type
}
```

## Argument Reference

The following arguments are supported:

* `resource_type` - (Required, String) Specifies the resource type.
  The valid values are as follows:
  + **ga-accelerators**
  + **ga-listeners**

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `tags` - The list of tags.

  The [tags](#tags_struct) structure is documented below.

<a name="tags_struct"></a>
The `tags` block supports:

* `key` - The key of the tag.

* `values` - The list of the tag values.
