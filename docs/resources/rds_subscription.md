---
subcategory: "Relational Database Service (RDS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_rds_subscription"
description: |-
  Manages an RDS subscription resource within HuaweiCloud.
---

# huaweicloud_rds_subscription

Manages an RDS subscription resource within HuaweiCloud.

## Example Usage

```hcl
variable "instance_id" {}
variable "publication_id" {}

resource "huaweicloud_rds_subscription" "test" {
  instance_id           = var.instance_id
  subscription_database = "test_db"
  subscription_type     = "push"
  initialize_at         = "immediate"

  local_subscription {
    publication_id   = var.publication_id
    publication_name = "test_pub"
  }
  
  job_schedule {
    job_schedule_type = "recurring"
    frequency {
      freq_type     = "daily"
      freq_interval = 1
    }
    daily_frequency {
      freq_subday_type   = "once"
      active_start_time = "12:00:00"
    }
    duration {
      active_start_date = "2024-01-01"
      active_end_date   = "2024-12-31"
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the subscription.  
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `instance_id` - (Required, String, NonUpdatable) Specifies the ID of the RDS instance.  

* `subscription_database` - (Required, String, NonUpdatable) Specifies the subscription database name.  

* `subscription_type` - (Required, String, NonUpdatable) Specifies the subscription type. The valid value is **push**.

* `initialize_at` - (Required, String, NonUpdatable) Specifies the initialization type.  
  The valid values are as follows:
  + **do_not**: Do not initialize.
  + **immediate**: Initialize immediately.
  + **at_first_sync**: Initialize at first synchronization.
  
* `job_schedule` - (Required, List) Specifies the synchronization strategy.  
  The [job_schedule](#job_schedule_struct) structure is documented below.

* `local_subscription` - (Required, List, NonUpdatable) Specifies the local subscription information.
  The [local_subscription](#local_subscription_struct) structure is documented below.

* `initialize_info` - (Optional, List, NonUpdatable) Specifies the initialization information.  
  The [initialize_info](#initialize_info_struct) structure is documented below.  

* `independent_agent` - (Optional, String, NonUpdatable) Specifies whether to use an independent distribution agent.  
  The valid values are **true** and **false**. Defaults to **true**.  

* `bak_file_name` - (Optional, String, NonUpdatable) Specifies the backup file name.  
  If this value is not empty, the subscription initialization method is through backup file initialization.  

* `bak_bucket_name` - (Optional, String, NonUpdatable) Specifies the OBS bucket name where the backup file is located.  

<a name="job_schedule_struct"></a>
The `job_schedule` block supports:

* `job_schedule_type` - (Optional, String) Specifies the schedule type.  
  The valid values are as follows:
  + **cpu_idle**: Start when CPU is idle.
  + **recurring**: Execute repeatedly.
  + **one_time**: Execute once.
  + **automatically**: Automatically start when SQL Server agent starts.
  
  Defaults to **recurring**.

* `one_time_occurrence` - (Optional, List) Specifies the one-time execution time.  
  The [one_time_occurrence](#one_time_occurrence_struct) structure is documented below.

* `frequency` - (Optional, List) Specifies the strategy interval period.  
  Required when strategy ID is empty.  
  The [frequency](#frequency_struct) structure is documented below.

* `daily_frequency` - (Optional, List) Specifies the strategy daily frequency.  
  Required when strategy ID is empty.  
  The [daily_frequency](#daily_frequency_struct) structure is documented below.

* `duration` - (Optional, List) Specifies the strategy validity period.  
  Required when strategy ID is empty.  
  The [duration](#duration_struct) structure is documented below.

<a name="one_time_occurrence_struct"></a>
The `one_time_occurrence` block supports:

* `active_start_date` - (Optional, String) Specifies the execution date. The format is **yyyy-MM-dd**.  
  The valid value range is from **1990-01-01** to **2099-12-31**.

* `active_start_time` - (Optional, String) Specifies the execution time. The format is **HH:mm:ss**.

<a name="frequency_struct"></a>
The `frequency` block supports:

* `freq_type` - (Optional, String) Specifies the strategy frequency type.  
  The valid values are as follows:
  + **daily**: By day.
  + **weekly**: By week.
  + **monthly_day**: By month, by day of month.
  + **monthly_week**: By month, by week of month.

* `freq_interval` - (Optional, Int) Specifies the execution interval.  
  The valid value range is from **1** to **99**.

* `freq_interval_weekly` - (Optional, List) Specifies which days of the week to execute.  
  Required when frequency type is weekly.  
  The valid values are **Monday**, **Tuesday**, **Wednesday**, **Thursday**, **Friday**, **Saturday**, **Sunday**.

* `freq_interval_day_monthly` - (Optional, Int) Specifies the date of each month to execute.  
  Required when frequency type is monthly_day.  
  The valid value range is from **1** to the total number of days in the current month.

* `freq_interval_monthly` - (Optional, String) Specifies which days of the week to execute in the current month.  
  Required when frequency type is monthly_week.  
  The valid values are **Sunday**, **Monday**, **Tuesday**, **Wednesday**, **Thursday**, **Friday**, **Saturday**,
  **day**, **weekday**, **weekend**.

* `freq_relative_interval_monthly` - (Optional, String) Specifies which week of each month to execute.  
  Required when frequency type is monthly_week.  
  The valid values are **first**, **second**, **third**, **fourth**, **last**.

<a name="daily_frequency_struct"></a>
The `daily_frequency` block supports:

* `freq_subday_type` - (Optional, String) Specifies the daily frequency type.  
  The valid values are as follows:
  + **once**: Once a day.
  + **multiple**: Multiple times a day.

* `active_start_time` - (Optional, String) Specifies the first execution time of each day.  
  When the daily frequency type is once, only this execution will occur. The format is **HH:mm:ss**.

* `active_end_time` - (Optional, String) Specifies the last execution time. The format is **HH:mm:ss**.  
  Required when executing multiple times a day.

* `freq_subday_interval` - (Optional, Int) Specifies the execution interval.  
  Required when executing multiple times a day.  
  The valid value range is from **1** to **99**.

* `freq_interval_unit` - (Optional, String) Specifies the unit of execution interval.  
  Required when executing multiple times a day.  
  The valid values are **second**, **minute**, **hour**.

<a name="duration_struct"></a>
The `duration` block supports:

* `active_start_date` - (Optional, String) Specifies the first execution date. The format is **yyyy-MM-dd**.  
  The valid value range is from **1990-01-01** to **2099-12-31**.

* `active_end_date` - (Optional, String) Specifies the last execution date. The format is **yyyy-MM-dd**.  
  Defaults to never ending.

<a name="initialize_info_struct"></a>
The `initialize_info` block supports:

* `file_source` - (Optional, String, NonUpdatable) Specifies the file source used for initialization.  
  The valid values are **OBS** and **BACKUP**.

* `backup_id` - (Optional, String, NonUpdatable) Specifies the backup file ID for initialization using a backup file.

* `bucket_name` - (Optional, String, NonUpdatable) Specifies the bucket name for restoring using a backup file in OBS.

* `file_path` - (Optional, String, NonUpdatable) Specifies the path of the backup file in the OBS bucket.

* `file_name` - (Optional, String, NonUpdatable) Specifies the name of the backup file in the OBS bucket.

* `overwrite_restore` - (Optional, String, NonUpdatable) Specifies whether to overwrite and restore the subscription
  database using the backup file. The valid values are **true** and **false**.

<a name="local_subscription_struct"></a>
The `local_subscription` block supports:

* `publication_id` - (Required, String, NonUpdatable) Specifies the publication ID.

* `publication_name` - (Optional, String, NonUpdatable) Specifies the publication name.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `job_schedule` - The synchronization strategy.  
  The [job_schedule](#job_schedule_attr) structure is documented below.

* `local_subscription` - The local subscription information.
  The [local_subscription](#local_subscription_attr) structure is documented below.

* `status` - The subscription status.  
  The valid values are **normal**, **abnormal**, **creating**, **createfail**.

* `is_cloud` - The source of the subscriber.

  <a name="job_schedule_attr"></a>
The `job_schedule` block supports:

* `id` - The schedule ID.

<a name="local_subscription_attr"></a>
The `local_subscription` block supports:

* `publication_instance_id` - The publisher instance ID

* `publication_instance_name` - The publisher name.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 30 minutes.
* `update` - Default is 30 minutes.
* `delete` - Default is 30 minutes.

## Import

RDS subscription can be imported using the `instance_id` and `id` separated by a slash, e.g.

```bash
$ terraform import huaweicloud_rds_subscription.test <instance_id>/<id>
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response, security or some other reason. The missing attributes include: `initialize_at`, `initialize_info`,
`independent_agent`, `bak_file_name`, `bak_bucket_name`. It is generally recommended running `terraform plan` after
importing the resource. You can then decide if changes should be applied to the subscription, or the resource definition
should be updated to align with the subscription. Also you can ignore changes as below.

```hcl
resource "huaweicloud_rds_subscription" "test" {
  ...

  lifecycle {
    ignore_changes = [
      initialize_at, initialize_info, independent_agent, bak_file_name, bak_bucket_name
    ]
  }
}
```
