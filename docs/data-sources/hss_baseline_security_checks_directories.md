---
subcategory: "Host Security Service (HSS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_hss_baseline_security_checks_directories"
description: |-
  Use this data source to get the list of HSS baseline security checks directories within HuaweiCloud.
---

# huaweicloud_hss_baseline_security_checks_directories

Use this data source to get the list of HSS baseline security checks directories within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_hss_baseline_security_checks_directories" "test" {
  support_os  = "Linux"
  select_type = "check_type"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `support_os` - (Required, String) Specifies the operating system for baseline checks. Valid values are:
  + **Linux**
  + **Windows**

* `select_type` - (Required, String) Specifies the order of the directory structure. Valid values are:
  + **check_type**: The secondary directory is the baseline name.
  + **tag**: The secondary directory is the type of check item.

* `group_id` - (Optional, String) Specifies the policy group ID to show which check items are selected.

* `enterprise_project_id` - (Optional, String, NonUpdatable) Specifies the enterprise project ID.  
  This parameter is valid only when the enterprise project is enabled.
  The default value is **0**, indicating the default enterprise project.
  If it is necessary to operate the asset under all enterprise projects, the value is **all_granted_eps**.
  If you only have permissions for a specific enterprise project, you need set the enterprise project ID. Otherwise,
  the operation may fail due to insufficient permissions.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID in UUID format.

* `task_condition` - The scheduled detection configuration.

  The [task_condition](#task_condition_struct) structure is documented below.

* `baseline_directory_list` - The baseline check policy directory list.

  The [baseline_directory_list](#baseline_directory_list_struct) structure is documented below.

* `pwd_directory_list` - The baseline check policy directory list for weak passwords.

  The [pwd_directory_list](#pwd_directory_list_struct) structure is documented below.

<a name="task_condition_struct"></a>
The `task_condition` block supports:

* `type` - The scheduled task type. Valid values are:
  + **fixed_weekday**: Fixed weekday.

* `day_of_week` - Which day of the week triggers it. Choose `0` or multiple options.

* `hour` - The hour when the task is triggered.

* `minute` - The minute when the task is triggered.

* `random_offset` - The random offset time.

<a name="baseline_directory_list_struct"></a>
The `baseline_directory_list` block supports:

* `type` - The meaning varies based on the value of `select_type`:
  + When `select_type` is **check_type**, this field represents the check type (baseline name).
  + When `select_type` is **tag**, this field represents the tag (type of baseline check item).

* `standard` - The standard type. Valid values are:
  + **cn_standard**: Compliance standard.
  + **hw_standard**: Cloud security practice standard.
  + **cis_standard**: General security standard.

* `data_list` - The third-level directory information for baseline check policies.

  The [data_list](#data_list_struct) structure is documented below.

<a name="data_list_struct"></a>
The `data_list` block supports:

* `name` - The meaning varies based on the value of `select_type`.
  + When `select_type` is **check_type**, this field represents the tag (type of baseline check item).
  + When `select_type` is **tag**, this field represents the check type (baseline name).

* `enable` - Whether the item is selected.

<a name="pwd_directory_list_struct"></a>
The `pwd_directory_list` block supports:

* `tag` - The primary tag for weak password and password complexity. Valid values are:
  + **weakpwd_pwdcomplexity**: Weak password and password complexity detection.
  + **weakpwd**: Weak password detection.

* `sub_tag` - The sub-tag for password checks. Valid values are:
  + **weak_pwd**: Classic weak password detection.
  + **pwd_complexity**: Password complexity policy check.

* `checked` - Whether the item is selected. The value can be **true** or **false**.

* `key` - The unique value in the directory. Valid values are:
  + **weak_pwd**: Classic weak password detection.
  + **pwd_complexity**: Password complexity policy check.
