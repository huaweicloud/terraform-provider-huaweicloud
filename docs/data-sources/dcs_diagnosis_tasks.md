---
subcategory: "Distributed Cache Service (DCS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dcs_diagnosis_tasks"
description: |-
  Use this data source to get the list of DCS diagnosis tasks.
---

# huaweicloud_dcs_diagnosis_tasks

Use this data source to get the list of DCS diagnosis tasks.

## Example Usage

```hcl
variable "instance_id" {}

data "huaweicloud_dcs_diagnosis_tasks" "test" {
  instance_id = var.instance_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the ID of the DCS instance.

* `task_id` - (Optional, String) Specifies the ID of the diagnosis task.

* `status` - (Optional, String) Specifies the status of the diagnosis task.
  Value options: **diagnosing**, **finished**.

* `begin_time` - (Optional, String) Specifies the start time of the diagnosis task, in RFC3339 format.

* `end_time` - (Optional, String) Specifies the end time of the diagnosis task, in RFC3339 format.

* `node_num` - (Optional, String) Specifies the number of diagnosed nodes.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `diagnosis_tasks` - Indicates the list of diagnosis reports.

  The [diagnosis_tasks](#diagnosis_tasks_struct) structure is documented below.

<a name="diagnosis_tasks_struct"></a>
The `diagnosis_tasks` block supports:

* `id` - Indicates the diagnosis task ID.

* `status` - Indicates the diagnosis task status.

* `begin_time` - Indicates the start time of the diagnosis task, in RFC3339 format.

* `end_time` - Indicates the end time of the diagnosis task, in RFC3339 format.

* `created_at` - Indicates the time when the diagnosis report is created.

* `node_num` - Indicates the number of diagnosed nodes.

* `abnormal_item_sum` - Indicates the total number of abnormal diagnosis items.

* `failed_item_sum` - Indicates the total number of failed diagnosis items.
