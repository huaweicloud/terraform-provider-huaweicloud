---
subcategory: "Config"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_rms_tracked_resource_tags"
description: |-
  Use this data source to query tags of tracked resources.
---

# huaweicloud_rms_tracked_resource_tags

Use this data source to query tags of tracked resources.

## Example Usage

```hcl
data "huaweicloud_rms_tracked_resource_tags" "test" {}
```

## Argument Reference

The following arguments are supported:

* `key` - (Optional, String) Specifies the tag key name.

* `resource_deleted` - (Optional, Bool) Specifies whether to query deleted resources. Defaults to false.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `tags` - The tag list.
  The [tags](#tags_struct) structure is documented below.

<a name="tags_struct"></a>
The `tags` block supports:

* `key` - The tag key.

* `value` - The tag value list.
