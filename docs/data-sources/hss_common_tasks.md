---
subcategory: "Host Security Service (HSS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_hss_common_tasks"
description: |-
  Use this data source to get the list of HSS common tasks within HuaweiCloud.
---

# huaweicloud_hss_common_tasks

Use this data source to get the list of HSS common tasks within HuaweiCloud.

## Example Usage

```hcl
variable "task_type" {}

data "huaweicloud_hss_common_tasks" "test" {
  task_type = var.task_type
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `task_type` - (Required, String) Specifies the task type. Valid values are **cluster_scan** and **iac_scan**.

* `task_id` - (Optional, String) Specifies the task ID to query.

* `task_name` - (Optional, String) Specifies the task name to fuzzy match.

* `start_create_time` - (Optional, Int) Specifies the start time of task creation time range query.

* `end_create_time` - (Optional, Int) Specifies the end time of task creation time range query.

* `trigger_type` - (Optional, String) Specifies the task trigger type. Valid values are **manual** and **schedule**.

* `task_status` - (Optional, String) Specifies the task status. Valid values are **ready**, **running**, and **finished**.

* `sort_key` - (Optional, String) Specifies the sort key. Support **start_time**.

* `sort_dir` - (Optional, String) Specifies the sort direction. Valid values are **desc** and **asc**.

* `cluster_scan_info` - (Optional, List) Specifies the cluster scan information.
  Only valid when `task_type` is **cluster_scan**.
  The [cluster_scan_info](#cluster_scan_info_struct) structure is documented below.

* `iac_scan_info` - (Optional, List) Specifies the IAC scan information.
  Only valid when `task_type` is **iac_scan**.
  The [iac_scan_info](#iac_scan_info_struct) structure is documented below.

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project ID.  
  This parameter is valid only when the enterprise project is enabled.
  The default value is **0**, indicating the default enterprise project.
  If you need to query data for all enterprise projects, the value is **all_granted_eps**.

<a name="cluster_scan_info_struct"></a>
The `cluster_scan_info` block supports:

* `scan_type_list` - (Optional, List) Specifies the list of scan types.
  Valid values are **cluster_vul**, **risk_assessment**, and **benchmark**.

<a name="iac_scan_info_struct"></a>
The `iac_scan_info` block supports:

* `file_type` - (Optional, String) Specifies the file type to scan. Valid values are **dockerfile** and **k8s_yaml**.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID in UUID format.

* `data_list` - The list of tasks.
  The [data_list](#data_list_struct) structure is documented below.

<a name="data_list_struct"></a>
The `data_list` block supports:

* `task_id` - The task ID.

* `task_type` - The task type. Valid values are **cluster_scan** and **iac_scan**.

* `task_name` - The task name.

* `trigger_type` - The task trigger type. Valid values are **manual** and **schedule**.

* `task_status` - The task status. Valid values are **ready**, **running**, and **finished**.

* `start_time` - The start time of the task in milliseconds.

* `end_time` - The end time of the task in milliseconds.

* `estimated_time` - The estimated remaining time in minutes.

* `cluster_scan_info` - The cluster scan information.
  The [cluster_scan_info](#cluster_scan_info_attr_struct) structure is documented below.

* `iac_scan_info` - The IAC scan information.
  The [iac_scan_info](#iac_scan_info_attr_struct) structure is documented below.

<a name="cluster_scan_info_attr_struct"></a>
The `cluster_scan_info` block supports:

* `scan_type_list` - The list of scan types.

* `scanning_cluster_num` - The number of clusters being scanned.

* `success_cluster_num` - The number of successfully scanned clusters.

* `failed_cluster_num` - The number of failed scanned clusters.

<a name="iac_scan_info_attr_struct"></a>
The `iac_scan_info` block supports:

* `file_type` - The file type that was scanned. Valid values are **dockerfile** and **k8s_yaml**.

* `scan_file_num` - The total number of scanned files.

* `success_file_num` - The number of successfully scanned files.

* `failed_file_num` - The number of failed scanned files.
