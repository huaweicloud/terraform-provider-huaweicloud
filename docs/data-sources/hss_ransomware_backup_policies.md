---
subcategory: "Host Security Service (HSS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_hss_ransomware_backup_policies"
description: |-
  Use this data source to get the list of HSS ransomware backup policies within HuaweiCloud.
---

# huaweicloud_hss_ransomware_backup_policies

Use this data source to get the list of HSS ransomware backup policies within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_hss_ransomware_backup_policies" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project ID.  
  This parameter is valid only when the enterprise project is enabled.
  The default value is **0**, indicating the default enterprise project.
  If you need to query data for all enterprise projects, the value is **all_granted_eps**.
  If you only have permissions for a specific enterprise project, you need set the enterprise project ID. Otherwise,
  the operation may fail due to insufficient permissions.

* `policy_id` - (Optional, String) Specifies the protection policy ID.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID in UUID format.

* `total_num` - The total number.

* `data_list` - The list of ransomware backup policies.

  The [data_list](#data_list_struct) structure is documented below.

<a name="data_list_struct"></a>
The `data_list` block supports:

* `enabled` - Is the policy enabled. The value can be **true** or **false**.

* `id` - The policy ID.

* `name` - The policy name.

* `operation_type` - The backup type. The value can be **backup**.

* `operation_definition` - The policy attribute retention rules.

  The [operation_definition](operation_definition_struct) structure is documented below.

* `trigger` - The backup policy time scheduling rules.

  The [trigger](#trigger_struct) structure is documented below.

<a name="operation_definition_struct"></a>
The `operation_definition` block supports:

* `day_backups` - The number of daily backups to be retained is not limited by the maximum number of backups to be
  retained. The value ranges from `0` to `100`.

* `max_backups` - The maximum number of backups that can be automatically backed up for a single backup object.
  `-1` represents not cleaning according to the number of backups. If both this field and the `retention_duration_days`
  field are empty, the backup will be permanently retained. The value is `-1` or `1` - `99,999`.

* `month_backups` - The keep monthly backups, which are not limited by the maximum number of backups to be kept.
  If this parameter is selected, the timezone must also be selected. The value ranges from `0` to `100`.

* `retention_duration_days` - The backup retention duration, in days. The maximum support is `99,999` days.
  `-1` represents not cleaning according to time. If both this field and the `max_backups` parameter are empty,
  the backup will be permanently retained. The minimum value is `-1` and the maximum value is `99,999`

* `timezone` - The user's time zone is in UTC+08:00 format.

* `week_backups` - The keep the number of backups per week, which is not limited by the maximum number of backups to be
  kept.

* `year_backups` - The retain the number of backups per year, which is not limited by the maximum number of backups
  retained. The value ranges from `0` to `100`.

<a name="trigger_struct"></a>
The `trigger` block supports:

* `id` - The scheduler ID.

* `name` - The scheduler name.

* `type` - The scheduler type. Currently only supports **time*, Scheduled scheduling.

* `properties` - The scheduler properties.

  The [properties](#properties_struct) structure is documented below.

<a name="properties_struct"></a>
The `properties` block supports:

* `pattern` - The scheduling strategy of scheduler.

* `start_time` - The scheduler start time.
