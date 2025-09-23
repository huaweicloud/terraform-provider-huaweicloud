---
subcategory: "Tag Management Service (TMS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_tms_resource_tag_values"
description: |-
  Use this data source to get the list of tag values by tag key.
---

# huaweicloud_tms_resource_tag_values

Use this data source to get the list of tag values by tag key.

## Example Usage

```hcl
data "huaweicloud_tms_resource_tag_values" "test" {
  key = "tag_key"
}
```

## Argument Reference

The following arguments are supported:

* `key` - (Required, String) Specifies the tag key.

* `region_id` - (Optional, String) Specifies the region ID.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `values` - Indicates the tag values.
