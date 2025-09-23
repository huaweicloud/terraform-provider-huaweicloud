---
subcategory: "Cloud Operations Center (COC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_coc_issue_tickets"
description: |-
  Use this data source to get the list of COC tickets.
---

# huaweicloud_coc_issue_tickets

Use this data source to get the list of COC tickets.

## Example Usage

```hcl
variable "ticket_id" {}

data "huaweicloud_coc_issue_tickets" "test" {
  ticket_type = "issues_mgmt"
  string_filters {
    operator = "="
    field    = "ticket_id"
    values   = [var.ticket_id]
  }
}
```

## Argument Reference

The following arguments are supported:

* `ticket_type` - (Required, String) Specifies the ticket type to search for.
  + When querying the issue ticket list, the value **issues_mgmt**.

* `string_filters` - (Required, List) Specifies the string condition filter for querying tickets.

  The [string_filters](#filters_struct) structure is documented below.

* `sort_filter` - (Optional, List) Specifies the sort condition filter for querying tickets.

  The [sort_filter](#filters_struct) structure is documented below.

* `contain_sub_ticket` - (Optional, Bool) Specifies whether the current page contains sub-orders.

* `ticket_types` - (Optional, List) Specifies the ticket type to search for on the current page.

<a name="filters_struct"></a>
The `string_filters` or `sort_filter` block supports:

* `operator` - (Required, String) Specifies the operator.
  The values can be **in**, **like**, **startwith**, **endwith**, **=**, **!=**, **>** or **<**.

* `field` - (Required, String) Specifies the name of the field to be operated.
  + When querying work orders, the field name is **ticket_id**.
  + When sorting the query results by creation time, the condition value is result field name.

* `values` - (Required, List) Specifies the value to be filtered.
  + When querying work orders, the condition value is the work order ID.
  + When sorting the query results by creation time, the condition value is result field name.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `tickets` - Indicates the queried work order information.

  The [tickets](#data_tickets_struct) structure is documented below.

<a name="data_tickets_struct"></a>
The `tickets` block supports:

* `issue_correlation_sla` - Indicates the SLA to which the problem is associated.

* `level` - Indicates the problem ticket level.

* `root_cause_cloud_service` - Indicates the problem ticket responsibility service.

* `root_cause_type` - Indicates the classification of single root causes of problems.

* `current_cloud_service_id` - Indicates the problem ticket service.

* `issue_contact_person` - Indicates the contact person for the problem ticket.

* `issue_version` - Indicates the version number where the problem was found.

* `source` - Indicates the source of the issue.

* `commit_upload_attachment` - Indicates the ID of the uploaded attachment.

* `source_id` - Indicates the source ID of the problem ticket.

* `enterprise_project_id` - Indicates the issue ticket enterprise project ID.

* `virtual_schedule_type` - Indicates the problem ticket scheduling type.

* `title` - Indicates the issue title.

* `regions` - Indicates the region information of the problem ticket.

* `description` - Indicates the problem ticket description.

* `root_cause_comment` - Indicates the root cause analysis of a problem work order.

* `solution` - Indicates the problem ticket solution.

* `regions_search` - Indicates the value representing a ticket region search.

* `level_approve_config` - Indicates the ticket level approval configuration.

* `suspension_approve_config` - Indicates the configuration of pending approval for a problem ticket.

* `handle_time` - Indicates the time it takes to handle a problem ticket.

* `found_time` - Indicates the time when the problem ticket was discovered.

* `is_common_issue` - Indicates whether it is a common problem.

* `is_need_change` - Indicates whether the issue ticket requires changes.

* `is_enable_suspension` - Indicates whether the ticket is open or pending.

* `is_start_process_async` - Indicates whether to start an asynchronous process.

* `is_update_null` - Indicates whether to resubmit empty fields.

* `creator` - Indicates the creator of the problem ticket.

* `operator` - Indicates the operator of the problem ticket.

* `is_return_full_info` - Indicates whether to return all field information.

* `is_start_process` - Indicates whether to start the process.

* `ticket_id` - Indicates the problem ticket number.

* `real_ticket_id` - Indicates the problem ticket number.

* `assignee` - Indicates the current person responsible for the issue ticket.

* `participator` - Indicates the ticket participant.

* `work_flow_status` - Indicates the status of the problem ticket.

* `engine_error_msg` - Indicates the process status.

* `baseline_status` - Indicates the baseline state.

* `ticket_type` - Indicates the work order type.

* `phase` - Indicates the current stage information of the issue ticket.

* `sub_tickets` - Indicates the change sub-order information.

  The [sub_tickets](#tickets_sub_tickets_struct) structure is documented below.

* `enum_data_list` - Indicates the enumerated list representing issue associations.

  The [enum_data_list](#tickets_enum_data_list_struct) structure is documented below.

* `id` - Indicates the primary key UUID of the issue.

* `meta_data_version` - Indicates the application version that caused the problem.

* `update_time` - Indicates the update time.

* `create_time` - Indicates the creation time.

* `is_deleted` - Indicates whether the work order is deleted.

* `ticket_type_id` - Indicates the work order type.

* `form_info` - Indicates the action information.

<a name="tickets_sub_tickets_struct"></a>
The `sub_tickets` block supports:

* `change_ticket_id` - Indicates the associated change order number.

* `change_ticket_id_sub` - Indicates the change of sub-order number.

* `whether_to_change` - Indicates whether a change is required.

* `is_deleted` - Indicates whether it has been deleted.

* `id` - Indicates the change work order ID.

* `main_ticket_id` - Indicates the primary key ID of the change work order.

* `parent_ticket_id` - Indicates the parent work order ID.

* `ticket_id` - Indicates the problem ticket ID.

* `real_ticket_id` - Indicates the problem ticket number.

* `ticket_path` - Indicates the work order path.

* `target_value` - Indicates the region information.

* `target_type` - Indicates the sub-order type.

* `create_time` - Indicates the creation time.

* `update_time` - Indicates the update time.

* `creator` - Indicates the creator.

* `operator` - Indicates the operator.

<a name="tickets_enum_data_list_struct"></a>
The `enum_data_list` block supports:

* `is_deleted` - Indicates whether it has been deleted.

* `match_type` - Indicates the matching enumeration type.

* `ticket_id` - Indicates the current work order ID.

* `real_ticket_id` - Indicates the work order number.

* `name_zh` - Indicates the Chinese name.

* `name_en` - Indicates the English name.

* `biz_id` - Indicates the unique ID corresponding to the enumeration value.

* `prop_id` - Indicates the type corresponding to the current enumeration value.

* `model_id` - Indicates the model ID corresponding to different background applications.
