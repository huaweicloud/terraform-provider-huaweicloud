---
subcategory: "Auto Scaling"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_as_group_tags"
description: |-
  Use this data source to get the list of all AS groups tags under the specified project.
---

# huaweicloud_as_group_tags

Use this data source to get the list of all AS groups tags under the specified project.

## Example Usage

```hcl
variable "resource_type" {}

data "huaweicloud_as_group_tags" "test" {
  resource_type = var.resource_type
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `resource_type` - (Required, String) Specifies the resource type.
  The valid values are as follows:
  + **scaling_group_tag**: Indicates the resource type is AS group.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `tags` - The list of the tags.

  The [tags](#tags_struct) structure is documented below.

<a name="tags_struct"></a>
The `tags` block supports:

* `key` - The key of the tag.

* `values` - The list of the tag values.
