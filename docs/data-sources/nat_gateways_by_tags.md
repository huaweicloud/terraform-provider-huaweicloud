---
subcategory: "NAT Gateway (NAT)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_nat_gateways_by_tags"
description: |-
  Use this data source to get NAT Gateways filtered by tags within HuaweiCloud.
---

# huaweicloud_nat_gateways_by_tags

Use this data source to get NAT Gateways filtered by tags within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_nat_gateways_by_tags" "test" {
  action = "filter"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

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
  The [resources](#nat_gateways_resources) structure is documented below.

<a name="nat_gateways_resources"></a>
The `resources` block supports:

* `resource_id` - The resource ID.

* `resource_name` - The resource name.

* `resource_detail` - The detail of the matched resources. The value is a resource object used for extension.
  This parameter is left blank by default.

* `tags` - The tag list.
  The [tags](#nat_gateways_tags) structure is documented below.

<a name="nat_gateways_tags"></a>
The `tags` block supports:

* `key` - The key of the resource tag.

  -> The key of tags has limits as follows:
    <br/>1. It can contain a maximum of `36` characters.
    <br/>2. It cannot be an empty string.
    <br/>3. Spaces before and after a key will be discarded.
    <br/>4. It cannot contain non-printable ASCII characters (`0`–`31`) and the following characters: =*<>,|/
    <br/>5. It can contain only letters, digits, hyphens (-), and underscores (_).

* `value` - The value of the resource tag.

  -> The value of tags has limits as follows:
    <br/>1. It is mandatory when a tag is added and optional when a tag is deleted.
    <br/>2. It can contain a maximum of `43` characters.
    <br/>3. It can be an empty string.
    <br/>4. Spaces before and after a value will be discarded.
    <br/>5. It cannot contain non-printable ASCII characters (`0`–`31`) and the following characters: =*<>,|/
    <br/>6. It can contain only letters, digits, hyphens (-), underscores (_), and periods (.).
