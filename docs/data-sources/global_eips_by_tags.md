---
subcategory: "Elastic IP (EIP)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_global_eips_by_tags"
description: |-
  Use this data source to get the list of global EIPs filtered by tags.
---

# huaweicloud_global_eips_by_tags

Use this data source to get the list of global EIPs filtered by tags.

## Example Usage

```hcl
data "huaweicloud_global_eips_by_tags" "test" {
  tags {
    key   = "foo"
    value = "bar"
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `tags` - (Required, List) Specifies the tag filter conditions. It can contain up to `20` tags.

  The [tags](#query_tags) structure is documented below.

<a name="query_tags"></a>
The `tags` block supports:

* `key` - (Required, String) Specifies the tag key.

* `value` - (Optional, String) Specifies the tag value.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `resources` - The list of resources.

  The [resources](#resources_struct) structure is documented below.

<a name="resources_struct"></a>
The `resources` block supports:

* `resource_id` - The resource ID.

* `resource_detail` - The resource detail object.

* `tags` - The tag list associated with the resource.

  The [tags](#tags_struct) structure is documented below.

* `resource_name` - The resource name.

<a name="tags_struct"></a>
The `tags` block supports:

* `key` - The key of the tag.

* `value` - The value of the tag.
