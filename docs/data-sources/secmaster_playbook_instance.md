---
subcategory: "SecMaster"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_secmaster_playbook_instance"
description: |-
  Use this data source to get the details of a specific SecMaster playbook instance within HuaweiCloud.
---

# huaweicloud_secmaster_playbook_instance

Use this data source to get the details of a specific SecMaster playbook instance within HuaweiCloud.

## Example Usage

```hcl
variable "workspace_id" {}
variable "instance_id" {}

data "huaweicloud_secmaster_playbook_instance" "test" {
  workspace_id = var.workspace_id
  instance_id  = var.instance_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `workspace_id` - (Required, String) Specifies the workspace ID.

* `instance_id` - (Required, String) Specifies the instance ID.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source, also known as the playbook instance ID.

* `name` - The playbook instance name.

* `project_id` - The project ID.

* `playbook` - The playbook information.

  The [playbook](#playbook_struct) structure is documented below.

* `dataclass` - The data class information.

  The [dataclass](#dataclass_struct) structure is documented below.

* `dataobject` - The data object key field information.

  The [dataobject](#dataobject_struct) structure is documented below.

* `status` - The playbook instance status. Valid values are **RUNNING**, **FINISHED**, **FAILED**, **RETRYING**,
  **TERMINATING**, **TERMINATED**.

* `trigger_type` - The trigger type.  
  The valid values are as follows:
  + **TIMER**: Timed trigger.
  + **EVENT**: Event trigger.

* `start_time` - The creation time.

* `end_time` - The update time.

<a name="playbook_struct"></a>
The `playbook` block supports:

* `id` - The playbook ID.

* `version_id` - The playbook version ID.

* `name` - The name.

* `version` - The version.

<a name="dataclass_struct"></a>
The `dataclass` block supports:

* `id` - The data class ID.

* `name` - The data class name.

<a name="dataobject_struct"></a>
The `dataobject` block supports:

* `id` - The unique identifier ID.

* `name` - The name.
