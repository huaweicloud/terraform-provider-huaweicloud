---
subcategory: "VPC Endpoint (VPCEP)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_vpcep_resources"
description: |-
  Use this data source to get vpcep endpoints filtered by tags within HuaweiCloud.
---

# huaweicloud_vpcep_resources_by_tags

Use this data source to get vpcep resources filtered by tags within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_vpcep_resources_by_tags" "test" {
  action = "filter"
  resource_type = "endpoint_service"
  tags = [
    {
      key    = "key_string"
      values = ["value_string"]
    }
  ]
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `resource_type` - (Required, String) Specifies the resource type to which the tags belong that to be queried.  
  The value can be **endpoint_service** or **endpoint**.

* `action` - (Required, String) Specifies the action name. Possible values are **count** and **filter**.
  + **count**: querying count of data filtered by tags.
  + **filter**: querying details of data filtered by tags.

* `tags` - (Optional, List) Specifies the list of included tags. Backups with these tags will be filtered.
  The [tags](#tags_struct) structure is documented below.

* `tags_any` - (Optional, List) Specifies the list of tags. Backups with any tags in this list will be filtered.
  The [tags_any](#tags_struct) structure is documented below.

* `not_tags` - (Optional, List) Specifies the list of excluded tags. Backups without these tags will be filtered.
  The [not_tags](#tags_struct) structure is documented below.

* `not_tags_any` - (Optional, List) Specifies the list of tags. Backups without any tags in this list will be filtered.
  The [not_tags_any](#tags_struct) structure is documented below.

-> For arguments above, include `tags`, `tags_any`, `not_tags`, `not_tags_any` have limits as follows:
  <br/>1. This list cannot be an empty list.
  <br/>2. The list can contain up to `10` keys.
  <br/>3. Keys in this list must be unique.
  <br/>4. If no tag filtering condition is specified, full data is returned.

* `without_any_tag` - (Optional, Bool) Specifies whether ignore tags params.
  If this parameter is set to **true**, all resources without tags are queried.
  In this case, the `tag`, `not_tags`, `tags_any`, and `not_tags_any` fields are ignored.

* `sys_tags` - (Optional, List) Specifies the system tags.
  The [sys_tags](#tags_struct) structure is documented below.

  -> The sys_tags has limits as follows:
  <br/>1. Only users with the op_service permission can obtain this field.
  <br/>2. Field `sys_tags` and tag filter conditions (`tags`, `tags_any`, `not_tags`, `not_tags_any`)
  cannot  be used at the same time.
  <br/>3. If no `sys_tags` exists, use other tag APIs for filtering. If no tag filter is specified, full data is returned.
  <br/>4. This list cannot be an empty list.

* `matches` - (Optional, List) Specifies the matches supported by resources. Keys in this list must be unique.
  Only one key is supported currently. Multiple-key support will be available later.
  The [matches](#matches_struct) structure is documented below.

<a name="tags_struct"></a>
The `tags` block supports:

* `key` - (Required, String) Specifies the key of the resource tag. It contains a maximum of `127` unicode characters.
  A tag key cannot be an empty string. Spaces before and after a key will be deprecated.

* `values` - (Required, List) Specifies the list of values corresponding to the key.

  -> The field has the following restrictions:
    <br/>1. The list can contain up to `10` values.
    <br/>2. A tag value contains up to `255` unicode characters. Spaces before and after a key will be deprecated.
    <br/>3. Values in this list must be unique.
    <br/>4. Values in this list are in an OR relationship.
    <br/>5. This list can be empty and each value can be an empty character string.
    <br/>6. If this list is left blank, it indicates that all values are included.
    <br/>7. The asterisk (*) is a reserved character in the system.
    If the value starts with (*), it indicates that fuzzy match is performed based on the value following (*).
    The value cannot contain only asterisks.

<a name="matches_struct"></a>
The `matches` block supports:

* `key` - (Required, String) Specifies the key of the resource tag.
  A key can only be set to **resource_name**, indicating the resource name.

* `value` - (Required, String) Specifies the value of the resource tag.
  A value consists of up to `255` characters.
  If key is **resource_name**, an empty string indicates exact match and any non-empty string indicates fuzzy match.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `total_count` - The total number of matched resources.

* `resources` - List of matched resources.
  The [resources](#vpcep_resources) structure is documented below.

<a name="vpcep_resources"></a>
The `resources` block supports:

* `resource_id` - The resource ID.

* `resource_name` - The resource name.

* `resource_detail` - The detail of the matched resources.
  The value is a resource object used for extension. This parameter is left blank by default.

* `tags` - The tag list.
  The [tags](#vpcep_resources_tags) structure is documented below.

<a name="vpcep_resources_tags"></a>
The `tags` block supports:

* `key` - The tag key.

* `value` - The tag value.
