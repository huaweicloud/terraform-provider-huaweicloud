---
subcategory: "Relational Database Service (RDS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_rds_subscriptions"
description: |-
  Use this data source to query the list of RDS subscriptions within HuaweiCloud.
---

# huaweicloud_rds_subscriptions

Use this data source to query the list of RDS subscriptions within HuaweiCloud.

## Example Usage

### Basic Usage

```hcl
variable "instance_id" {}

data "huaweicloud_rds_subscriptions" "test" {
  instance_id = var.instance_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the subscriptions.  
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the ID of the RDS instance.

* `publication_id` - (Optional, String) Specifies the publication ID to filter.  
  If specified, queries subscriptions under this publication. If not specified, queries local subscriptions of the instance.

* `is_cloud` - (Optional, String) Specifies the subscription server source to filter.  
  The valid values are as follows:
  + **true**: The subscription server is a cloud instance.
  + **false**: The subscription server is not a cloud instance.

* `publication_name` - (Optional, String) Specifies the publication name to filter (fuzzy match).

* `subscription_db_name` - (Optional, String) Specifies the subscription database name to filter (fuzzy match).

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `subscriptions` - The list of subscriptions.  
  The [subscriptions](#subscriptions_subscriptions) structure is documented below.

<a name="subscriptions_subscriptions"></a>
The `subscriptions` block supports:

* `id` - The subscription ID.

* `status` - The subscription status.  
  The valid values are as follows:
  + **normal**: Normal.
  + **abnormal**: Abnormal.
  + **creating**: Creating.
  + **createfail**: Creation failed.

* `publication_id` - The publication ID.

* `publication_name` - The publication name.

* `is_cloud` - Whether the subscription server is a cloud instance.

* `subscription_database` - The target database name.

* `subscription_type` - The subscription type. The valid value is **push**.

* `publication_subscription` - The publication subscription information.  
  The [publication_subscription](#subscriptions_publication_subscription) structure is documented below.

* `local_subscription` - The local subscription information.  
  The [local_subscription](#subscriptions_local_subscription) structure is documented below.

* `job_schedule` - The job schedule information.  
  The [job_schedule](#subscriptions_job_schedule) structure is documented below.

<a name="subscriptions_publication_subscription"></a>
The `publication_subscription` block supports:

* `subscription_instance_name` - The subscription server name.

* `subscription_instance_ip` - The subscription server IP.

* `subscription_instance_id` - The subscription instance ID.

<a name="subscriptions_local_subscription"></a>
The `local_subscription` block supports:

* `publication_instance_id` - The publication instance ID when the publication server is a cloud instance.

* `publication_instance_name` - The publication server name.

<a name="subscriptions_job_schedule"></a>
The `job_schedule` block supports:

* `id` - The job schedule ID.

* `job_schedule_type` - The job schedule type.  
  The valid values are as follows:
  + **automatically**: Start automatically when SQL Server agent starts.
  + **cpu_idle**: Start when CPU is idle.
  + **recurring**: Execute repeatedly.
  + **one_time**: Execute once.

* `one_time_occurrence` - The one-time execution time.  
  The [one_time_occurrence](#subscriptions_one_time_occurrence) structure is documented below.

* `frequency` - The frequency interval.  
  The [frequency](#subscriptions_frequency) structure is documented below.

* `daily_frequency` - The daily frequency.  
  The [daily_frequency](#subscriptions_daily_frequency) structure is documented below.

* `duration` - The duration.  
  The [duration](#subscriptions_duration) structure is documented below.

<a name="subscriptions_one_time_occurrence"></a>
The `one_time_occurrence` block supports:

* `active_start_date` - The execution date, format: yyyy-MM-dd.

* `active_start_time` - The execution time, format: HH:mm:ss.

<a name="subscriptions_frequency"></a>
The `frequency` block supports:

* `freq_type` - The frequency type.  
  The valid values are as follows:
  + **daily**: Daily.
  + **weekly**: Weekly.
  + **monthly_day**: Monthly by day.
  + **monthly_week**: Monthly by week.

* `freq_interval` - The execution interval. The valid value is from 1 to 99.

* `freq_interval_weekly` - The days of the week to execute.  
  Valid values are Monday, Tuesday, Wednesday, Thursday, Friday, Saturday, Sunday.

* `freq_interval_day_monthly` - The day of the month to execute.  
  The valid value is from 1 to the total days of the month.

* `freq_interval_monthly` - The days of the week to execute monthly by week.  
  Valid values are Sunday, Monday, Tuesday, Wednesday, Thursday, Friday, Saturday, day, weekday, weekend.

* `freq_relative_interval_monthly` - The week of the month to execute.  
  Valid values are first, second, third, fourth, last.

<a name="subscriptions_daily_frequency"></a>
The `daily_frequency` block supports:

* `freq_subday_type` - The daily frequency type.  
  The valid values are as follows:
  + **once**: Once per day.
  + **multiple**: Multiple times per day.

* `active_start_time` - The first execution time of the day, format: HH:mm:ss.

* `active_end_time` - The last execution time of the day, format: HH:mm:ss.

* `freq_subday_interval` - The execution interval. The valid value is from 1 to 99.

* `freq_interval_unit` - The execution interval unit.  
  The valid values are as follows:
  + **second**: Second.
  + **minute**: Minute.
  + **hour**: Hour.

<a name="subscriptions_duration"></a>
The `duration` block supports:

* `active_start_date` - The first execution date, format: yyyy-MM-dd.

* `active_end_date` - The last execution date, format: yyyy-MM-dd.
