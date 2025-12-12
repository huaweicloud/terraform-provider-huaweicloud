---
subcategory: "Resource Access Manager (RAM)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_ram_resource_instance_filter"
description: |
  Use this data source to filter the resource instance by tags.
---

# huaweicloud_ram_resource_instances_filter

Use this data source to filter the resource instance by tags.

## Example Usage

```hcl
variable "without_any_tag" {}
variable "key" {}
variable "value" {}
variable "match_key" {}
variable "match_value" {}

data "huaweicloud_ram_resource_instances_filter" "test" {
  without_any_tag = var.without_any_tag
  tags {
    key   = var.key
    values = [var.value] 
  }
  matches {
    key   = var.match_key
    value = [var.match_value]
  }
}
```

## Argument Reference

The following arguments are supported:

* `without_any_tag` - (Optional, Bool) Specifies the flag to query instances without tags.
  When this flag is set to true, it queries all resources without tags.

* `tags` - (Optional, List) Specifies the list of tags.

  The [tags](#tags) structure is documented below.

* `matches` - (Optional, List) Specifies the name of RAM permission in which to query the data source.

  The [matches](#matches) structure is documented below.

<a name="tags"></a>
The `tags` block supports:

* `key` - (Required, String) Specifies the key of tags.

* `values` - (Optional, List) Specifies all values of the key in the tags.

<a name="matches"></a>
The `matches` block supports:

* `key` - (Required, String) Specifies the key of matched tags.

* `value` - (Required, String) Specifies all values of the key in the matched tags.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - Indicates the resource ID.

* `resources` - Indicates the list of the RAM permissions
  The [resources](#RAM_Permissions) structure is documented below.

* `total_count` - The total number of the RAM permissions.

<a name="RAM_Permissions"></a>
The `resources` block supports:

* `resource_id` - Indicates the ID of resource.

* `resource_name` - Indicates the name of resource.

* `resource_detail` - Indicates the details of resource.

* `tags` - Indicates the list of tags.

  The [tags](#tags) structure is documented below.

<a name="tags"></a>
The `tags` block supports:

* `key` - Indicates the key of tags.

* `value` - Indicates all values of the key in the tags.
