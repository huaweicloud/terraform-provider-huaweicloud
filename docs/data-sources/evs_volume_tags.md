---
subcategory: "Elastic Volume Service (EVS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_evs_volume_tags"
description: |-
  Use this data source to query the list of tags of the EVS volume within HuaweiCloud.
---

# huaweicloud_evs_volume_tags

Use this data source to query the list of tags of the EVS volume within HuaweiCloud.

## Example Usage

```hcl
variable "volume_id" {}

data "huaweicloud_evs_volume_tags" "test" {
  volume_id = var.volume_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `volume_id` - (Required, String) Specifies the EVS volume ID.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `tags` - The tags list.
  The [tags](#tags_structure) structure is documented below.

<a name="tags_structure"></a>
The `tags` block supports:

* `key` - The key of the tag.

* `value` - The value of the tag.
