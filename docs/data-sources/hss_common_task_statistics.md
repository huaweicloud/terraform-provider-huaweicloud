---
subcategory: "Host Security Service (HSS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_hss_common_task_statistics"
description: |-
  Use this data source to get the task statistics.
---

# huaweicloud_hss_common_task_statistics

Use this data source to get the task statistics.

## Example Usage

```hcl
data "huaweicloud_hss_common_task_statistics" "test" {
  task_type = "cluster_scan"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `task_type` - (Required, String) Specifies the asset type.
  The valid values are as follows:
  + **cluster_scan**: Indicates cluster scanning task.
  + **iac_scan**: Indicates iac scanning task.

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project ID.
  This parameter is valid only when the enterprise project is enabled.
  The default value is **0**, indicating the default enterprise project.
  If you need to query data for all enterprise projects, the value is **all_granted_eps**.
  If you only have permissions for a specific enterprise project, you need set the enterprise project ID. Otherwise,
  the operation may fail due to insufficient permissions.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `total_num` - The total number of tasks executed cumulatively.

* `running_num` - The number of tasks in the scan.

* `last_task_start_time` - The creation time of the most recent scanning task.
