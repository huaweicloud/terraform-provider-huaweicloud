---
subcategory: "SecMaster"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_secmaster_playbook_instances"
description: |-
  Use this data source to get the list of SecMaster playbook instances.
---

# huaweicloud_secmaster_playbook_instances

Use this data source to get the list of SecMaster playbook instances.

## Example Usage

```hcl
variable "workspace_id" {}

data "huaweicloud_secmaster_playbook_instances" "test" {
  workspace_id = var.workspace_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `workspace_id` - (Required, String) Specifies the workspace ID.

* `from_date` - (Optional, String) Specifies the search start time.
  The supported time formats are as follows:
  + **2023-04-27T13:00:00.000+08:00**
  + **2024-08-27T11:00:00.000Z+0800**
  + **2024-08-27T11:00:00Z+0800**

* `to_date` - (Optional, String) Specifies the search end time.
  The supported time formats are as follows:
  + **2023-04-27T13:00:00.000+08:00**
  + **2024-08-27T11:00:00.000Z+0800**
  + **2024-08-27T11:00:00Z+0800**

* `status` - (Optional, String) Specifies the playbook instance status.
  The value can be **RUNNING**, **FINISHED**, **FAILED**, **TERMINATING** or **TERMINATED**.

* `data_class_name` - (Optional, String) Specifies the data class name.

* `data_object_name` - (Optional, String) Specifies the data object name.

* `trigger_type` - (Optional, String) Specifies the triggering type.
  + **TIMER**: indicates scheduled triggering,
  + **EVENT**: indicates event triggering.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `instances` - The playbook instance list.

  The [instances](#instances_struct) structure is documented below.

<a name="instances_struct"></a>
The `instances` block supports:

* `id` - The playbook instance ID.

* `name` - The playbook instance name.

* `status` - The playbook instance status.

* `trigger_type` - The triggering type.

* `created_at` - The creation time.

* `updated_at` - The update time.

* `project_id` - The project ID.

* `playbook` - The playbook information of the instance.

  The [playbook](#instances_playbook_struct) structure is documented below.

* `data_object` - The data object of the instance.

  The [data_object](#instances_data_object_struct) structure is documented below.

* `data_class` - The data class of the instance.

  The [data_class](#instances_data_class_struct) structure is documented below.

<a name="instances_playbook_struct"></a>
The `playbook` block supports:

* `id` - The playbook ID.

* `name` - The playbook name.

* `version` - The playbook version.

* `version_id` - The playbook version ID.

<a name="instances_data_object_struct"></a>
The `data_object` block supports:

* `id` - The data object ID of the instance.

* `name` - The data object name of the instance.

* `created_at` - The creation time of the data object.

* `updated_at` - The update time of the data object.

* `project_id` - The project ID of the data object.

* `data_class_id` - The data class ID of the data object.

* `content` - The data content of the data object.

<a name="instances_data_class_struct"></a>
The `data_class` block supports:

* `id` - The data class ID of the instance.

* `name` - The data class name of the instance.
