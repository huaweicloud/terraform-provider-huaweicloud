---
subcategory: "GaussDB(DWS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dws_workload_plan_execution"
description: |-
  Manages a GaussDB(DWS) workload plan execution resource within HuaweiCloud.
---

# huaweicloud_dws_workload_plan_execution

Manages a GaussDB(DWS) workload plan execution resource within HuaweiCloud.

-> 1. Only one workload plan can be started for each cluster.
  <br/> 2. A workload plan must have at least two plan stages before it can be started.

## Example Usage

```hcl
variable "cluster_id" {}
variable "plan_id" {}
variable "stage_id" {}

resource "huaweicloud_dws_workload_plan_execution" "test" {
  cluster_id = var.cluster_id
  plan_id    = var.plan_id
  stage_id   = var.stage_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `cluster_id` - (Required, String, ForceNew) Specifies the cluster ID of to which the workload plan to execute belongs.
  Changing this parameter will create a new resource.

* `plan_id` - (Required, String, ForceNew) Specifies the ID of the workload plan to be executed.
  Changing this parameter will create a new resource.

* `stage_id` - (Optional, String) Specifies the plan stage ID to be executed after the successful start of the workload
  plan. If omitted, the workload plan will not take effect even if it is successfully started.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID, same as the workload plan ID.
