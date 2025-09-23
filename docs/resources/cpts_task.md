---
subcategory: "Cloud Performance Test Service (CPTS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cpts_task"
description: ""
---

# huaweicloud_cpts_task

Manages a pressure Test Task resource within HuaweiCloud CPTS.
The task resource only supports serial execution mode.

## Example Usage

```hcl
resource "huaweicloud_cpts_project" "test" {
  name = "tf_demo_project"
}

resource "huaweicloud_cpts_task" "test" {
  name       = "tf_demo_task"
  project_id = huaweicloud_cpts_project.test.id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the task resource. If omitted, the
  provider-level region will be used. Changing this parameter will create a new resource.

* `name` - (Required, String) Specifies the name of the task, which can contain a maximum of `42` characters.

* `project_id` - (Required, Int, ForceNew) Specifies the CPTS project ID which the task belongs to.
  Changing this parameter will create a new resource.

* `benchmark_concurrency` - (Optional, Int) Specifies benchmark concurrency of the task, the value range is `0` to
  `2,000,000`. The default value is `100`. Reference for the calculation of the number of concurrent users.  
  `Number of concurrent users` = `benchmark concurrency` * `concurrency ratio`.

* `operation` - (Optional, String) Specifies whether to enable the task or stop the task. The options are as follows:
  + **enable**: Starting the pressure test task.
  + **stop**: Stop the pressure test tasks.

 -> Starting the pressure test task Only after add all the test cases to task.

* `cluster_id` - (Optional, Int) Specifies a cluster ID of the CPTS resource group. If the number of concurrent users
  is less than 1000, you can use a shared resource group for testing and do not have to create a resource group.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `status` - The status of the task. The options are as follows:
  + **0**: The task is running.
  + **2**: The task is idle.

## Timeouts

This resource provides the following timeouts configuration options:

* `update` - Default is 30 minutes.

## Import

Tasks can be imported using the `id`, e.g.

```bash
$ terraform import huaweicloud_cpts_task.test 1090
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response, security or some other reason. The missing attributes include: `operation`.
It is generally recommended running `terraform plan` after importing an instance.
You can then decide if changes should be applied to the instance, or the resource definition should be updated to
align with the instance. Also you can ignore changes as below.

```hcl
resource "huaweicloud_cpts_task" "test" {
    ...

  lifecycle {
    ignore_changes = [
      operation,
    ]
  }
}
```
