---
subcategory: "Relational Database Service (RDS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_rds_publications"
description: |-
  Use this data source to get the list of RDS instance publications.
---

# huaweicloud_rds_publications

Use this data source to get the list of RDS instance publications.

## Example Usage

```hcl
variable "instance_id" {}

data "huaweicloud_rds_publications" "test" {
  instance_id = var.instance_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the ID of the RDS instance.

* `publication_name` - (Optional, String) Specifies the publication name. Fuzzy search is supported.

* `publication_db_name` - (Optional, String) Specifies the publication database name. Fuzzy search is supported.

* `subscriber_instance_id` - (Optional, String) Specifies the ID of the subscriber instance.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `publications` - Indicates the list of publications.
  The [publications](#publications_struct) structure is documented below.

<a name="publications_struct"></a>
The `publications` block supports:

* `id` - Indicates the ID of publication.

* `status` - Indicates the status of publication. The value can be: **normal**, **abnormal**, **creating**, **modifying**,
  **createfail**.

* `publication_name` - Indicates the name of publication.

* `publication_database` - Indicates the name of publication database.

* `subscription_count` - Indicates the count of subscription.

* `subscription_options` - Indicates the options of subscription.
  The [subscription_options](#subscription_options_struct) structure is documented below.

* `job_schedule` - Indicates the schedule info.
  The [job_schedule](#job_schedule_struct) structure is documented below.

* `is_select_all_table` - Indicates whether select all table.

* `extend_tables` - Indicates the extent tables when select all table.

* `tables` - Indicates the tables of publication.
  The [tables](#tables_struct) structure is documented below.

<a name="subscription_options_struct"></a>
The `subscription_options` block supports:

* `independent_agent` - Indicates whether is independent agent.

* `snapshot_always_available` - Indicates whether the snapshot is always available.

* `replicate_ddl` - Indicates whether the replicate ddl can be changed.

* `allow_initialize_from_backup` - Indicates whether allow to initialize from backup.

<a name="job_schedule_struct"></a>
The `job_schedule` block supports:

* `id` - Indicates the schedule ID.

* `job_schedule_type` - Indicates the schedule type.

* `one_time_occurrence` - Indicates the occurrence time for one time.
  The [one_time_occurrence](#one_time_occurrence_struct) structure is documented below.

* `frequency` - Indicates the frequency of the schedule.
  The [frequency](#frequency_struct) structure is documented below.

* `daily_frequency` - Indicates the daily frequency of the schedule.
  The [daily_frequency](#daily_frequency_struct) structure is documented below.

* `duration` - Indicates the duration of the schedule.
  The [duration](#duration_struct) structure is documented below.

<a name="one_time_occurrence_struct"></a>
The `one_time_occurrence` block supports:

* `active_start_date` - Indicates the active start date.

* `active_start_time` - Indicates the active start time.

<a name="frequency_struct"></a>
The `frequency` block supports:

* `freq_type` - Indicates the frequency type. The value can be: **daily**, **weekly**, **monthly_day**, **monthly_week**.

* `freq_interval` - Indicates the frequency interval.

* `freq_interval_weekly` - Indicates the days of the week will be occurred. The value can be:**Monday**, **Tuesday**,
  **Wednesday**, **Thursday**, **Friday**, **Saturday**, **Sunday**.

* `freq_interval_day_monthly` - Indicates the monthly execution date.

* `freq_interval_monthly` - Indicates the days of the month are implemented on a weekly basis. The value can be:
  **Sunday**, **Monday**, **Tuesday**, **Wednesday**, **Thursday**, **Friday**, **Saturday**, **day**, **weekday**,
**weekend**.

* `freq_relative_interval_monthly` - Indicates the week of the month will it be occurred. The value can be: **first**,
  **second**, **third**, **fourth**, **last**.

<a name="daily_frequency_struct"></a>
The `daily_frequency` block supports:

* `freq_subday_type` - Indicates the daily frequency type.

* `active_start_time` - Indicates the daily frequency active start time.

* `active_end_time` - Indicates the daily frequency active end time.

* `freq_subday_interval` - Indicates the daily frequency interval.

* `freq_interval_unit` - Indicates the daily frequency interval unit. The value can be: **second**, **minute**, **hour**.

<a name="duration_struct"></a>
The `duration` block supports:

* `active_start_date` - Indicates the active start date.

* `active_end_date` - Indicates the active end date.

<a name="tables_struct"></a>
The `tables` block supports:

* `table_name` - Indicates the table name.

* `schema` - Indicates the schema.

* `columns` - Indicates the publication columns.

* `primary_key` - Indicates the primary key of the table.

* `filter_statement` - Indicates the filter statement.

* `filter` - Indicates the filter.
  The [filter](#filter_struct) structure is documented below.

* `article_properties` - Indicates the article properties.
  The [article_properties](#article_properties_struct) structure is documented below.

<a name="filter_struct"></a>
The `filter` block supports:

* `relation` - Indicates the relation of the filter. An empty value indicates that the current filter is the lowest level
  filter; a non-empty value indicates that the current item has lower level filters.

* `column` - Indicates the column of the filter.

* `condition` - Indicates the condition of the filter.

* `value` - Indicates the value of the filter.

* `filters` - Indicates the sub-filter of the filter. The value is a json string.

<a name="article_properties_struct"></a>
The `article_properties` block supports:

* `destination_object_name` - Indicates the destination object name.

* `destination_object_owner` - Indicates the destination object owner.

* `insert_delivery_format` - Indicates the insert delivery format. The value can be:
  + **do_not_insert**: o not execute the insert statement.
  + **insert**: execute the insert statement.
  + **insert_without_column_list**: the fields in the insert statement remain in their original order.
  + **call_procedure**: execute the stored procedure to pass all values for all columns.

* `insert_stored_procedure` - Indicates the insert stored procedure.

* `update_delivery_format` - Indicates the update delivery format. The value can be:
   + **do_not_update**: do not execute the update statement.
   + **update**: execute the update statement.
   + **call_procedure**: execute the stored procedure to pass all values for all columns.
   + **mcall_procedure**: execute the stored procedure, only the values of the affected columns are passed; it also
     includes a bitmask representing the modified columns.
   + **xcall_procedure**: execute the stored procedure, passing all columns (whether affected or not) along with the old
     data values for each column.
   + **scall_procedure**: execute a stored procedure only passes the values of the columns that are actually affected by
    the update.

* `update_stored_procedure` - Indicates the update stored procedure.

* `delete_delivery_format` - Indicates the delete delivery format. The value can be:
   + **do_not_delete**: do not execute the delete statement.
   + **delete**: execute the call_procedure statement.
   + **call_procedure**: execute the stored procedure to pass all values for all columns.
   + **xcall_procedure**: execute the stored procedure, passing all columns (whether affected or not) along with the old
    data values for each column.

* `delete_stored_procedure` - Indicates the delete stored procedure.
