---
subcategory: "Cloud Search Service (CSS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_css_scan_tasks"
description: |-
  Use this data source to get the list of CSS scan tasks.
---

# huaweicloud_css_scan_tasks

Use this data source to get the list of CSS scan tasks.

## Example Usage

```hcl
variable "cluster_id" {}

data "huaweicloud_css_scan_tasks" "test" {
  cluster_id = var.cluster_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `cluster_id` - (Required, String) Specifies the ID of the CSS cluster.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `scan_tasks` - The scan tasks.

  The [scan_tasks](#scan_tasks_struct) structure is documented below.

<a name="scan_tasks_struct"></a>
The `scan_tasks` block supports:

* `id` - The scan task ID.

* `name` - The scan task name.

* `description` - The scan task description.

* `status` - The execution status of the cluster scan task.

* `smn_status` - The SMN alarm sending status after the cluster scan task is completed.

* `created_at` - The scan task creation time.

* `smn_fail_reason` - The reason for failure in sending SMN alarm.

* `task_risks` - The risk found by the cluster scan task.

  The [task_risks](#scan_tasks_task_risks_struct) structure is documented below.

* `summary` - The risk summary after the cluster scan task is completed.

  The [summary](#scan_tasks_summary_struct) structure is documented below.

<a name="scan_tasks_task_risks_struct"></a>
The `task_risks` block supports:

* `risk` - The risk item.

* `level` - The level of the risk item.

* `description` - The description of the risk item.

* `suggestion` - The suggestion on how to resolve this risk item.

<a name="scan_tasks_summary_struct"></a>
The `summary` block supports:

* `suggestion` - The number of suggestions found by the cluster scan task.

* `high` - The number of high-risk items found by the cluster scan task.

* `medium` - The number of medium-risk items found by the cluster scan task.
