---
subcategory: "Cloud Trace Service (CTS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cts_resource_tags"
description: |-
  Use this data source to get resource tag list of CTS service within HuaweiCloud.
---

# huaweicloud_cts_resource_tags

Use this data source to get resource tag list of CTS service within HuaweiCloud.

## Example Usage

```hcl
variable "tracker_id" {}

data "huaweicloud_cts_resource_tags" "test" {
  resource_id = var.tracker_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource tags.  
  If omitted, the provider-level region will be used.

* `resource_type` - (Required, String) Specifies the resource type to be queried.  
  The valid value is **cts-tracker**.

* `resource_id` - (Required, String) Specifies the resource ID to be queried.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `tags` - The list of tags that matched filter parameters.  
  The [tags](#cts_resource_tags_attr) structure is documented below.

<a name="cts_resource_tags_attr"></a>
The `tags` block supports:

* `key` - The tag key.

* `value` - The tag value.
