---
subcategory: "SecMaster"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_secmaster_workflow_instances"
description: |-
  Use this data source to get the list of SecMaster workflow instances.
---

# huaweicloud_secmaster_workflow_instances

Use this data source to get the list of SecMaster workflow instances.

## Example Usage

```hcl
variable "workspace_id" {}

data "huaweicloud_secmaster_workflow_instances" "test" {
  workspace_id = var.workspace_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `workspace_id` - (Required, String) Specifies the workspace ID.

* `workflow_id` - (Optional, String) Specifies the workflow ID.

* `id` - (Optional, String) Specifies the workflow instance ID.

* `name` - (Optional, String) Specifies the workflow instance name.

* `dataclass_id` - (Optional, String) Specifies the data class ID.

* `playbook_id` - (Optional, String) Specifies the playbook ID.

* `defence_id` - (Optional, String) Specifies the defence ID.

* `status` - (Optional, String) Specifies the workflow status.
  The valid values are as follows:
  + **CREATED**
  + **RUNNING**
  + **FINISHED**
  + **RETRYING**
  + **TERMINATING**
  + **TERMINATED**
  + **FAILED**

* `trigger_type` - (Optional, String) Specifies the workflow trigger mode.
  The valid values are as follows:
  + **TIMER**: Indicates scheduled triggering.
  + **EVENT**: Indicates event triggering.

* `from_date` - (Optional, String) Specifies the search start time.
  The time is RFC3339 format. e.g. **2024-08-27T11:00:00.000Z+0800**.

* `to_date` - (Optional, String) Specifies the search end time.
  The time is RFC3339 format. e.g. **2024-08-27T11:00:00.000Z+0800**.

* `sort_key` - (Optional, String) Specifies sorting field.
 The value can be **start_time** or **end_time**.

* `sort_dir` - (Optional, String) Specifies sorting order.
  The valid values are as follows:
  + **ASC**: Ascending order.
  + **DESC**: Descending order.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `instance` - The list of the workflow instances.

  The [instance](#instances_struct) structure is documented below.

<a name="instances_struct"></a>
The `instances` block supports:

* `id` - The workflow instance ID.

* `name` - The workflow instance name.

* `workflow` - The workflow information of the instance.

  The [workflow](#instances_workflow_struct) structure is documented below.

* `dataclass` - The data class of the instance.

  The [dataclass](#instances_dataclass_struct) structure is documented below.

* `playbook` - The playbook information of the instance.

  The [playbook](#instances_playbook_struct) structure is documented below.

* `trigger_type` - The workflow trigger mode.

* `status` - The workflow instance status.

* `start_time` - The start time.

* `end_time` - The end time.

* `retry_count` - The workflow instance retry count.

* `defense_id` - The defense ID.

* `dataobject_id` - The data object ID.

<a name="instances_workflow_struct"></a>
The `workflow` block supports:

* `id` - The workflow ID.

* `name` - The workflow name.

* `name_en` - The workflow English name.

* `version` - The workflow version.

<a name="instances_dataclass_struct"></a>
The `dataclass` block supports:

* `id` - The data class ID.

* `name` - The data class name.

* `en_name` - The data class English name.

<a name="instances_playbook_struct"></a>
The `playbook` block supports:

* `id` - The playbook ID.

* `name` - The playbook name.

* `en_name` - The playbook English name.
