---
subcategory: "Cloud Operations Center (COC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_coc_ticket_operation_histories"
description: |-
  Use this data source to get the list of COC ticket operation histories.
---

# huaweicloud_coc_ticket_operation_histories

Use this data source to get the list of COC ticket operation histories.

## Example Usage

```hcl
variable "ticket_type" {}

data "huaweicloud_coc_ticket_operation_histories" "test" {
  ticket_type = var.ticket_type
}
```

## Argument Reference

The following arguments are supported:

* `ticket_type` - (Required, String) Specifies the type of work order to be queried.
  The valid values are as follows:
  + **incident**: Incident ticket.
  + **issues_mgmt**: Problem ticket.

* `string_filters` - (Optional, List) Specifies the string search condition, based on which you can search for specific
  work orders.

  The [string_filters](#filters_struct) structure is documented below.

* `sort_filter` - (Optional, List) Specifies the sorting criteria for the queried historical records.

  The [sort_filter](#filters_struct) structure is documented below.

<a name="filters_struct"></a>
The `string_filters` or `sort_filter` block supports:

* `operator` - (Optional, String) Specifies the expression operators.
  The value can be **in**, **like**, **desc**, **startwith**, **endwith**, **=**, **!=**, **>**, **<**, etc.

* `field` - (Optional, String) Specifies the name of the field that needs to be operated.
  + When querying tickets, the field name is **ticket_id**.
  + When sorting and filtering the query results by creation time, the field name is **start_time**.

* `name` - (Optional, String) Specifies the name of the condition that requires action.
  + When querying tickets, the field name is **ticket_id**.
  + When sorting and filtering the query results by creation time, the field name is **start_time**.

* `values` - (Optional, List) Specifies the conditional value that needs to be met when performing a certain operation
  on a field.
  + When querying work orders, the conditional value is the work order ID.
  + When sorting and filtering query results by creation time, the conditional value is **start_time**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `data` - Indicates the details of the operation record.

  The [data](#data_struct) structure is documented below.

<a name="data_struct"></a>
The `data` block supports:

* `action_id` - Indicates the operation identifier.

* `action` - Indicates the action information of the current operation.

* `sub_action` - Indicates the sub-action of the current operation.

* `operator` - Indicates the operator number of the current operation.

* `comment` - Indicates the comment information of the current operation.

* `id` - Indicates the primary key ID of the operation record.

* `ticket_id` - Indicates the work order number.

* `start_time` - Indicates the start time of the current operation.

* `stop_time` - Indicates the stop time of the current operation.

* `target_type` - Indicates the target type.

* `target_value` - Indicates the target value.

* `is_deleted` - Indicates whether the current operation record has been deleted.

* `update_time` - Indicates the update time.

* `action_name_zh` - Indicates the Chinese name of the current action.

* `action_name_en` - Indicates the English name of the current action.

* `action_template_zh` - Indicates the Chinese template of the current operation action.

* `action_template_en` - Indicates the English template of the current operation action.

* `status` - Indicates the work order status corresponding to the current operation.

* `final_sub_action` - Indicates the final sub-action.

* `enum_data_list` - Indicates the list of enumeration data.

  The [enum_data_list](#data_enum_data_list_struct) structure is documented below.

<a name="data_enum_data_list_struct"></a>
The `enum_data_list` block supports:

* `is_deleted` - Indicates whether the current enumeration value has been deleted.

* `match_type` - Indicates the match type.

* `ticket_id` - Indicates the order number.

* `real_ticket_id` - Indicates the real order number.

* `name_zh` - Indicates the Chinese name corresponding to the enumeration value.

* `name_en` - Indicates the English name corresponding to the enumeration value.

* `user_name` - Indicates the operator name.

* `biz_id` - Indicates the unique ID corresponding to the enumeration value.

* `prop_id` - Indicates the type corresponding to the current enumeration value.

* `model_id` - Indicates the model ID corresponding to different background applications.
