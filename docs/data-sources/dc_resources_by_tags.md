---
subcategory: "Direct Connect (DC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dc_resources_by_tags"
description: |-
  Use this data source to get the list of DC resources by tag.
---

# huaweicloud_dc_resources_by_tags

Use this data source to get the list of DC resources by tag.

## Example Usage

```hcl
variable "resource_type"{}

data "huaweicloud_dc_resources_by_tags" "test" {
  resource_type = var.resource_type
  action        = "filter"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `resource_type` - (Required, String) Specifies the direct Connect resource type. Value options:
  + **dc-directconnect**: direct connect connection
  + **dc-vgw**: virtual gateway
  + **dc-vif**: virtual interface

* `action` - (Required, String) Specifies the operate action. Value options:
  + **filter**: indicates pagination query
  + **count**: indicates that the total number of query results meeting the search criteria will be returned.

* `matches` - (Optional, List) Specifies the search criteria. The tag key is the parameter to match, for example,
  **resource_name**. The tag value indicates the value to be matched. This field is a fixed dictionary value. Determine
  whether fuzzy match is required based on different fields. For example, if key is **resource_name**, fuzzy search
  (case-insensitive) is used by default. If value is an empty string, exact match is used. If key is **resource_id**,
  exact match is used. Only **resource_name** for key is supported. Other key values will be available later.
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

* `sys_tags` - (Optional, List) Specifies the system tags. Only users with the op_service permission can use this parameter
  to filter resources. Only one tag structure is contained when this API is called by Tag Management Service (TMS). The
  key is **_sys_enterprise_project_id**, and the value is the enterprise project ID list. Currently, each key can contain
  only one value. 0 indicates the default enterprise project. sys_tags and tenant tag filtering conditions
  (`without_any_tag`, `tags`, `tags_any`, `not_tags`, and `not_tags_any`) cannot be used at the same time.
  The [sys_tags](#tags_struct) structure is documented below.

<a name="matches_struct"></a>
The `matches` block supports:

* `key` - (Required, String) Specifies the tag key.

* `value` - (Required, String) Specifies the tag value.

<a name="tags_struct"></a>
The `tags`, `not_tags`, `tags_any`, `not_tags_any` and `sys_tags` block supports:

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

* `sys_tags` - Indicates the list of queried system tags.
  The [sys_tags](#resources_tags_struct) structure is documented below.

<a name="resources_tags_struct"></a>
The `tags`, `sys_tags` block supports:

* `key` - Indicates the tag key.

* `values` - Indicates the tag value.
