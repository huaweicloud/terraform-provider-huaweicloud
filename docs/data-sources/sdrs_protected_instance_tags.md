---
subcategory: "Storage Disaster Recovery Service (SDRS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_sdrs_protected_instance_tags"
description: |-
  Use this data source to query SDRS protected instance tags within HuaweiCloud.
---

# huaweicloud_sdrs_protected_instance_tags

Use this data source to query SDRS protected instance tags within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_sdrs_protected_instance_tags" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `tags` - The tag list.

  The [tags](#tags_struct) structure is documented below.

<a name="tags_struct"></a>
The `tags` block supports:

* `key` - The tag key.

* `values` - The tag values.
