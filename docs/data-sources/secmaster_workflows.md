---
subcategory: "SecMaster"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_secmaster_workflows"
description: |-
  Use this data source to get the list of SecMaster workflows.
---

# huaweicloud_secmaster_workflows

Use this data source to get the list of SecMaster workflows.

## Example Usage

```hcl
variable "workspace_id" {}

data "huaweicloud_secmaster_workflows" "test" {
  workspace_id = var.workspace_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `workspace_id` - (Required, String) Specifies the workspace ID.

* `name` - (Optional, String) Specifies the workflow name.

* `type` - (Optional, String) Specifies the workflow type.
  The value can be **NORMAL**, **SURVEY**, **HEMOSTASIS** and **EASE**.

* `description` - (Optional, String) Specifies the workflow description.

* `data_class_id` - (Optional, String) Specifies the data class ID.

* `data_class_name` - (Optional, String) Specifies the data class name.

* `enabled` - (Optional, String) Specifies whether the version is enabled. The value can be **true** and **false**.

* `last_version` - (Optional, String) Specifies whether the version is the latest. The value can be **true** and **false**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `workflows` - The workflow list.

  The [workflows](#workflows_struct) structure is documented below.

<a name="workflows_struct"></a>
The `workflows` block supports:

* `id` - The workflow ID.

* `name` - The workflow name.

* `type` - The workflow type.

* `description` - The workflow description.

* `data_class_id` - The ID of the data class.

* `enabled` - Whether the version is enabled.

* `version_id` - The workflow version ID.

* `engine_type` - The type of engine.

* `current_approval_version_id` - The version to be approved currently.

* `current_rejected_version_id` - The version that has been rejected currently.

* `edit_role` - The edit role.

* `approve_role` - The approval role.

* `use_role` - The user role.

* `workspace_id` - The workspace ID.
