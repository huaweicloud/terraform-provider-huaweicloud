---
subcategory: "Resource Access Manager (RAM)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_ram_resource_instances_count"
description: |-
  Use this data source to get the count of resource instance by tags.
---

# huaweicloud_ram_resource_instances_count

Use this data source to get the count of resource instance by tags.

## Example Usage

```hcl
data "huaweicloud_ram_resource_instances_count" "test" {}
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

* `total_count` - The total number of matched resource instances.
