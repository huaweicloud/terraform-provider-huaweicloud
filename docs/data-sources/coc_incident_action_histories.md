---
subcategory: "Cloud Operations Center (COC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_coc_incident_action_histories"
description: |-
  Use this data source to get the list of COC incident action histories.
---

# huaweicloud_coc_incident_action_histories

Use this data source to get the list of COC incident action histories.

## Example Usage

```hcl
variable "incident_id" {}

data "huaweicloud_coc_incident_action_histories" "test" {
  incident_id = var.incident_id
}
```

## Argument Reference

The following arguments are supported:

* `incident_id` - (Required, String) Specifies the event ticket number.

* `count_filters` - (Optional, List) Specifies the counting condition filter for querying event tickets.

  The [count_filters](#filters_struct) structure is documented below.

* `fields` - (Optional, List) Specifies the fields returned by the specified interface.

* `group_by_filter` - (Optional, List) Specifies the group filter for querying event tickets.

  The [group_by_filter](#filters_struct) structure is documented below.

* `int_filters` - (Optional, List) Specifies the int condition filter for querying event tickets.

  The [int_filters](#filters_struct) structure is documented below.

* `string_filters` - (Optional, List) Specifies the string condition filter for querying event tickets.

  The [string_filters](#filters_struct) structure is documented below.

* `sort_filter` - (Optional, List) Specifies the sort condition filter for querying event tickets.

  The [sort_filter](#filters_struct) structure is documented below.

* `condition` - (Optional, String) Specifies the expression to assemble complex expressions.
  The value can be parentheses **()**, and **&**, or **|**.
  For example: ( filterName1 & filterName2) | filterName3. The **filterName** taken from `string_filters.name`.
  If left blank, the conditions in string_filters default to an AND relationship.

<a name="count_filters_struct"></a>
The `count_filters` block supports:

* `name` - (Optional, String) Specifies the filter name.

* `filters` - (Optional, List) Specifies the list of filters.

  The [filters](#filters_struct) structure is documented below.

<a name="filters_struct"></a>
The `string_filters`, `sort_filter`, `filters`, `group_by_filter`, and `int_filters` block supports:

* `field` - (Required, String) Specifies the field key. For values, refer to the return field name.

* `operator` - (Required, String) Specifies the operator.
  The values can be **in**, **like**, **startwith**, **endwith**, **=**, **!=**, **>** or **<**.

* `values` - (Required, List) Specifies the value to be filtered.
  The value content supports String, timestamp, enumeration, and integer, all of which are passed in the form of strings.

* `name` - (Optional, String) Specifies the filter name for condition expression splicing.

* `group` - (Optional, String) Specifies the group to query information.

* `match_type` - (Optional, String) Specifies the matching method.

* `priority_type` - (Optional, String) Specifies the priority processing method.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `data` - Indicates the list of incident action historical information.

  The [data](#data_info_struct) structure is documented below.

<a name="data_info_struct"></a>
The `data` block supports:

* `action` - Indicates the operation type.

* `action_name_zh` - Indicates the Chinese operation type.

* `action_name_en` - Indicates the English operation type.

* `operator` - Indicates the operator ID.

* `status` - Indicates the current status.

* `start_time` - Indicates the start time of the operation.

* `stop_time` - Indicates the end time of the operation.

* `comment` - Indicates the comment information.

* `enum_data_list` - Indicates the enumeration data list.

  The [enum_data_list](#info_enum_data_list_struct) structure is documented below.

<a name="info_enum_data_list_struct"></a>
The `enum_data_list` block supports:

* `prop_id` - Indicates the field identified by the field key.

* `biz_id` - Indicates the enumeration key.

* `name_zh` - Indicates the Chinese name.

* `name_en` - Indicates the English name.
