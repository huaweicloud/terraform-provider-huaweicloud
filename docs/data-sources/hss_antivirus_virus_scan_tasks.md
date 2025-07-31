---
subcategory: "Host Security Service (HSS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_hss_antivirus_virus_scan_tasks"
description: |-
  Use this data source to get the list of HSS virus scan tasks within HuaweiCloud.
---
# huaweicloud_hss_antivirus_virus_scan_tasks

Use this data source to get the list of HSS virus scan tasks within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_hss_antivirus_virus_scan_tasks" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `whether_paid_task` - (Required, Bool) Specifies whether this scanning task is paid or not.  
  The valid values are as follows:
  + **true**
  + **false**

* `task_name` - (Optional, String) Specifies the task name.

* `last_days` - (Optional, Int) Specifies the number of days within the query time range. Which is mutually exclusive
  with parameters `begin_time` and `end_time`.

* `begin_time` - (Optional, String) Specifies the starting time of the query time period, with a millisecond-level
  timestamp. It is mutually exclusive with parameter `last_days`, and it needs to satisfy the condition that the
  difference between `end_time` and `begin_time` is less than or equal to `2` days.

* `end_time` - (Optional, String) Specifies the end time of the query time period, with a millisecond-level
  timestamp. It is mutually exclusive with parameter `last_days`, and it needs to satisfy the condition that the
  difference between `end_time` and `begin_time` is less than or equal to `2` days.

* `task_status` - (Optional, String) Specifies the task status.  
  The valid values are as follows:
  + **scanning**
  + **cancel**
  + **fail**
  + **finish**

* `host_name` - (Optional, String) Specifies the host name.

* `private_ip` - (Optional, String) Specifies the host private IP address.

* `public_ip` - (Optional, String) Specifies the host public IP address.

* `host_task_status` - (Optional, List) Specifies the list of host scanning status.  
  The valid values are as follows:
  + **scanning**
  + **success**
  + **fail**
  + **cancel**

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project ID.  
  This parameter is valid only when the enterprise project is enabled.
  The default value is **0**, indicating the default enterprise project.
  If you need to query data for all enterprise projects, the value is **all_granted_eps**.
  If you only have permissions for a specific enterprise project, you need set the enterprise project ID. Otherwise,
  the operation may fail due to insufficient permissions.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID in UUID format.

* `total_num` - The total number of tasks.

* `data_list` - The list of virus scan tasks details.
  The [data_list](#data_list_struct) structure is documented below.

<a name="data_list_struct"></a>
The `data_list` block supports:

* `task_id` - The task ID.

* `task_name` - The task name.

* `scan_type` - The scan type.  
  The valid values are as follows:
  + **quick**
  + **full**
  + **custom**

* `start_type` - The startup type.  
  The valid values are as follows:
  + **now**
  + **later**
  + **period**

* `action` - The disposal action.  
  The valid values are as follows:
  + **auto**
  + **manual**

* `start_time` - The startup time in milliseconds.

* `task_status` - The task status.  
  The valid values are as follows:
  + **scanning**
  + **finish**

* `host_num` - The number of associated hosts.

* `success_host_num` - The number of hosts successfully scanned.

* `fail_host_num` - The number of hosts with scanning failures.

* `cancel_host_num` - The number of hosts cancelled.

* `host_info_list` - The list of host details.
  The [host_info_list](#host_info_list_struct) structure is documented below.

* `rescan` - Do you need to rescan.

* `whether_paid_task` - Is this scanning task paid or not.

<a name="host_info_list_struct"></a>
The `host_info_list` block supports:

* `host_id` - The host ID.

* `host_name` - The host name.

* `private_ip` - The host private IP address.

* `public_ip` - The host public IP address.

* `asset_value` - The importance of assets.  
  The valid values are as follows:
  + **important**: Deleted.
  + **common*: Not deleted.
  + **test*: Not deleted.

* `start_time` - The start time in milliseconds.

* `run_duration` - The running time in seconds.

* `scan_progress` - The scan progress.

* `virus_num` - The number of hosts cancelled.

* `scan_file_num` - The number of scanned files.

* `host_task_status` - The host scanning status.

* `fail_reason` - The failure reason.

* `deleted` - Do you want to delete it.  
  The valid values are as follows:
  + **true**: Deleted.
  + **false*: Not deleted.

* `whether_using_quota` - Whether to use virus scanning and removal quota on a per-use basis.

* `agent_id` - The agent ID.

* `os_type` - The operating system type.
  The valid values are as follows:
  + **Linux**
  + **Windows*

* `host_status` - The host status.

* `agent_status` - The agent status.  
  The valid values are as follows:
  + **installed**
  + **not_installed*
  + **online**
  + **offline**
  + **install_failed**
  + **installing**
  + **not_online**

* `protect_status` - The protection status.  
  The valid values are as follows:
  + **closed**
  + **opened**

* `os_name` - The operating system name.

* `os_version` - The operating system version.
