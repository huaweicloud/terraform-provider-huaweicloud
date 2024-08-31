---
subcategory: "GaussDB(DWS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dws_workload_plan"
description: |-
  Manages a GaussDB(DWS) workload plan resource within HuaweiCloud.
---

# huaweicloud_dws_workload_plan

Manages a GaussDB(DWS) workload plan resource within HuaweiCloud.

## Example Usage

```hcl
variable "cluster_id" {}
variable "plan_name" {}

resource "huaweicloud_dws_workload_plan" "test" {
  cluster_id = var.cluster_id
  name       = var.plan_name
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `cluster_id` - (Required, String, ForceNew) Specifies the cluster ID of to which the workload plan belongs.
  Changing this parameter will create a new resource.

-> Currently, only regular cluster is supported, and logical cluster is temporarily not supported.

* `name` - (Required, String, ForceNew) Specifies the name of the workload plan, which must be unique and contains
  `3` to `28` characters, composed only of lowercase letters, numbers, or underscores (_), and must start with a
  lowercase letter. Changing this parameter will create a new resource.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID in UUID format.

* `status` - The workload plan status. The valid values are as follows:  
  + **enabled**: The workload plan has been started.
  + **disabled**: The workload plan has not been started.

* `current_stage_name` - The name of the current plan stage of the workload plan.

* `stages` - All plan stages under the workload plan.
  The [stages](#DWS_WorkLoadPlan_stages) structure is documented below.

-> The `current_stage_name` and `stages` attributes are available when the workload plan has been started.

<a name="DWS_WorkLoadPlan_stages"></a>
The `stages` block supports:

* `id` - The plan stage ID.

* `name` - The plan stage name.

## Import

The workload plan can be imported using `cluster_id` and `name`, separated by a slash, e.g.

```bash
$ terraform import huaweicloud_dws_workload_plan.test <cluster_id>/<name>
```
