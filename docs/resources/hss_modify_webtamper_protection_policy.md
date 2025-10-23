---
subcategory: "Host Security Service (HSS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_hss_modify_webtamper_protection_policy"
description: |-
  Manages a resource to update a web tamper protection policy within HuaweiCloud.
---

# huaweicloud_hss_modify_webtamper_protection_policy

Manages a resource to update a web tamper protection policy within HuaweiCloud.

-> This resource is a one-time action resource. Deleting this resource will not clear the corresponding request record,
  but will only remove the resource information from the tf state file.

## Example Usage

### The server OS is Linux

```hcl
variable "host_id" {}
variable "protect_directories" {
  type = list(object({
    protect_dir      = string  
    local_backup_dir = string
  }))
}
variable "week_off_list" {
  type = list(number)
}
variable "time_ranges" {
  type = list(object({
    time_range = string  
  }))
}

resource "huaweicloud_hss_modify_webtamper_protection_policy" "test" {
  host_id = var.host_id

  protect_dir_info {
    dynamic "protect_dir_list" {
      for_each = var.protect_directories
  
      content {
        protect_dir      = protect_dir_list.value.protect_dir
        local_backup_dir = protect_dir_list.value.local_backup_dir
      }
    }

    exclude_file_type = "log;text"
    protect_mode      = "recovery"
  }

  enable_timing_off = true

  timing_off_config_info {
    week_off_list = var.week_off_list

    dynamic "timing_range_list" {
      for_each = var.time_ranges
  
      content {
        time_range = timing_range_list.value
      }
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.
  Changing this parameter will create a new resource.

* `host_id` - (Required, String) Specifies the server ID.
  For this parameter to be valid, the server must be protected by WTP.

* `protect_dir_info` - (Required, List, NonUpdatable) Specifies the protected directory information.
  The [protect_dir_info](#policy_protect_dir_info_struct) structure is documented below.

* `enable_timing_off` - (Optional, Bool, NonUpdatable) Specifies the scheduled protection switch status.
  The valid values are as follows:
  + **true**: Indicates scheduled protection is enabled.
  + **false**: Indicates scheduled protection is disabled. Defaults value.

* `timing_off_config_info` - (Optional, List, NonUpdatable) Specifies the details of scheduled switch configuration.
  The [timing_off_config_info](#policy_timing_off_config_info_struct) structure is documented below.

  -> This parameter is valid and required when `enable_timing_off` is set to **true**.

* `enable_rasp_protect` - (Optional, Bool, NonUpdatable) Specifies whether to enable dynamic web tamper protection.
  The valid values are as follows:
  + **true**: Indicates the dynamic web tamper protection is enabled.
  + **false**: Indicates the dynamic web tamper protection is disabled. Defaults value.

  -> This parameter is only valid for Linux servers.

* `rasp_path` - (Optional, String, NonUpdatable) Specifies the Tomcat bin directory of dynamic WTP.
  + The value must start with a slash (/) and cannot end with a slash (/). Only letters, digits, underscores (_),
    hyphens (-), and periods (.) are allowed.
  + The valid length is `1` to `256` characters.

  -> This parameter is valid and required when `enable_rasp_protect` is set to **true**.

* `enable_privileged_process` - (Optional, Bool, NonUpdatable) Specifies the privileged process status.
  The valid values are as follows:
  + **true**: Indicates the privileged process is enabled.
  + **false**: Indicates the privileged process is disabled. Defaults value.

* `privileged_process_info` - (Optional, List, NonUpdatable) Specifies the privileged process configuration details.
  The [privileged_process_info](#policy_privileged_process_info_struct) structure is documented below.

  -> This parameter is valid and required when `enable_privileged_process` is set to **true**.

* `enterprise_project_id` - (Optional, String, NonUpdatable) Specifies the enterprise project ID.  
  This parameter is valid only when the enterprise project is enabled.
  The default value is **0**, indicating the default enterprise project.
  If it is necessary to operate the hosts under all enterprise projects, the value is **all_granted_eps**.
  If you only have permissions for a specific enterprise project, you need set the enterprise project ID. Otherwise,
  the operation may fail due to insufficient permissions.

<a name="policy_protect_dir_info_struct"></a>
The `protect_dir_info` block supports:

* `protect_dir_list` - (Required, List, NonUpdatable) Specifies the protected directory list.
  The [protect_dir_list](#protect_dir_list_struct) structure is documented below.

  -> At least one, up to `50` items included.

* `exclude_file_type` - (Optional, String, NonUpdatable) Specifies the excluded file type.
  The file type can only contain letters and digits, a maximum of `10` file types are supported, and each file type
  not exceeding `10` characters, multiple file types should be separated by semicolons (;).

* `protect_mode` - (Optional, String, NonUpdatable) Specifies the protection mode.
  The valid values are as follows:
  + **recovery**: Indicates interception mode.
  + **alarm**: Indicates alarm mode. Only Linux server are supported.

<a name="protect_dir_list_struct"></a>
The `protect_dir_list` block supports:

* `protect_dir` - (Required, String, NonUpdatable) Specifies the protected directory.
  + The valid length is `1` to `256` characters.
  + For Linux servers, the value must start with a slash (/) and cannot end with a slash (/). Only letters, digits,
    underscores (_), hyphens (-), and periods (.) are allowed.
  + For Windows servers, the directory name cannot start with a space or end with a backslash (\\), and cannot contain
    the following special characters: ;/*?"<>|.

* `local_backup_dir` - (Optional, String, NonUpdatable) Specifies the local backup path.
  + The local backup path cannot start with a space or end with a slash (/), and cannot contain semicolons (;).
  + The value consists of a maximum of `256` characters.

  -> This parameter is only valid and required for Linux servers.

* `exclude_child_dir` - (Optional, String, NonUpdatable) Specifies the excluded subdirectory.
  + The subdirectory must be a relative path under the protected directory.
  + The maximum length of the subdirectory is `256` characters. A maximum of `10` subdirectories can be added,
    separated by semicolons (;).
  + A subdirectory name on a Linux server cannot start or end with a slash (/). A subdirectory name on a Windows server
    cannot start or end with a backslash (\\).

* `exclude_file_path` - (Optional, String, NonUpdatable) Specifies the excluded file path.
  + The file path must be a relative path under the protected directory.
  + The maximum length of the file path is `256` characters. A maximum of `50` paths can be added,
    separated by semicolons (;).
  + The file path cannot start or end with a slash (/).

  -> This parameter is only valid for Linux servers.

<a name="policy_timing_off_config_info_struct"></a>
The `timing_off_config_info` block supports:

* `week_off_list` - (Required, List, NonUpdatable) Specifies the automatically close the protection cycle.
  The valid value from `1` to `7`, `1` indicates Monday, `2` indicates Tuesday and so on.

* `timing_range_list` - (Required, List, NonUpdatable) Specifies the automatically turn off the protection period.
  The [timing_range_list](#timing_range_list_struct) structure is documented below.

  -> At least one, up to `5` items included.

<a name="timing_range_list_struct"></a>
The `timing_range_list` block supports:

* `time_range` - (Required, String, NonUpdatable) Specifies the time range of automatically turn off the protection
  period.
  + The start time and end time are separated by a hyphen (-). e.g. **15:00-15:30**.
  + A time range must be at least `5` minutes.
  + Time ranges cannot overlap and must have at least a 5-minute interval.

* `description` - (Optional, String, NonUpdatable) Specifies the description of automatically turn off the protection
  period.

<a name="policy_privileged_process_info_struct"></a>
The `privileged_process_info` block supports:

* `privileged_process_path_list` - (Required, List, NonUpdatable) Specifies the list of privileged process file paths.

  -> At least one, up to `10` items included.

* `privileged_child_status` - (Optional, Bool, NonUpdatable) Specifies the privileged sub-process trusted status.
  The valid values are as follows:
  + **true**: Indicates the sub-process is trusted.
  + **false**: Indicates the sub-process is not trusted. Defaults value.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.
