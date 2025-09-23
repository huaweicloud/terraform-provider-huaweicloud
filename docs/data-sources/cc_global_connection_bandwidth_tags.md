---
subcategory: "Cloud Connect (CC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cc_global_connection_bandwidth_tags"
description: |-
  Use this data source to get all global connection bandwidth tags in an account.
---

# huaweicloud_cc_global_connection_bandwidth_tags

Use this data source to get all global connection bandwidth tags in an account.

## Example Usage

```hcl
data "huaweicloud_cc_global_connection_bandwidth_tags" "test"{}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `tags` - Indicates the tag list.

  The [tags](#tags_struct) structure is documented below.

<a name="tags_struct"></a>
The `tags` block supports:

* `key` - Indicates the tag key.

* `values` - Indicates the list of tag values.
