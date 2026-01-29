---
subcategory: "Resource Access Manager (RAM)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_ram_resource_instances_filter"
description: |-
  Use this data source to filter the resource instances by tags.
---

# huaweicloud_ram_resource_instances_filter

Use this data source to filter the resource instances by tags.

## Example Usage

```hcl
data "huaweicloud_ram_resource_instances_filter" "test" {
  without_any_tag = true
}
```

## Argument Reference

The following arguments are supported:

* `without_any_tag` - (Optional, Bool) Specifies the flag to query instances without tags.
  When this flag is set to **true**, it queries all resources without tags.

* `tags` - (Optional, List) Specifies the list of tags.

  The [tags](#tags_struct) structure is documented below.

* `matches` - (Optional, List) Specifies the name of RAM permission in which to query the data source.

  The [matches](#matches_struct) structure is documented below.

<a name="tags_struct"></a>
The `tags` block supports:

* `key` - (Required, String) Specifies the key of tags.

* `values` - (Optional, List) Specifies all values of the key in the tags.

<a name="matches_struct"></a>
The `matches` block supports:

* `key` - (Required, String) Specifies the key of matched tags.

* `value` - (Required, String) Specifies the value of the key in the matched tags.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID in UUID format.

* `resources` - The list of resource information.

  The [resources](#resources_struct) structure is documented below.

* `total_count` - The total number.

<a name="resources_struct"></a>
The `resources` block supports:

* `resource_id` - Indicates the ID of resource.

* `resource_name` - Indicates the name of resource.

* `tags` - Indicates the list of tags.

  The [tags](#resources_tags_struct) structure is documented below.

* `resource_detail` - Indicates the details of resource.

<a name="resources_tags_struct"></a>
The `tags` block supports:

* `key` - Indicates the key of tags.

* `value` - Indicates the value of the key in the tags.
