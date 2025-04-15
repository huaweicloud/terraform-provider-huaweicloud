---
subcategory: "Tag Management Service (TMS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_tms_resource_tag_keys"
description: |-
  Use this data source to get the list of tag keys.
---

# huaweicloud_tms_resource_tag_keys

Use this data source to get the list of tag keys.

## Example Usage

```hcl
data "huaweicloud_tms_resource_tag_keys" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region_id` - (Optional, String) Specifies the region ID.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `keys` - Indicates the tag keys.
