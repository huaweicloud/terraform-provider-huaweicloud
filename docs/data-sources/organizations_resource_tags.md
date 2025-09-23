---
subcategory: "Organizations"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_organizations_resource_tags"
description: |-
  Use this data source to get the list of tags attached to the specified resource.
---

# huaweicloud_organizations_resource_tags

Use this data source to get the list of tags attached to the specified resource.

## Example Usage

```hcl
data "huaweicloud_organizations_resource_tags" "test"{
  resource_type = "organizations:ous"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `resource_type` - (Required, String) Specifies the resource type.
  Value options:
  + **organizations:policies**
  + **organizations:ous**
  + **organizations:accounts**
  + **organizations:roots**

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `tags` - Indicates the list of tags.

  The [tags](#tags_struct) structure is documented below.

<a name="tags_struct"></a>
The `tags` block supports:

* `key` - Indicates the key of the tag.

* `values` - Indicates the list of values of the tag.
