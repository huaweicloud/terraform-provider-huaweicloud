---
subcategory: "TaurusDB"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_taurusdb_tags"
description: |-
  Use this data source to query all tags of TaurusDB instances in a specified project.
---

# huaweicloud_taurusdb_tags

Use this data source to query all tags of TaurusDB instances in a specified project.

## Example Usage

```hcl
data "huaweicloud_taurusdb_tags" "test" {
  engine_name = "taurus"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `engine_name` - (Optional, String) Specifies the engine name. The value defaults to **taurus**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `tags` - The list of tags.

  The [tags](#tags_struct) structure is documented below.

<a name="tags_struct"></a>
The `tags` block supports:

* `key` - Indicates the tag key.

* `values` - Indicates the tag values.
