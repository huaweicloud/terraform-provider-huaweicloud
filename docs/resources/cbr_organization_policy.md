---
subcategory: "Cloud Backup and Recovery (CBR)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cbr_organization_policy"
description: |-
  Manages a CBR organization policy resource within HuaweiCloud.
---

# huaweicloud_cbr_organization_policy

Manages a CBR organization policy resource within HuaweiCloud.

## Example Usage

### Create a backup organization policy with retention settings

```hcl
variable "organization_policy_name" {}
variable "policy_name" {}

resource "huaweicloud_cbr_organization_policy" "test" {
  name           = var.organization_policy_name
  description    = "Created by terraform script"
  operation_type = "backup"
  policy_name    = var.policy_name
  policy_enabled = false

  policy_operation_definition {
    day_backups          = 5
    max_backups          = 30
    month_backups        = 1
    week_backups         = 2
    year_backups         = 1
    timezone             = "UTC+08:00"
    full_backup_interval = 10
  }

  policy_trigger {
    properties {
      pattern = ["FREQ=WEEKLY;BYDAY=WE,TH,FR;BYHOUR=16;BYMINUTE=00"]
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `name` - (Required, String) Specifies the organization policy name.

* `operation_type` - (Required, String, ForceNew) Specifies the organization policy type.  
  The valid values are as follows:
  + **backup**
  + **replication**.

* `policy_name` - (Required, String) Specifies the policy name.  
  The CBR service will automatically create a corresponding policy based on this name.

* `policy_enabled` - (Required, Bool) Specifies whether the policy is enabled.

* `policy_operation_definition` - (Required, List) Specifies the policy operation definition for backup and replication.
  The [policy_operation_definition](#cbr_organization_policy_operation_definition) structure is documented below.

* `policy_trigger` - (Required, List) Specifies the policy execution time rule.
  The [policy_trigger](#cbr_organization_policy_trigger) structure is documented below.

* `description` - (Optional, String) Specifies the organization policy description.

* `effective_scope` - (Optional, String) Specifies the effective scope of the organization policy.

<a name="cbr_organization_policy_operation_definition"></a>
The `policy_operation_definition` block supports:

* `day_backups` - (Optional, Int) Specifies the maximum number of daily backups that can be retained.
  The latest backup of each day is saved in the long term. The value ranges from **0** to **100**.

* `destination_project_id` - (Optional, String) Specifies the destination project ID for replication.
  This parameter is **mandatory** for cross-region replication.

* `destination_region` - (Optional, String) Specifies the destination region for replication.
  This parameter is **mandatory** for cross-region replication.

* `enable_acceleration` - (Optional, String) Specifies whether to enable acceleration to shorten replication time for
  cross-region replication.  
  The valid values are as follows:
  + **true**
  + **false**

* `max_backups` - (Optional, Int) Specifies the maximum number of backups that can be automatically created for a
  backup object. The value can be **-1** or ranges from **0** to **99999**. If the value is set to **-1**, backups
  will not be cleared by quantity limit.

* `month_backups` - (Optional, Int) Specifies the maximum number of monthly backups that can be retained.
  The latest backup of each month is saved in the long term. The value ranges from **0** to **100**.

* `retention_duration_days` - (Optional, Int) Specifies the duration of retaining a backup, in days.
  The maximum value is **99999**. If the value is set to **-1**, backups will not be cleared by retention duration.

* `week_backups` - (Optional, Int) Specifies the maximum number of weekly backups that can be retained.
  The latest backup of each week is saved in the long term. The value ranges from **0** to **100**.

* `year_backups` - (Optional, Int) Specifies the maximum number of yearly backups that can be retained.
  The latest backup of each year is saved in the long term. The value ranges from **0** to **100**.

* `timezone` - (Optional, String) Specifies the time zone where the user is located, for example, **UTC+08:00**.

* `full_backup_interval` - (Optional, Int) Defines how often a full backup is performed after incremental backups.
  If **-1** is specified, full backup will not be performed. The value ranges from **-1** to **100**.

<a name="cbr_organization_policy_trigger"></a>
The `policy_trigger` block supports:

* `properties` - (Required, List) Specifies the properties of policy trigger.
  The [properties](#cbr_organization_policy_trigger_properties) structure is documented below.

<a name="cbr_organization_policy_trigger_properties"></a>
The `properties` block supports:

* `pattern` - (Required, List) Specifies the scheduling rules for policy execution. Up to 24 rules are supported.
  The scheduling rules follow the iCalendar RFC 2445 specification, supporting parameters like **FREQ**, **BYDAY**,
  **BYHOUR**, **BYMINUTE**, and **INTERVAL**. For example: **FREQ=WEEKLY;BYDAY=MO,TU,WE,TH,FR,SA,SU;BYHOUR=14;BYMINUTE=00**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `status` - The organization policy status.

* `domain_id` - The ID of the account to which the organization policy belongs.

* `domain_name` - The account to which the organization policy belongs.

## Import

The CBR organization policy can be imported using the `id` or `name`, e.g.

```bash
$ terraform import huaweicloud_cbr_organization_policy.test <id>
```

```bash
$ terraform import huaweicloud_cbr_organization_policy.test <name>
```
