---
subcategory: "GaussDB(DWS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dws_workload_plans"
description: |-
  Use this data source to query the list of workload plans within HuaweiCloud.
---

# huaweicloud_dws_workload_plans

Using this data source to query the list of workload plans within HuaweiCloud.

## Example Usage

```hcl
variable "cluster_id" {}

data "huaweicloud_dws_workload_plans" "test" {
  cluster_id = var.cluster_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `cluster_id` - (Required, String) Specifies the DWS cluster ID to which the workload plans belong.

* `logical_cluster_name` - (Optional, String) Specifies the logical cluster name to which the workload plans belong.
  This parameter is only available when the DWS cluster is logical cluster.  

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `plans` - All workload plans that match the filter parameters.

  The [plans](#plans_struct) structure is documented below.

<a name="plans_struct"></a>
The `plans` block supports:

* `id` - The ID of the workload plan.

* `name` - The name of the workload plan.

* `cluster_id` - The cluster ID to which the workload plan belongs.

* `current_stage_name` - The name of the current plan stage corresponding to the workload plan.

* `status` - The status of the workload plan.
  + **enabled**: The workload plan has been started.
  + **disabled**: The workload plan has not been started.
