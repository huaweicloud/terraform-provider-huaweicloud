---
subcategory: "IoT Device Access (IoTDA)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_iotda_batchtask"
description: |-
  Manages an IoTDA batch task resource within HuaweiCloud.
---

# huaweicloud_iotda_batchtask

Manages an IoTDA batch task resource within HuaweiCloud.

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
variable "task_name" {}
variable "task_type" {}
variable "targets_file" {}

resource "huaweicloud_iotda_batchtask" "test" { 
  name         = var.task_name
  type         = var.task_type
  targets_file = var.targets_file
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the IoTDA batch task resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `name` - (Required, String, ForceNew) Specifies the batch task name. The length does not exceed `128`, and only
  combinations of Chinese characters, letters, numbers, and underscores (_) are allowed.
  Changing this parameter will create a new resource.

* `type` - (Required, String, ForceNew) Specifies the batch task type.
  Changing this parameter will create a new resource.  
  The valid values are as follows:
  + **createDevices**: Batch create devices task.
  + **updateDevices**: Batch update devices task.
  + **deleteDevices**: Batch deletion of devices task.
  + **freezeDevices**: Batch freeze devices task.
  + **unfreezeDevices**: Batch unfreeze devices task.

* `space_id` - (Optional, String, ForceNew) Specifies the resource space ID to which the batch task belongs.
  Changing this parameter will create a new resource.

* `targets` - (Optional, List, ForceNew) Specifies an array of target device IDs for executing the batch task, which can
  include up to `30,000` device IDs. This parameter is supported when the `type` is **deleteDevices**,
  **freezeDevices**, or **unfreezeDevices**. Changing this parameter will create a new resource.

* `targets_filter` - (Optional, List, ForceNew) Specifies the batch task target filtering parameters.
  Using this parameter, batch tasks will filter out devices that meet the criteria as targets.
  This parameter is supported when the `type` is **deleteDevices**, **freezeDevices**, or **unfreezeDevices**.
  Changing this parameter will create a new resource.
  The [targets_filter](#IoTDA_targets_filter) structure is documented below.

* `targets_file` - (Optional, String, ForceNew) Specifies the batch task file path to be used for creating batch task.
  Currently, only the **xlsx/xls** file format is supported, and the maximum number of lines in the file is `30000`.
  This parameter is supported when the `type` is **createDevices**, **updateDevices**, **deleteDevices**,
  **freezeDevices**, or **unfreezeDevices**. Changing this parameter will create a new resource.
  Please following [reference](https://support.huaweicloud.com/intl/en-us/usermanual-iothub/iot_01_0032.html),
  download the template file and fill it out.

-> Exactly one of `targets`, `targets_filter`, or `targets_file` should be specified.

<a name="IoTDA_targets_filter"></a>
The `targets_filter` block supports:

* `group_ids` - (Required, List, ForceNew) Specifies the list of device group IDs for executing batch task. Batch task
  will filter out devices within the groups as targets.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `status` - The resource status. The valid values are as follows:
  + **Success**: Batch task execution successful.
  + **Fail**: Batch task execution failed.
  + **PartialSuccess**: Batch task partially executed successfully, some subtasks have been executed successfully,
    while others have failed.

* `status_desc` - The batch task status description, indicating the error message of the main task failure.

* `created_at` - The time of batch task creation. The format is **yyyyMMdd'T'HHmmss'Z**. e.g. **20190528T153000Z**.

* `task_progress` - Subtask execution statistics results.
  The [task_progress](#iotda_task_progress) structure is documented below.

* `task_details` - List of subTask details.
  The [task_details](#iotda_task_details) structure is documented below.

<a name="iotda_task_progress"></a>
The `task_progress` block supports:

* `total` - The total number of subtasks.

* `success` - The number of successfully executed subtasks.

* `fail` - The number of subtasks that failed to execute.

<a name="iotda_task_details"></a>
The `task_details` block supports:

* `target` - The goal of executing subtask. The value includes product ID and node ID.

* `status` - The execution status of subtask. The value can be **Success** or **Fail**.

* `output` - The output information of subtask execution. The value only exists when the subtask is successfully
  executed, including device ID, space ID, device secret, and device fingerprint.

* `error` - Subtask execution failure information. The value only exists when the subtask fails.
  The [task_details](#iotda_task_details_error) structure is documented below.

<a name="iotda_task_details_error"></a>
The `error` block supports:

* `error_code` - Subtask execution failure error code.

* `error_msg` - Subtask execution failure error message.

## Import

The batch task can be imported using the `id`, e.g.

```bash
$ terraform import huaweicloud_iotda_batchtask.test <id>
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response, security or some other reason. The missing attributes include: `space_id`, `targets`, `targets_filter`,
 `targets_file`. It is generally recommended running `terraform plan` after importing a resource.
You can then decide if changes should be applied to the resource, or the resource definition
should be updated to align with the resource. Also, you can ignore changes as below.

```hcl
resource "huaweicloud_iotda_batchtask" "test" { 
  ...
  
  lifecycle {
    ignore_changes = [
      space_id, targets, targets_filter, targets_file,
    ]
  }
}
```
