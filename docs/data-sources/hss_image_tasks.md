---
subcategory: "Host Security Service (HSS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_hss_image_tasks"
description: |-
  Use this data source to get the list of image scan tasks.
---

# huaweicloud_hss_image_tasks

Use this data source to get the list of image scan tasks.

## Example Usage

```hcl
variable "type" {}

data "huaweicloud_hss_image_tasks" "test" {
  type = var.type
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `type` - (Required, String) Specifies the task type.
  The valid values are as follows:
  + **image_sync**: Image synchronization.
  + **image_scan**: Image scan.

* `global_image_type` - (Optional, String) Specifies the image type.
  The valid values are as follows:
  + **local**
  + **registry**

* `task_type` - (Optional, String) Specifies the task type.
  The valid values are as follows:
  + **cycle**: Scheduled scan.
  + **manual**: Manual scan.
  + **autoSync**: Scheduled synchronization.
  + **manualSync**: Manual synchronization.

* `task_name` - (Optional, String) Specifies the task name fuzzy match.

* `task_id` - (Optional, String) Specifies the task ID.

* `create_time` - (Optional, Int) Specifies the task creation time, in milliseconds.

* `end_time` - (Optional, Int) Specifies the task end time, in milliseconds.

* `task_status` - (Optional, String) Specifies the task status.
  The valid values are as follows:
  + **scanning**
  + **finished**

* `scan_scope` - (Optional, String) Specifies the scan risk type.
  The valid values are as follows:
  + **0**: None.
  + **0x7fffffff**: All.
  + **0x000f0000**: Vulnerability.
  + **0x0000f000**: Baseline check.
  + **0x00000f00**: Malicious file.
  + **0x000000f0**: Sensitive information.
  + **0x0000000f**: Software compliance.

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project ID.
  This parameter is valid only when the enterprise project is enabled.
  The default value is **0**, indicating the default enterprise project.
  If you need to query data for all enterprise projects, the value is **all_granted_eps**.
  If you only have permissions for a specific enterprise project, you need set the enterprise project ID. Otherwise,
  the operation may fail due to insufficient permissions.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `data_list` - The list of image scan tasks.

  The [data_list](#data_list_struct) structure is documented below.

<a name="data_list_struct"></a>
The `data_list` block supports:

* `task_id` - The task ID.

* `policy_id` - The policy ID.

* `task_name` - The task name.

* `begin_time` - The task start time.

* `end_time` - The task end time.

* `remain_min` - The task remain time.

* `task_type` - The task type.

* `image_type` - The image type.

* `task_status` - The task status.
  The valid values are as follows:
  + **100**: The scan is completed.
  + **0-100**: Scan progress.
  + **-1**: The scan is terminated.
  + **-2**: The scan times out.
  + **-3**: The scan is abnormal.

* `scan_scope` - The scan risk type.

* `rate_limit` - The scan speed limit (unit:images/hour)

* `is_all` - Whether to scan all images.

* `failed_num` - The number of images that fail to be scanned.

* `success_num` - The number of images that are successfully scanned.

* `total_num` - The total number of images associated with a task.

* `risky_num` - The total number of images with vulnerability risks, baseline risks, and malicious files.

* `sync_task_type` - The synchronization task type.

* `failed_reason` - The failed reason.

* `failed_images` - The list of failed images.

  The [failed_images](#data_list_images_struct) structure is documented below.

<a name="data_list_images_struct"></a>
The `failed_images` block supports:

* `id` - The failed image ID.

* `registry_id` - The image repository ID.

* `registry_name` - The image repository name.

* `image_name` - The image name.

* `image_version` - The image version.

* `namespace` - The name space.

* `registry_type` - The image repository type.
  The valid values are as follows:
  + **SwrPrivate**: SWR private repository.
  + **SwrShared**: SWR shared repository.
  + **SwrEnterprise**: SWR enterprise repository.
  + **Harbor**: Harbor repository.
  + **Jfrog**: JFfog repository.
  + **Other**: Other repository.

* `failed_reason` - The failed reason.
