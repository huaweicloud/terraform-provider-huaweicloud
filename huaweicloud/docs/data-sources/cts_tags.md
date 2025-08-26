---
subcategory: "Cloud Trace Service (CTS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cts_tags"
description: |-
  Use this data source to get the list of CTS tags within HuaweiCloud.
---

# huaweicloud_cts_tags

Use this data source to get the list of CTS tags within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_cts_tags" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.  
  If omitted, the provider-level region will be used.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `tags` - The list of tags.  
  The [tags](#cts_tags_attr) structure is documented below.

<a name="cts_tags_attr"></a>
The `tags` block supports:

* `key` - The tag key.

* `values` - The list of tag values.
