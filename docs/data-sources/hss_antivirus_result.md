---
subcategory: "Host Security Service (HSS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_hss_antivirus_result"
description: |-
  Use this data source to get the list of HSS antivirus results within HuaweiCloud.
---

# huaweicloud_hss_antivirus_result

Use this data source to get the list of HSS antivirus results within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_hss_antivirus_result" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project ID.  
  This parameter is valid only when the enterprise project is enabled.
  The default value is **0**, indicating the default enterprise project.
  If you need to query data for all enterprise projects, the value is **all_granted_eps**.
  If you only have permissions for a specific enterprise project, you need to set the enterprise project ID. Otherwise,
  the operation may fail due to insufficient permissions.

* `host_name` - (Optional, String) Specifies the server name to filter the results.

* `private_ip` - (Optional, String) Specifies the private IP to filter the results.

* `public_ip` - (Optional, String) Specifies the public IP to filter the results.

* `handle_status` - (Optional, String) Specifies the handling status to filter the results.  
  The valid values are as follows:
  + **unhandled**: Unhandled.
  + **handled**: Handled.

* `severity_list` - (Optional, List) Specifies the threat level list to filter the results.  
  The valid values are as follows:
  + **Low**: Low risk.
  + **Medium**: Medium risk.
  + **High**: High risk.
  + **Critical**: Critical risk.

* `asset_value` - (Optional, String) Specifies the asset importance to filter the results.  
  The valid values are as follows:
  + **important**: Important assets.
  + **common**: Common assets.
  + **test**: Test assets.

* `malware_name` - (Optional, String) Specifies the malware name to filter the results.

* `file_path` - (Optional, String) Specifies the file path to filter the results.

* `file_hash` - (Optional, String) Specifies the file hash (SHA256) to filter the results.

* `task_name` - (Optional, String) Specifies the task name to filter the results.

* `manual_isolate` - (Optional, Bool) Specifies whether to use the manual isolation button.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID in UUID format.

* `data_list` - The list of antivirus results.

  The [data_list](#data_list_struct) structure is documented below.

<a name="data_list_struct"></a>
The `data_list` block supports:

* `result_id` - The antivirus result ID.

* `malware_type` - The malware type.

* `malware_name` - The malware name.

* `severity` - The threat level. The valid values are **Low**, **Medium**, **High**, and **Critical**.

* `task_id` - The task ID.

* `task_name` - The task name.

* `file_info` - The file information.

  The [file_info](#data_list_file_info_struct) structure is documented below.

* `resource_info` - The resource information.

  The [resource_info](#data_list_resource_info_struct) structure is documented below.

* `event_type` - The event type.

* `occur_time` - The occurrence time in milliseconds.

* `handle_status` - The handling status. The valid values are **unhandled** and **handled**.

* `handle_method` - The handling method.  
  The valid values are as follows:
  + **mark_as_handled**: Manual handling.
  + **ignore**: Ignore.
  + **add_to_alarm_whitelist**: Add to alarm whitelist.
  + **isolate_and_kill**: Isolate file.

* `memo` - The memo information.

* `operate_accept_list` - The list of supported handling operations.

* `operate_detail_list` - The operation details information list.

  The [operate_detail_list](#data_list_operate_detail_list_struct) structure is documented below.

* `isolate_tag` - The automatic isolation scan tag.

<a name="data_list_file_info_struct"></a>
The `file_info` block supports:

* `file_path` - The file path.

* `file_hash` - The file hash.

* `file_size` - The file size.

* `file_owner` - The file owner.

* `file_attr` - The file attributes.

* `file_ctime` - The file creation time.

* `file_mtime` - The file update time.

<a name="data_list_resource_info_struct"></a>
The `resource_info` block supports:

* `host_name` - The server name.

* `host_id` - The host ID.

* `agent_id` - The agent ID.

* `private_ip` - The private IP address.

* `public_ip` - The public IP address.

* `os_type` - The operating system type. The valid values are **Linux** and **Windows**.

* `host_status` - The server status.  
  The valid values are as follows:
  + **ACTIVE**: Running.
  + **SHUTOFF**: Shutdown.
  + **BUILDING**: Creating.
  + **ERROR**: Fault.

* `agent_status` - The agent status.  
  The valid values are as follows:
  + **installed**: Installed.
  + **not_installed**: Not installed.
  + **online**: Online.
  + **offline**: Offline.
  + **install_failed**: Installation failed.
  + **installing**: Installing.

* `protect_status` - The protection status.  
  The valid values are as follows:
  + **closed**: Not protected.
  + **opened**: Protected.

* `asset_value` - The asset importance. The valid values are **important**, **common**, and **test**.

* `os_name` - The operating system name.

* `os_version` - The operating system version.

<a name="data_list_operate_detail_list_struct"></a>
The `operate_detail_list` block supports:

* `keyword` - The alarm event keywords, only used for alarm whitelist.

* `hash` - The alarm event hash, only used for alarm whitelist.
