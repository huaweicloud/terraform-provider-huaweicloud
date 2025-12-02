---
subcategory: "Dedicated Host (DeH)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_deh_instance_tags"
description: |-
  Use this data source to get the list of GeH instance tags.
---

# huaweicloud_deh_instance_tags

Use this data source to get the list of GeH instance tags.

## Example Usage

```hcl
data "huaweicloud_deh_instance_tags" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `tags` - The list of all tags for resources of the DeH instance.

  The [tags](#tags_struct) structure is documented below.

<a name="tags_struct"></a>
The `tags` block supports:

* `key` - The key of the resource tag.

* `values` - All values corresponding to the key.
