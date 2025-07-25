---
subcategory: "Host Security Service (HSS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_hss_backup_policy"
description: |-
  Use this data source to get the HSS backup policy within HuaweiCloud.
---

# huaweicloud_hss_backup_policy

Use this data source to get the HSS backup policy within HuaweiCloud.

## Example Usage

```hcl
variable "policy_id" {}

data "huaweicloud_hss_backup_policy" "test" {
  policy_id = var.policy_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource. If omitted, the provider-level
  region will be used.

* `policy_id` - (Required, String) Specifies the backup policy ID.

* `enterprise_project_id` - (Optional, String) Specifies the ID of the enterprise project that the server belongs to.
  The value **0** indicates the default enterprise project. To query servers in all enterprise projects, set this parameter
  to **all_granted_eps**. If you have only the permission on an enterprise project, you need to transfer the enterprise
  project ID to query the server in the enterprise project. Otherwise, an error is reported due to insufficient permission.

  -> An enterprise project can be configured only after the enterprise project function is enabled.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID in UUID format.

* `enabled` - Indicates whether the backup policy is enabled.

* `name` - The name of the backup policy.

* `operation_type` - The backup type. Currently, only **backup** is supported.

* `operation_definition` - The policy attribute.
  The [operation_definition](#operation_definition_struct) structure is documented below.

* `trigger` - The backup policy scheduling rule.
  The [trigger](#trigger_struct) structure is documented below.

<a name="operation_definition_struct"></a>
The `operation_definition` block supports:

* `day_backups` - The maximum number of daily backups that can be retained.

* `max_backups` - The maximum number of backups that can be automatically created for a backup object.
  If the value is `-1`, backups will not be cleared by quantity limit. If this parameter and `retention_duration_days`
  are left blank at the same time, the backups will be retained permanently. Minimum value: `1`. Maximum value: `99999`.
  Default value: `-1`.

* `month_backups` - The maximum number of monthly backups that can be retained.

* `retention_duration_days` - The duration of retaining a backup, in days. The maximum value is 99999. If the value is
  `-1`, backups will not be cleared by retention duration. If this parameter and `max_backups` are left blank at the
  same time, the backups will be retained permanently.

* `timezone` - The time zone where the user is located. For example, **UTC+08:00**.

* `week_backups` - The maximum number of weekly backups that can be retained.

* `year_backups` - The maximum number of yearly backups that can be retained.

<a name="trigger_struct"></a>
The `trigger` block supports:

* `id` - The scheduler ID.

* `name` - The scheduler name.

* `type` - The scheduler type. Currently, only **time** is supported.

* `properties` - The scheduler attribute.
  The [properties](#properties_struct) structure is documented below.

<a name="properties_struct"></a>
The `properties` block supports:

* `pattern` - The scheduling policy. The value contains a maximum of `10,240` characters and complies with iCalendar RFC
  `2445`. However, only FREQ, BYDAY, BYHOUR, and BYMINUTE are supported. FREQ can be set to only WEEKLY or DAILY. BYDAY
  can be set to the seven days in a week (MO, TU, WE, TH, FR, SA and SU). BYHOUR can be set to `0` to `23` hours. BYMINUTE
  can be set to `0` to `59` minutes. The interval between time points cannot be less than one hour. Multiple backup time
  points can be set in a backup policy, and up to `24` time points can be set for a day.

* `start_time` - The scheduler start time. Example: **2020-01-08 09:59:49**.
