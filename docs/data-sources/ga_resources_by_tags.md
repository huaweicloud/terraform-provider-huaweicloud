---
subcategory: "Global Accelerator (GA)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_ga_resources_by_tags"
description: |- 
  Use this data source to get the list of GA resources by tags.
---

# huaweicloud_ga_resources_by_tags

Use this data source to get the list of GA resources by tags.

## Example Usage

```hcl
data "huaweicloud_ga_resources_by_tags" "test" {
  resource_type = "ga-accelerators"
  
  tags {
    key    = "foo"
    values = ["bar"]
  }
}
```

## Argument Reference

The following arguments are supported:

* `resource_type` - (Required, String) Specifies the resource type.
  Valid values include **ga-accelerators** and **ga-listeners**.

* `tags` - (Optional, List) Specifies the list of tags used for filtering resources.
  The [tags](#tags_struct) structure is documented below.

* `matches` - (Optional, List) Specifies the list of matches used for filtering resources.
  The [matches](#matches_struct) structure is documented below.

<a name="tags_struct"></a>
The `tags` block supports:

* `key` - (Required, String) Specifies the key of the tag.

* `values` - (Required, List) Specifies the list of the tag values.

<a name="matches_struct"></a>
The `matches` block supports:

* `key` - (Required, String) Specifies the key for matching a resource instance. The value can be: **resource_name**.

* `value` - (Required, String) Specifies the value for matching a resource instance.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `resources` - The list of target resources that matched filter parameters.
  The [resources](#resources_struct) structure is documented below.

* `total_count` - The total count of the resources.

<a name="resources_struct"></a>
The `resources` block supports:

* `resource_id` - The ID of the resource.

* `resource_name` - The name of the resource.

* `tags` - The list of tags associated with the resource.
  The [tags](#tags_struct) structure is documented below.

<a name="tags_struct"></a>
The `tags` block supports:

* `key` - The key of the tag.

* `value` - The value of the tag.
