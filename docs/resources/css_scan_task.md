---
subcategory: "Cloud Search Service (CSS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_css_scan_task"
description: ""
---

# huaweicloud_css_scan_task

Manages a CSS cluster scan task resource within HuaweiCloud.

## Example Usage

```hcl
variable "cluster_id" {}
variable "scan_task_name" {}

resource "huaweicloud_css_scan_task" "test" {
  cluster_id = var.cluster_id
  name       = var.scan_task_name
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.
  Changing this creates a new resource.

* `cluster_id` - (Required, String, ForceNew) Specifies ID of the CSS cluster.
  Changing this creates a new resource.

* `name` - (Required, String, ForceNew) Specifies the name of the cluster scan task.
  Changing this creates a new resource.

* `description` - (Optional, String, ForceNew) Specifies the description of the cluster scan task.
  Changing this creates a new resource.

* `alarm` - (Optional, List, ForceNew) Specifies sending SMN alarm message configuration
  after the cluster scan task is completed.
  Changing this creates a new resource.
  The [alarm](#css_scan_task_alarm) structure is documented below.

<a name="css_scan_task_alarm"></a>
The `alarm` block supports:

* `level` - (Required, String) Specifies the level of alarm messages found by the cluster scan task.
  The valid values are **high**, **medium**, **suggestion** and **noRisk**.

* `smn_topic` - (Required, String) Specifies the name of the SMN topic.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `status` - The execution status of the cluster scan task.

* `created_at` - The creation time of the cluster scan task.

* `smn_status` - The SMN alarm sending status after the cluster scan task is completed.

* `smn_fail_reason` - The reason for failure in sending SMN alarm.

* `summary` - The risk summary after the cluster scan task is completed.
  The [summary](#css_scan_task_summary_attr) structure is documented below.

* `task_risks` - The risk found by the cluster scan task.
  The [task_risks](#css_scan_task_risks_attr) structure is documented below.

<a name="css_scan_task_summary_attr"></a>
The `summary` block supports:

* `high_num` - The number of high-risk items found by the cluster scan task.

* `medium_num` - The number of medium-risk items found by the cluster scan task.

* `suggestion_num` - The number of suggestions found by the cluster scan task.

<a name="css_scan_task_risks_attr"></a>
The `task_risks` block supports:

* `risk` - The risk item.

* `level` - The level of the risk item.

* `description` - The description of the risk item.

* `suggestion` - The suggestion on how to resolve this risk item.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 30 minutes.

## Import

The CSS cluster scan task can be imported using `cluster_id` and `name` separated by a slash, e.g.

```bash
$ terraform import huaweicloud_css_scan_task.test <cluster_id>/<name>
```

Note that the imported state may not be identical to your resource definition, due to the attribute missing from the
API response. The missing attribute is: `alarm`.
It is generally recommended running `terraform plan` after importing a scan task.
You can then decide if changes should be applied to the scan task, or the resource definition should be updated to align
with the scan task. Also you can ignore changes as below.

```hcl
resource "huaweicloud_css_scan_task" "test" {
  ...

  lifecycle {
    ignore_changes = [
      alarm,
    ]
  }
}
```
