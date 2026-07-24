---
subcategory: "Data Security Center (DSC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dsc_db_mask_tasks"
description: |-
  Use this data source to get the list of DSC database mask tasks within HuaweiCloud.
---

# huaweicloud_dsc_db_mask_tasks

Use this data source to get the list of DSC database mask tasks within HuaweiCloud.

## Example Usage

```hcl
variable "template_id" {}

data "huaweicloud_dsc_db_mask_tasks" "test" {
  template_id = var.template_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `template_id` - (Required, String) Specifies the template ID.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `tasks` - The database mask task list.

  The [tasks](#tasks_struct) structure is documented below.

<a name="tasks_struct"></a>
The `tasks` block supports:

* `db_type` - The database type.

* `end_time` - The task end time.

* `execute_line` - The number of executed lines.

* `id` - The task ID.

* `progress` - The task progress.

* `run_status` - The task running status.

* `start_time` - The task start time.

* `task_template_id` - The task template ID.

* `type` - The task type.
