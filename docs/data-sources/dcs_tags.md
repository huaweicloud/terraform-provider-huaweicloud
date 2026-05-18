---
subcategory: "Distributed Cache Service (DCS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dcs_tags"
description: |-
  Use this data source to query all resource tags of DCS instances within a specified project in HuaweiCloud.
---

# huaweicloud_dcs_tags

Use this data source to query all resource tags of DCS instances within a specified project in HuaweiCloud.

```hcl
data "huaweicloud_dcs_tags" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the tags. If omitted, the provider-level region
  will be used.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `tags` - The list of tags.
  The [tags](#tags_struct) structure is documented below.

<a name="tags_struct"></a>
The `tags` block supports:

* `key` - The tag key.

* `values` - The list of tag values.
