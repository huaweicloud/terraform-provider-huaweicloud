---
subcategory: Dedicated Load Balance (Dedicated ELB)
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_elb_loadbalancers_by_tags"
description: |-
  Use this data source to get the list of ELB load balancers filtered by tags within HuaweiCloud.
---


# huaweicloud_elb_loadbalancers_by_tags

Use this data source to get the list of ELB load balancers filtered by tags within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_elb_loadbalancers_by_tags" "test" {
  action = "filter"
  
  tags = [
    {
      key    = "key_string"
      values = ["value_string"]
    }
  ]

  matches = [
    {
      key   = "resource_name"
      value = "shared01"
    }
  ]
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

-> `tags` have limits as follows:
  <br/>1. This list cannot be an empty list.
  <br/>2. The list can contain up to `20` keys.
  <br/>3. Keys in this list must be unique.
  <br/>4. If no tag filtering condition is specified, full data is returned.

* `matches` - (Optional, List) Specifies the matches supported by resources. Keys in this list must be unique.
  Only two key is supported currently. Other key values will be available later.
  The [matches](#matches_struct) structure is documented below.

<a name="tags_struct"></a>
The `tags` block supports:

* `key` - (Required, String) Specifies the key of the resource tag. It contains a maximum of `127` unicode characters.
  A tag key cannot be an empty string. Spaces before and after a key will be deprecated.

* `values` - (Required, List) Specifies the list of values corresponding to the key.

  -> The field has the following restrictions:
    <br/>1. The list can contain up to `20` values.
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
  Currently, only **resource_name**, **resource_id** for key is supported.

* `value` - (Required, String) Specifies the value of the resource tag.
  A value consists of up to `255` characters.
  If key is **resource_name**, an empty string indicates exact match and any non-empty string indicates fuzzy match.
  If key is **resource_id**, indicates exact match.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `total_count` - The total number of matched resources.

* `resources` - List of matched resources.
  The [resources](#elb_loadbalancers_resources) structure is documented below.

<a name="elb_loadbalancers_resources"></a>
The `resources` block supports:

* `resource_id` - The resource ID.

* `super_resource_id` - The parent resource ID.

* `resource_name` - The resource name.

* `resource_detail` - The detail of the matched resources. The value is a resource object used for extension.
  This parameter is left blank by default.

* `tags` - The tag list.
  The [tags](#elb_loadbalancers_tags) structure is documented below.

<a name="elb_loadbalancers_tags"></a>
The `tags` block supports:

* `key` - The tag key.

* `value` - The tag value.
