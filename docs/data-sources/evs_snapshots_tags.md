---
subcategory: "Elastic Volume Service (EVS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_evs_snapshots_tags"
description: |-
  Use this data source to query the list of tags of all EVS snapshots within HuaweiCloud.
---

# huaweicloud_evs_snapshots_tags

Use this data source to query the list of tags of all EVS snapshots within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_evs_snapshots_tags" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `tags` - The tags list.
  The [tags](#tags_structure) structure is documented below.

<a name="tags_structure"></a>
The `tags` block supports:

* `key` - The key of the tag.

* `value` - The value of the tag.
