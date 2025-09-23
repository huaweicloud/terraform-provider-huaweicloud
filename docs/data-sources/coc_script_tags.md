---
subcategory: "Cloud Operations Center (COC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_coc_script_tags"
description: |-
  Use this data source to get the list of COC script resource tags.
---

# huaweicloud_coc_script_tags

Use this data source to get the list of COC script resource tags.

## Example Usage

```hcl
data "huaweicloud_coc_script_tags" "test" {
  resource_type = "coc:script"
}
```

## Argument Reference

The following arguments are supported:

* `resource_type` - (Required, String) Specifies the resource type. Valid value is **coc:script**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `tags` - Indicates the list of tag.

  The [tags](#tags_struct) structure is documented below.

<a name="tags_struct"></a>
The `tags` block supports:

* `key` - Indicates the key of the tag.

* `values` - Indicates the values of the tag.
