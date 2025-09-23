---
subcategory: "IoT Device Access (IoTDA)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_iotda_batchtasks"
description: |-
  Use this data source to get the list of IoTDA batch tasks within HuaweiCloud.
---

# huaweicloud_iotda_batchtasks

Use this data source to get the list of IoTDA batch tasks within HuaweiCloud.

-> When accessing an IoTDA **standard** or **enterprise** edition instance, you need to specify the IoTDA service
  endpoint in `provider` block.
  You can login to the IoTDA console, choose the instance **Overview** and click **Access Details**
  to view the HTTPS application access address. An example of the access address might be
  **9bc34xxxxx.st1.iotda-app.ap-southeast-1.myhuaweicloud.com**, then you need to configure the
  `provider` block as follows:

  ```hcl
  provider "huaweicloud" {
    endpoints = {
      iotda = "https://9bc34xxxxx.st1.iotda-app.ap-southeast-1.myhuaweicloud.com"
    }
  }
  ```

## Example Usage

```hcl
variable "type" {}

data "huaweicloud_iotda_batchtasks" "test" {
  type = var.type
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the batch tasks.
  If omitted, the provider-level region will be used.

* `type` - (Required, String) Specifies the type of the batch task.
  The valid values are as follows:
  + **createDevices**: Batch create devices task.
  + **updateDevices**: Batch update devices task.
  + **deleteDevices**: Batch deletion of devices task.
  + **freezeDevices**: Batch freeze devices task.
  + **unfreezeDevices**: Batch unfreeze devices task.

* `space_id` - (Optional, String) Specifies the space ID.
  If omitted, query all batch tasks under the current instance.

* `status` - (Optional, String) Specifies the status of the batch task.
  The valid values are as follows:
  + **Initializing**
  + **Waitting**
  + **Processing**
  + **Success**
  + **Fail**
  + **PartialSuccess**
  + **Stopped**
  + **Stopping**

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID in UUID format.

* `batchtasks` - All batch tasks that match the filter parameters.
  The [batchtasks](#iotda_batchtasks) structure is documented below.

<a name="iotda_batchtasks"></a>
The `batchtasks` block supports:

* `id` - The batch task ID.

* `name` - The batch task name.

* `type` - The batch task type.

* `targets` - The target device ID array for executing the batch task.

* `targets_filter` - The batch task target filtering parameters. Json format, which contains key value pairs and (K, V)
  format to identify the parameters required for filtering targets. Currently, the supported K format includes
  group_ids, and the task will select devices that meet the conditions of the group as targets.

* `task_policy` - The task execution strategy.
  The [task_policy](#iotda_task_policy) structure is documented below.

* `status` - The status of the batch task.

* `status_desc` - The batch task status description, including main task failure error information.

* `task_progress` - The subtask execution statistics results.
  The [task_progress](#iotda_task_progress) structure is documented below.

* `created_at` - The creation time of the batch task. The format is **yyyyMMdd'T'HHmmss'Z**. e.g. **20190528T153000Z**.

<a name="iotda_task_policy"></a>
The `task_policy` block supports:

* `schedule_time` - The batch task specifies execution time.
  The format is **yyyyMMdd'T'HHmmss'Z**. e.g. **20190528T153000Z**.
  The valid value is within `7` days. If it is empty, it means task will be executed immediately.

* `retry_count` - The automatic retry times for subtasks of batch tasks.
  The valid value is range form `1` to `5`.

* `retry_interval` - The time interval for automatic retry after a subtask of a batch task fails. Unit in minutes,
  the valid value is range form `0` to `1,440`, the `0` means no retry.

<a name="iotda_task_progress"></a>
The `task_progress` block supports:

* `total` - The total number of subtasks.

* `processing` - The number of subtasks currently being executed.

* `success` - The number of successfully executed subtasks.

* `fail` - The number of subtasks that failed to execute.

* `waitting` - The number of subtasks waiting to be executed.

* `fail_wait_retry` - The number of subtasks waiting for retry due to failure.

* `stopped` - The number of stopped subtasks.

* `removed` - The number of subtasks removed.
