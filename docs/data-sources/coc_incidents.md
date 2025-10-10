---
subcategory: "Cloud Operations Center (COC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_coc_incidents"
description: |-
  Use this data source to get the list of COC incidents.
---

# huaweicloud_coc_incidents

Use this data source to get the list of COC incidents.

## Example Usage

```hcl
data "huaweicloud_coc_incidents" "test" {}
```

## Argument Reference

The following arguments are supported:

* `contain_sub_ticket` - (Optional, Bool) Specifies whether to include sub-order information in the query results.

* `string_filters` - (Optional, List) Specifies the string condition filter for querying event tickets.

  The [string_filters](#filters_struct) structure is documented below.

* `sort_filter` - (Optional, List) Specifies the sort condition filter for querying event tickets.

  The [sort_filter](#filters_struct) structure is documented below.

* `condition` - (Optional, String) Specifies the expression to assemble complex expressions.
  The value can be parentheses **()**, and **&**, or **|**. The default value is **&**.
  For example: ( filterName1 & filterName2) | filterName3. The **filterName** taken from ``string_filters.name``.

* `count_filters` - (Optional, List) Specifies the counting condition filter for querying event tickets.

  The [count_filters](#count_filters_struct) structure is documented below.

* `fields` - (Optional, List) Specifies the fields returned by the specified interface.

* `group_by_filter` - (Optional, List) Specifies the group filter for querying event tickets.

  The [group_by_filter](#filters_struct) structure is documented below.

* `int_filters` - (Optional, List) Specifies the int condition filter for querying event tickets.

  The [int_filters](#filters_struct) structure is documented below.

* `ticket_types` - (Optional, List) Specifies the work order type to be queried.
  When obtaining an event list, the value passed here is **incident**.

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

* `tickets` - Indicates the list of result data.

  The [tickets](#data_tickets_struct) structure is documented below.

<a name="data_tickets_struct"></a>
The `tickets` block supports:

* `current_cloud_service_id` - Indicates the cloud service ID.

* `level_id` - Indicates the event level.
  [For details](https://support.huaweicloud.com/api-coc/coc_api_04_03_001_006.html#coc_api_04_03_001_006__section289718103710)

* `mtm_region` - Indicates the region.

* `source_id` - Indicates the event source.
  [For details](https://support.huaweicloud.com/api-coc/coc_api_04_03_001_006.html#coc_api_04_03_001_006__section10172124616391)

* `forward_rule_id` - Indicates the forwarding rules.

* `enterprise_project_id` - Indicates the enterprise project ID.

* `mtm_type` - Indicates the event category.
  [For details](https://support.huaweicloud.com/api-coc/coc_api_04_03_001_006.html#coc_api_04_03_001_006__section137541363918)

* `title` - Indicates the event title.

* `description` - Indicates the event description.

* `ticket_id` - Indicates the event ticket number.

* `is_service_interrupt` - Indicates whether the service is interrupted.

* `work_flow_status` - Indicates the process status.

* `phase` - Indicates the process stage.

* `assignee` - Indicates the assignee.

* `creator` - Indicates the creator.

* `operator` - Indicates the last operator.

* `update_time` - Indicates the update time.

* `create_time` - Indicates the creation time.

* `start_time` - Indicates the fault start time.

* `handle_time` - Indicates the processing time.

* `incident_ownership` - Indicates the incident ownership.

* `enum_data_list` - Indicates the enumeration list.

  The [enum_data_list](#tickets_enum_data_list_struct) structure is documented below.

<a name="tickets_enum_data_list_struct"></a>
The `enum_data_list` block supports:

* `prop_id` - Indicates the field identified by the field key.

* `biz_id` - Indicates the enumeration key.

* `name_zh` - Indicates the Chinese name.

* `name_en` - Indicates the English name.
