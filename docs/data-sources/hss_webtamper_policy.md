---
subcategory: "Host Security Service (HSS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_hss_webtamper_policy"
description: |-
  Use this data source to get the web tamper policy of HSS within HuaweiCloud.
---

# huaweicloud_hss_webtamper_policy

Use this data source to get the web tamper policy of HSS within HuaweiCloud.

## Example Usage

```hcl
variable "host_id" {}

data "huaweicloud_hss_webtamper_policy" "test" {
  host_id = var.host_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `host_id` - (Required, String) Specifies the host ID.

* `enterprise_project_id` - (Optional, String, NonUpdatable) Specifies the enterprise project ID.  
  This parameter is valid only when the enterprise project is enabled.
  The default value is **0**, indicating the default enterprise project.
  If it is necessary to operate the asset under all enterprise projects, the value is **all_granted_eps**.
  If you only have permissions for a specific enterprise project, you need set the enterprise project ID. Otherwise,
  the operation may fail due to insufficient permissions.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID in UUID format.

* `protect_dir_num` - The number of protected directories.

* `protect_dir_info` - The protected directory information.

  The [protect_dir_info](#protect_dir_info_struct) structure is documented below.

* `enable_timing_off` - Whether the scheduled shutdown function is enabled.  
  The valid values are as follows:
  + **true**: Activate the timed shutdown protection function.
  + **false**: The timed shutdown protection function is not enabled.

* `timing_off_config_info` - The scheduled shutdown configuration details.

  The [timing_off_config_info](#timing_off_config_info_struct) structure is documented below.

* `enable_rasp_protect` - Whether dynamic web page tamper protection is enabled. Only Linux servers support
  dynamic web tamper protection.  
  The valid values are as follows:
  + **true**: Enable dynamic web tamper protection.
  + **false**: Dynamic web tamper protection is not enabled.

* `rasp_path` - The Tomcat bin directory for dynamic web page tamper protection.

* `enable_privileged_process` - Whether the privileged process is enabled.  
  The valid values are as follows:
  + **true**: Activate the privileged process.
  + **false**: The privileged process has not been enabled.

* `privileged_child_status` - Whether the privileged process child process is trusted. This requires enabling the
  privileged process first.  
  The valid values are as follows:
  + **true**: Enable the trusted sub process of the privileged process.
  + **false**: The privileged process sub process is not enabled and trusted.

* `privileged_process_path_list` - The list of privileged process file paths.

<a name="protect_dir_info_struct"></a>
The `protect_dir_info` block supports:

* `protect_dir_list` - The list of protected directories.

  The [protect_dir_list](#protect_dir_list_struct) structure is documented below.

* `exclude_file_type` - The excluded file types.

* `protect_mode` - The protection mode.  
  + **recovery**: Intercept mode.
  + **alarm**: Alarm mode, only Linux servers support alarm mode.

<a name="protect_dir_list_struct"></a>
The `protect_dir_list` block supports:

* `protect_dir` - The protected directory.

* `exclude_child_dir` - The excluded subdirectories.

* `exclue_file_path` - The excluded file path.

* `local_backup_dir` - The local backup path. Only Linux servers support setting a local backup path.

* `protect_status` - The protection status.  
  The valid values are as follows:
  + **closed**: Not enabled.
  + **opened**: Protection in progress.
  + **opening**: Enabling protection.
  + **closing**: Disabling protection.
  + **open_failed**: Protection failed.

* `error` - The failure reason. This exists when the protection status is **open_failed**.

<a name="timing_off_config_info_struct"></a>
The `timing_off_config_info` block supports:

* `week_off_list` - The automatic shutdown protection cycle list. `1` represents Monday, `2` represents Tuesday,
  `3` represents Wednesday, `4` represents Thursday, `5` represents Friday, `6` represents Saturday,
  `7` represents Sunday.

* `timing_range_list` - The automatic shutdown protection time periods.

  The [timing_range_list](#timing_range_list_struct) structure is documented below.

<a name="timing_range_list_struct"></a>
The `timing_range_list` block supports:

* `time_range` - The automatic shutdown protection time range.

* `description` - The description of the automatic shutdown protection time period.
