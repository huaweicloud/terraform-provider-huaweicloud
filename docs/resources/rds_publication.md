---
subcategory: "Relational Database Service (RDS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_rds_publication"
description: |-
  Manage an RDS publication resource within HuaweiCloud.
---

# huaweicloud_rds_publication

Manage an RDS publication resource within HuaweiCloud.

## Example Usage

```hcl
variable "instance_id" {}

resource "huaweicloud_rds_publication" "test" {
  instance_id                    = var.instance_id
  publication_name               = "test_publication_name"
  publication_database           = "test_db"
  is_create_snapshot_immediately = "true"

  subscription_options {
    independent_agent            = "true"
    snapshot_always_available    = "false"
    replicate_ddl                = "true"
    allow_initialize_from_backup = "false"
  }

  job_schedule{
    job_schedule_type = "recurring"

    one_time_occurrence {
      active_start_date = "2026-05-06"
      active_start_time = "07:20:30"
    }

    frequency {
      freq_type                      = "monthly_week"
      freq_interval                  = 2
      freq_interval_monthly          = "weekday"
      freq_relative_interval_monthly = "second"
    }

    daily_frequency {
      freq_subday_type     = "multiple"
      active_start_time    = "02:00:00"
      active_end_time      = "20:00:00"
      freq_subday_interval = "3"
      freq_interval_unit   = "minute"
    }

    duration {
      active_start_date = "2010-07-15"
      active_end_date   = "2089-10-20"
    }
  }

  is_select_all_table = "false"

  tables {
    table_name  = "test_table_1"
    schema      = "test_schema_1"
    columns     = ["id", "name", "address"]
    primary_key = ["id"]

    filter {
      relation = "AND"
      filters  = [
        jsonencode({
          "column": "id",
          "condition": "=",
          "value": "111"
        }),
        jsonencode({
          "relation": "AND",
          "filters": [
            {
              "column": "id",
              "condition": "=",
              "value": "222"
            },
            {
              "column": "name",
              "condition": "=",
              "value": "123"
            }
          ]
        })
      ]
    }

    article_properties {
      destination_object_name  = "test_table_1"
      destination_object_owner = "test_schema_1"
      insert_delivery_format   = "call_procedure"
      insert_stored_procedure  = "sp_MSins_test_schema_1test_table_1"
      update_delivery_format   = "scall_procedure"
      update_stored_procedure  = "sp_MSupd_test_schema_1test_table_1"
      delete_delivery_format   = "call_procedure"
      delete_stored_procedure  = "sp_MSdel_test_schema_1test_table_1"
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) The region in which to create the rds PostgreSQL SQL limit resource. If omitted,
  the provider-level region will be used. Changing this creates a new resource.

* `instance_id` - (Required, String, NonUpdatable) Specifies the ID of RDS instance.

* `publication_name` - (Required, String, NonUpdatable) Specifies the publication name. The publication name must consist
  of 5 to 64 characters. Only letters, digits, and underscores (_) are allowed.

* `publication_database` - (Required, String, NonUpdatable) Specifies the publication database name.

* `is_create_snapshot_immediately` - (Required, String, NonUpdatable) Specifies whether a snapshot is created immediately.
  Value options:
  + **true**: A snapshot is created immediately.
  + **false**: A snapshot is not created immediately.

* `subscription_options` - (Optional, List) Specifies the subscription options.
  The [subscription_options](#subscription_options_struct) structure is documented below.

* `job_schedule` - (Required, List) Specifies the schedule details.
  The [job_schedule](#job_schedule_struct) structure is documented below.

* `is_select_all_table` - (Optional, String) Specifies whether all data tables are selected.
  Value options:
  + **true**: All tables are selected.
  + **false(default)**: Not all tables are selected.

* `extend_tables` - (Optional, List) Specifies the tables to be removed after all tables are selected.

* `tables` - (Optional, List) Specifies the published tables.
  The [tables](#tables_struct) structure is documented below.

<a name="subscription_options_struct"></a>
The `subscription_options` block supports:

* `independent_agent` - (Optional, String) Specifies whether an independent distribution agent is used. Value options:
  + **true(default)**: An agent is used.
  + **false**: No agent is used.

* `snapshot_always_available` - (Optional, String) Specifies whether snapshots are always available. Value options:
  + **true(default)**: Snapshots are always available.
  + **false**: Snapshots are not always available.

* `replicate_ddl` - (Optional, String) Specifies whether schema changes can be replicated. Value options:
  + **true(default)**: Schema changes can be replicated.
  + **false**: Schema changes cannot be replicated.

* `allow_initialize_from_backup` - (Optional, String) Specifies whether backup files can be used for initialization.
  Value options:
  + **true**: Backup files can be used for initialization.
  + **false(default)**: Backup files cannot be used for initialization.

<a name="job_schedule_struct"></a>
The `job_schedule` block supports:

* `job_schedule_type` - (Optional, String) Specifies the agent schedule type. Value options:
  + **recurring**: The task is executed repeatedly.
  + **one_time**: The task is executed only once.

* `one_time_occurrence` - (Optional, List) Specifies the execution time when the task is executed only once. This parameter
  is mandatory when `job_schedule_type` is set to **one_time**.
  The [one_time_occurrence](#one_time_occurrence_struct) structure is documented below.

* `frequency` - (Optional, List) Specifies the interval of the schedule. This parameter is mandatory only when
  `job_schedule_type` is set to **recurring**.
  The [frequency](#frequency_struct) structure is documented below.

* `daily_frequency` - (Optional, List) Specifies the daily frequency of the schedule. This parameter is mandatory when
  `job_schedule_type` is set to **recurring**.
  The [daily_frequency](#daily_frequency_struct) structure is documented below.

* `duration` - (Optional, List) Specifies the validity period of the schedule. This parameter is mandatory when
  `job_schedule_type` is set to **recurring**.
  The [duration](#duration_struct) structure is documented below.

<a name="one_time_occurrence_struct"></a>
The `one_time_occurrence` block supports:

* `active_start_date` - (Optional, String) Specifies the execution date, in the format of **yyyy-MM-dd**.

* `active_start_time` - (Optional, String) Specifies the execution time, in the format of **HH:mm:ss**.

<a name="frequency_struct"></a>
The `frequency` block supports:

* `freq_type` - (Optional, String) Specifies the frequency type of the schedule. Value options:
  + **daily**: by day
  + **weekly**: by week
  + **monthly_day**: by month and by day in each month
  + **monthly_week**: by month and by week in each month

* `freq_interval` - (Optional, Int) Specifies the execution interval. Value range: **1–99**.

* `freq_interval_weekly` - (Optional, List) Specifies the days in a week when the task is executed. This parameter is
  mandatory when `freq_type` is set to **weekly**. If `freq_type` is not set to **weekly**, this parameter does not take
  effect. Value options: **Monday**, **Tuesday**, **Wednesday**, **Thursday**, **Friday**, **Saturday**, **Sunday**.

* `freq_interval_day_monthly` - (Optional, Int) Specifies the monthly execution dates. his parameter is mandatory when
  `freq_type` is set to **monthly_day**. If `freq_type` is not set to **monthly_day**, this parameter does not take effect.
  Value range: **1** to the total number of days in the month, for example, **1** to **31**.

* `freq_interval_monthly` - (Optional, String) Specifies the days in a week when the task is executed in the current month.
  This parameter is mandatory when `freq_type` is set to **monthly_week**. If `freq_type` is not set to **monthly_week**,
  this parameter does not take effect. Value options: **Sunday**, **Monday**, **Tuesday**, **Wednesday**, **Thursday**,
  **Friday**, **Saturday**, **day**, **weekday**, **weekday**.

* `freq_relative_interval_monthly` - (Optional, String) Specifies the week in a month when the task is executed. This
  parameter is mandatory when `freq_type` is set to **monthly_week**. If `freq_type` is not set to **monthly_week**,
  this parameter does not take effect. Value options: **first**, **second**, **third**, **fourth**, **last**.

<a name="daily_frequency_struct"></a>
The `daily_frequency` block supports:

* `freq_subday_type` - (Required, String) Specifies the daily frequency type. Value options:
  + **once**: once a day
  + **multiple**: multiple times a day

* `active_start_time` - (Optional, String) Specifies the time of the first execution on each day. If `freq_subday_type`
  is set to **once**, the task is executed only once a day. The value is in **HH:mm:ss** format.

* `active_end_time` - (Optional, String) Specifies the last execution time, in the format of **HH:mm:ss**. This parameter
  is mandatory when `freq_subday_type` is set to **multiple**. It does not take effect when `freq_subday_type` is set to
  **once**.

* `freq_subday_interval` - (Optional, Int) Specifies the execution interval. This parameter is mandatory when
  `freq_subday_type` is set to **multiple**. It does not take effect when `freq_subday_type` is set to **once**.
  Value range: **1–99**.

* `freq_interval_unit` - (Optional, String) Specifies the execution interval unit. This parameter is mandatory when
  `freq_subday_type` is set to **multiple**. It does not take effect when `freq_subday_type` is set to **once**.
  Value options: **second**, **minute**, **hour**

<a name="duration_struct"></a>
The `duration` block supports:

* `active_start_date` - (Optional, String) Specifies the first execution date, in the format of **yyyy-MM-dd**.
  Value range: **1990-01-01** to **2099-12-31**.

* `active_end_date` - (Optional, String) Specifies the last execution date, in the format of **yyyy-MM-dd**. If this
  parameter is not specified, the execution does not end.

<a name="tables_struct"></a>
The `tables` block supports:

* `table_name` - (Required, String) Specifies the table name.

* `schema` - (Optional, String) Specifies the schema name. Defaults to **dbo**.

* `columns` - (Optional, List) Specifies the published fields. If this parameter is empty, all fields are selected.

* `primary_key` - (Optional, List) Specifies the primary key.

* `filter_statement` - (Optional, String) Specifies the filter statement.

* `filter` - (Optional, List) Specifies the filter. This parameter is valid only when `filter_statement` is empty.
  If the filter statement and filter are both empty, no filter condition is available.
  The [filter](#filter_struct) structure is documented below.

* `article_properties` - (Optional, List) Specifies the project properties.
  The [article_properties](#article_properties_struct) structure is documented below.

<a name="filter_struct"></a>
The `filter` block supports:

* `relation` - (Optional, String) Specifies the Filter relationship. If the value is empty, the current filter is the
  lowest-level filter. If the value is not empty, there are lower-level filters. Value options: **AND**, **OR**.

* `column` - (Optional, String) Specifies the filter field. This parameter is valid only when `relation` is empty. If
  `relation` is empty, this parameter is mandatory.

* `condition` - (Optional, String) Specifies the filter condition. This parameter is valid only when `relation` is empty.
  If `relation` is empty, this parameter is mandatory. Value options: **=**, **!=**, **>**, **<**, **>=**, **<=**,
  **LIKE**, **NOT LIKE**, **IN**.

* `value` - (Optional, String) Specifies the filter value. This parameter is valid only when `relation` is empty. If
  `relation` is empty, this parameter is mandatory.

* `filters` - (Optional, List) Specifies the lower-level filter. This parameter is valid only when `relation` is not empty.
  If `relation` is not empty, this parameter is mandatory. The element is a json string.

<a name="article_properties_struct"></a>
The `article_properties` block supports:

* `destination_object_name` - (Optional, String) Specifies the name of the target object.

* `destination_object_owner` - (Optional, String) Specifies the namespace of the target object.

* `insert_delivery_format` - (Optional, String) Specifies the INSERT delivery format. Value options:
  + **do_not_insert**: Do not execute the INSERT statement.
  + **insert**: Execute the INSERT statement.
  + **insert_without_column_list**: The fields of the INSERT statement remain in the original order.
  + **call_procedure**: Execute the stored procedure to pass all values for all columns.

* `insert_stored_procedure` - (Optional, String) Specifies the INSERT stored procedure. This parameter is mandatory when
  `insert_delivery_format` is set to **call_procedure**.

* `update_delivery_format` - (Optional, String) Specifies the UPDATE delivery format. Value options:
  + **do_not_update**: Do not execute the UPDATE statement.
  + **update**: Execute the UPDATE statement.
  + **call_procedure**: Execute the stored procedure to pass all values for all columns.
  + **mcall_procedure**: Execute the stored procedure to only pass values for affected columns. It also includes a bitmask
    representing the changed columns.
  + **xcall_procedure**: Execute the stored procedure to pass all columns (whether affected or not) and the old data
    values for each column.
  + **scall_procedure**: Execute the stored procedure to pass values only for the columns that were actually affected by
    the update.

* `update_stored_procedure` - (Optional, String) Specifies the UPDATE stored procedure. This parameter is mandatory when
  `update_delivery_format` is set to **call_procedure**, **mcall_procedure**, **xcall_procedure**, or **scall_procedure**.

* `delete_delivery_format` - (Optional, String) Specifies the DELETE delivery format. Value options:
  + **do_not_delete**: Do not execute the DELETE statement.
  + **delete**: Execute the DELETE statement.
  + **call_procedure**: Execute the stored procedure to pass all values for all columns.
  + **xcall_procedure**: Execute the stored procedure to pass all columns (whether affected or not) and the old data
    values for each column.

* `delete_stored_procedure` - (Optional, String) Specifies the DELETE stored procedure. This parameter is mandatory when
  `delete_delivery_format` is set to **call_procedure** or **xcall_procedure**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - Indicates the resource ID which is same as `instance_id`.

* `job_schedule` - Indicates the schedule details.
  The [job_schedule](#job_schedule_attribute) structure is documented below.

* `status` - Indicates the publication status.

* `subscription_count` - Indicates the number of subscriptions.

<a name="job_schedule_attribute"></a>
The `job_schedule` block supports:

* `id` -Indicates the schedule ID.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 30 minutes.
* `update` - Default is 30 minutes.
* `delete` - Default is 30 minutes.

## Import

The RDS publication can be imported using the `instance_id` and `id` separated by a slash, e.g.

```bash
$ terraform import huaweicloud_rds_publication.test <instance_id>/<id>
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response, security or some other reason. The missing attributes include: `is_create_snapshot_immediately`. It is
generally recommended running `terraform plan` after importing the RDS publication. You can then decide if changes should
be applied to the RDS publication, or the resource definition should be updated to align with the RDS publication. Also
you can ignore changes as below.

```hcl
resource "huaweicloud_rds_publication" "test" {
    ...

  lifecycle {
    ignore_changes = [
      is_create_snapshot_immediately,
    ]
  }
}
```
