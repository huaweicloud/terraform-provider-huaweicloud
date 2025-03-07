---
subcategory: "GaussDB"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_gaussdb_opengauss_predefined_tags"
description: |-
  Use this data source to get the list of predefined tags.
---

# huaweicloud_gaussdb_opengauss_predefined_tags

Use this data source to get the list of predefined tags.

## Example Usage

```hcl
data "huaweicloud_gaussdb_opengauss_predefined_tags" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `tags` - The tags.

  The [tags](#tags_struct) structure is documented below.

<a name="tags_struct"></a>
The `tags` block supports:

* `key` - Indicates the tag key.

* `values` - Indicates the tag values.
