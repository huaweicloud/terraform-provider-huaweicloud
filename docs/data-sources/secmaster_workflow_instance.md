---
subcategory: "SecMaster"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_secmaster_workflow_instance"
description: |-
  Use this data source to query a specific workflow instance detail.
---

# huaweicloud_secmaster_workflow_instances

Use this data source to query a specific workflow instance detail.

## Example Usage

```hcl
variable "workspace_id" {}
variable "instance_id" {}

data "huaweicloud_secmaster_workflow_instance" "test" {
  workspace_id = var.workspace_id
  instance_id  = var.instance_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `workspace_id` - (Required, String) Specifies the workspace ID.

* `instance_id` - (Required, String) Specifies the workflow instance ID.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `workflow_instance_id` - The workflow instance ID.

* `name` - The workflow instance name.

* `workflow` - The workflow information of the instance.

  The [workflow](#workflow_struct) structure is documented below.

* `dataclass` - The data class information of the instance.

  The [dataclass](#dataclass_struct) structure is documented below.

* `playbook` - The playbook information of the instance.

  The [playbook](#playbook_struct) structure is documented below.

* `trigger_type` - The workflow trigger mode.

* `status` - The workflow instance status.

* `start_time` - The start time.

* `end_time` - The end time.

* `retry_count` - The workflow instance retry count.

* `defense_id` - The defense ID.

* `dataobject_id` - The data object ID.

<a name="workflow_struct"></a>
The `workflow` block supports:

* `id` - The workflow ID.

* `name` - The workflow Chinese name.

* `name_en` - The workflow English name.

* `version` - The workflow version.

<a name="dataclass_struct"></a>
The `dataclass` block supports:

* `id` - The data class ID.

* `name` - The data class Chinese name.

* `en_name` - The data class English name.

<a name="playbook_struct"></a>
The `playbook` block supports:

* `id` - The playbook ID.

* `name` - The playbook name.

* `en_name` - The playbook English name.
