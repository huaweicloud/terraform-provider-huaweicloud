---
subcategory: "Cloud Backup and Recovery (CBR)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cbr_policy"
description: ""
---

# huaweicloud_cbr_policy

Manages a backup policy for backing up vault objects within HuaweiCloud.

## Example Usage

### Create a backup policy (weekly backup)

```hcl
variable "policy_name" {}

resource "huaweicloud_cbr_policy" "test" {
  name        = var.policy_name
  type        = "backup"
  time_period = 20
  time_zone   = "UTC+08:00"

  long_term_retention {
    daily                = 10
    weekly               = 10
    monthly              = 1
    full_backup_interval = -1
  }

  backup_cycle {
    days            = "SA,SU"
    execution_times = ["08:00", "20:00"]
  }
}
```

### Create a replication policy (periodic backup)

```hcl
variable "policy_name" {}
variable "destination_region" {}
variable "destination_project_id" {}

resource "huaweicloud_cbr_policy" "test" {
  name                   = var.policy_name
  type                   = "replication"
  destination_region     = var.destination_region
  destination_project_id = var.destination_project_id
  backup_quantity        = 20

  backup_cycle {
    interval        = 5
    execution_times = ["21:00"]
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region where the policy is located. If omitted, the
  provider-level region will be used. Changing this will create a new policy.

* `name` - (Required, String) Specifies the policy name.  
  This parameter can contain a maximum of `64`
  characters, which may consist of chinese characters, letters, digits, underscores(_) and hyphens (-).

* `type` - (Required, String, ForceNew) Specifies the protection type of the policy.
  Valid values are **backup** and **replication**.
  Changing this will create a new policy.

* `backup_cycle` - (Required, List) Specifies the scheduling rule for the policy backup execution.
  The [object](#cbr_policy_backup_cycle) structure is documented below.

* `enabled` - (Optional, Bool) Specifies whether to enable the policy. Default to **true**.

* `destination_region` - (Optional, String) Specifies the name of the replication destination region, which is mandatory
  for cross-region replication. Required if `protection_type` is **replication**.

* `destination_project_id` - (Optional, String) Specifies the ID of the replication destination project, which is
  mandatory for cross-region replication. Required if `protection_type` is **replication**.

* `backup_quantity` - (Optional, Int) Specifies the maximum number of retained backups. The value ranges from `2` to
  `99,999`. This parameter and `time_period` are alternative.

* `time_period` - (Optional, Int) Specifies the duration (in days) for retained backups. The value ranges from `2` to
  `99,999`.

-> **NOTE:** If this `backup_quantity` and `time_period` are both left blank, the backups will be retained permanently.

* `long_term_retention` - (Optional, List) Specifies the long-term retention rules, which is an advanced options of
  the `backup_quantity`. The [object](#cbr_policy_long_term_retention) structure is documented below.

-> The configuration of `long_term_retention` and `backup_quantity` will take effect together.
  When the number of retained backups exceeds the preset value (number of `backup_quantity`), the system automatically
  deletes the earliest backups. By default, the system automatically clears data every other day.

* `time_zone` - (Optional, String) Specifies the UTC time zone, e.g. `UTC+08:00`.
  Only available if `long_term_retention` is set.

<a name="cbr_policy_backup_cycle"></a>
The `backup_cycle` block supports:

* `days` - (Optional, String) Specifies the weekly backup day of backup schedule.  
  It supports seven days a week (MO, TU, WE, TH, FR, SA, SU) and this parameter is separated by a comma (,) without
  spaces, between date and date during the configuration.

* `interval` - (Optional, Int) Specifies the interval (in days) of backup schedule.  
  The valid value ranges from `1` to `30`.

-> Parameter `backup_cycle.days` and `backup_cycle.interval` are alternative.

* `execution_times` - (Required, List) Specifies the backup time. Automated backups will be triggered at the backup
  time. The current time is in the UTC format (HH:MM). The minutes in the list must be set to **00** and the hours
  cannot be repeated. In the replication policy, you are advised to set one time point for one day.

<a name="cbr_policy_long_term_retention"></a>
The `long_term_retention` block supports:

* `daily` - (Optional, Int) - Specifies the latest backup of each day is saved in the long term.

* `weekly` - (Optional, Int) - Specifies the latest backup of each week is saved in the long term.

* `monthly` - (Optional, Int) - Specifies the latest backup of each month is saved in the long term.

* `yearly` - (Optional, Int) - Specifies the latest backup of each year is saved in the long term.

-> A maximum of 10 backups are retained for failed periodic backup tasks. They are retained for one month and can be
  manually deleted on the web console.

* `full_backup_interval` - (Optional, Int) Specifies how often (after how many incremental backups) a full backup is
  performed. The valid value ranges from `-1` to `100`.
  If `-1` is specified, full backup will not be performed.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - A resource ID in UUID format.

## Import

Policies can be imported by their `id`, e.g.

```bash
$ terraform import huaweicloud_cbr_policy.test <id>
```
