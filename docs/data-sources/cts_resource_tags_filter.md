---
subcategory: "Cloud Trace Service (CTS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cts_resource_tags_filter"
description: |-
  Use this data source to filter CTS resources by tags within Huaweicloud.
---

# huaweicloud_cts_resource_tags_filter

Use this data source to filter CTS resources by tags within Huaweicloud.

## Example Usage

```hcl
data "huaweicloud_cts_resource_tags_filter" "test" {
  resource_type = "cts-tracker"
  
  tags {
    key    = "foo"
    values = ["bar", "bax"]
  }
  
  tags {
    key    = "test"
    values = ["alpha", "beta"]
  }
  
  matches {
    key   = "resource_name"
    value = "example"
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.  
  If omitted, the provider-level region will be used.

* `resource_type` - (Required, String) Specifies the resource type to be queried.  
  The valid value is **cts-tracker**.

* `tags` - (Optional, List) Specifies the tag list for filtering resources.  
  The [tags](#cts_filter_resource_tags_arg) structure is documented below.

* `matches` - (Optional, List) Specifies the match conditions for filtering resources.  
  The [matches](#cts_filter_resource_matches_arg) structure is documented below.  
  -> It matches exactly when `matches.value` is empty string. Otherwise, it matches fuzzily.

<a name="cts_filter_resource_tags_arg"></a>
The `tags` block supports:

* `key` - (Optional, String) Specifies the tag key.

* `values` - (Optional, List) Specifies the tag values.

<a name="cts_filter_resource_matches_arg"></a>
The `matches` block supports:

* `key` - (Optional, String) Specifies the match key.  
  The valid values is **resource_name**.

* `value` - (Optional, String) Specifies the match value.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `resources` - The list of resources that match the filter conditions.  
  The [resources](cts_filter_resources_attr) structure is documented below.

<a name="cts_filter_resources_attr"></a>
The `resources` block supports:

* `id` - The ID of the resource.

* `detail` - The detailed information of the resource.

* `name` - The name of the resource.

* `tags` - The tags associated with the resource.  
  The [tags](#cts_filter_resource_tags_attr) structure is documented below.

<a name="cts_filter_resource_tags_attr"></a>
The `tags` block supports:

* `key` - The tag key.

* `value` - The tag value.
