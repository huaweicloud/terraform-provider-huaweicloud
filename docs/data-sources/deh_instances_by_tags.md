---
subcategory: "Dedicated Host (DeH)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_deh_instances_by_tags"
description: |-
  Use this data source to get the list of DeH instances by tag.
---

# huaweicloud_deh_instances_by_tags

Use this data source to get the list of DeH instances  by tag.

## Example Usage

```hcl
variable "action"{}

data "huaweicloud_deh_instances_by_tags" "test" {
  action = var.action
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `action` - (Required, String) Specifies the operate action. Value options:
  + **filter**: indicates pagination query
  + **count**: indicates that the total number of query results meeting the search criteria will be returned.

* `matches` - (Optional, List) Specifies the search criteria. Only **resource_name** for key is supported.
  The [matches](#matches_struct) structure is documented below.

* `tags` - (Optional, List) Specifies the tags. A maximum of 10 keys can be queried at a time, and each key can contain
  a maximum of 10 values. The structure body must be included. The tag key cannot be left blank or be an empty string.
  Each tag key must be unique, and each tag value of a tag must also be unique. Resources identified by different keys
  are in AND relationship, and values in one tag are in OR relationship. If no tag filtering criteria is specified, full
  data is returned.
  The [tags](#tags_struct) structure is documented below.

* `not_tags` - (Optional, List) Specifies the excluded tags. Each tag contains a maximum of 10 keys, and each key contains
  a maximum of 10 values. The structure body cannot be missing, and the key cannot be left blank or set to an empty string.
  Each tag key must be unique, and each tag value of a tag must also be unique. Resources not identified by different keys
  are in AND relationship, and values in one tag are in OR relationship. If not_tags_any is not specified, all resources
  will be returned.
  The [not_tags](#tags_struct) structure is documented below.

* `tags_any` - (Optional, List) Specifies any included tags. Each tag contains a maximum of 10 keys, and each key contains
  a maximum of 10 values. The structure body cannot be missing, and the key cannot be left blank or set to an empty string.
  Each tag key must be unique, and each tag value of a tag must also be unique. Resources identified by different keys are
  in OR relationship, and values in one tag are in OR relationship. If not_tags_any is not specified, all resources will
  be returned.
  The [tags_any](#tags_struct) structure is documented below.

* `not_tags_any` - (Optional, List) Specifies any excluded tags. Each tag contains a maximum of 10 keys, and each key
  contains a maximum of 10 values. The structure body cannot be missing, and the key cannot be left blank or set to an
  empty string. Each tag key must be unique, and each tag value of a tag must also be unique. Resources not identified
  by different keys are in OR relationship, and values in one tag are in OR relationship. If not_tags_any is not specified,
  all resources will be returned.
  The [not_tags_any](#tags_struct) structure is documented below.

<a name="matches_struct"></a>
The `matches` block supports:

* `key` - (Required, String) Specifies the tag key.

* `value` - (Required, String) Specifies the tag value.

<a name="tags_struct"></a>
The `tags`, `not_tags`, `tags_any` and `not_tags_any` block supports:

* `key` - (Required, String) Specifies the key of the tag.

* `values` - (Required, List) Specifies the values of the tag.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `resources` - Indicates the list of resources.
  The [resources](#gresources_struct) structure is documented below.

* `total_count` - Indicates the total number of resources.

<a name="gresources_struct"></a>
The `resources` block supports:

* `resource_id` - Indicates the resource ID.

* `resource_name` - Indicates the resource name.

* `resource_detail` - Indicates the provides details about the resource.

* `tags` - Indicates the list of queried tags.
  The [tags](#resources_tags_struct) structure is documented below.

<a name="resources_tags_struct"></a>
The `tags` block supports:

* `key` - Indicates the tag key.

* `values` - Indicates the tag value.
