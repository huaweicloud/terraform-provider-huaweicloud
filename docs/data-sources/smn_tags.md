---
subcategory: "Simple Message Notification (SMN)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_smn_tags"
description: |-
  Use this data source to get the list of all tags of a resource type in a specified region.
---

# huaweicloud_smn_tags

Use this data source to get the list of all tags of a resource type in a specified region.

## Example Usage

```hcl
data "huaweicloud_smn_tags" "test" {
  resource_type = "smn_topic"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `resource_type` - (Required, String) Specifies the resource type.
  Value options:
  + **smn_topic**: topic
  + **smn_sms**: SMS
  + **smn_application**: mobile push

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `tags` - Indicates the resource tags.

  The [tags](#tags_struct) structure is documented below.

<a name="tags_struct"></a>
The `tags` block supports:

* `key` - Indicates the tag key.

* `values` - Indicates the tag values.
