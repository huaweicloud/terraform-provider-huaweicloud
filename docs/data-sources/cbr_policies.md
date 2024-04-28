---
subcategory: "Cloud Backup and Recovery (CBR)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cbr_policies"
description: ""
---

# huaweicloud_cbr_policies

Use this data source to get available CBR policies within HuaweiCloud.

## Example Usage

```hcl
variable policy_name {}

data "huaweicloud_cbr_policies" "test" {
  name = var.policy_name
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the policies.
  If omitted, the provider-level region will be used.

* `policy_id` - (Optional, String) Specifies the policy ID used to query.

* `name` - (Optional, String) Specifies the policy name used to query.

* `type` - (Optional, String) Specifies the policy type used to query. The valid values are as follows:
  + **backup**: Backup policy
  + **replication**: Replication policy

* `enabled` - (Optional, Bool) Specifies the policy enabling status to query. The valid values are as follows:
  + **true**: Policy enabled
  + **false**: Policy not enabled

* `vault_id` - (Optional, String) Specifies the vault ID of the associated policy used to query.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `policies` - All CBR policies that match the filter parameters.
  The [policies](#CBR_Policies) structure is documented below.

<a name="CBR_Policies"></a>
The `policies` block supports:

* `id` - The policy ID.

* `name` - The policy name.

* `type` - The protection type of the policy. The valid values are as follows:
  + **backup**: Backup policy
  + **replication**: Replication policy

* `enabled` - Whether to enable the policy. The valid values are as follows:
  + **true**: Policy enabled
  + **false**: Policy not enabled

* `backup_cycle` - The scheduling rule for the policy backup execution.
  The [backup_cycle](#CBR_Policies_BackupCycle) structure is documented below.

* `destination_region` - The name of the replication destination region.

* `enable_acceleration` - Whether to enable the acceleration function to shorten the replication time for cross-region.
  The valid values are as follows:
  + **true**: Enabled acceleration
  + **false**: Not enabled acceleration

* `destination_project_id` - The ID of the replication destination project.

* `backup_quantity` - The maximum number of retained backups. The value ranges from `2` to `99,999`.
  This parameter and `time_period` are alternative.

* `time_period` - The duration (in days) for retained backups. The value ranges from `2` to `99,999`.

* `long_term_retention` - The long-term retention rules, which is an advanced options of the `backup_quantity`.
  The [long_term_retention](#CBR_Policies_LongTermRetention) structure is documented below.

* `time_zone` - The UTC time zone, e.g. `UTC+08:00`. Only available when `long_term_retention` is set.

* `associated_vaults` - The vault associated with the CBR policy.
  The [associated_vaults](#CBR_Policies_AssociatedVaults) structure is documented below.

<a name="CBR_Policies_BackupCycle"></a>
The `backup_cycle` block supports:

* `interval` - The interval (in days) of backup schedule. The value range is `1` to `30`.

* `days` - The weekly backup day of backup schedule. It supports seven days a week (MO, TU, WE, TH, FR, SA, SU)
  and this parameter is separated by a comma (,) without spaces between the date and date.

* `execution_times` - The backup time. Automated backups will be triggered at the backup
  time. The current time is in the UTC format (HH:MM).

<a name="CBR_Policies_LongTermRetention"></a>
The `long_term_retention` block supports:

* `daily` - The latest backup of each day is saved in the long term.

* `weekly` - The latest backup of each week is saved in the long term.

* `monthly` - The latest backup of each month is saved in the long term.

* `yearly` - The latest backup of each year is saved in the long term.

* `full_backup_interval` - How often (after how many incremental backups) a full backup is performed.
  The valid value ranges from `-1` to `100`. If `-1` is specified, full backup will not be performed.

<a name="CBR_Policies_AssociatedVaults"></a>
The `associated_vaults` block supports:

* `vault_id` - The vault ID of the associated CBR policy.

* `destination_vault_id` - The destination vault ID associated with CBR policy.
