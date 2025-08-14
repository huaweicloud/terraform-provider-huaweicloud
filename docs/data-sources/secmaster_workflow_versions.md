---
subcategory: "SecMaster"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_secmaster_workflow_versions"
description: |-
  Use this data source to get the list of workflow versions under a specified workflow.
---

# huaweicloud_secmaster_workflow_versions

Use this data source to get the list of workflow versions under a specified workflow.

## Example Usage

```hcl
variable "workspace_id" {}
variable "workflow_id" {}

data "huaweicloud_secmaster_workflow_versions" "test" {
  workspace_id = var.workspace_id
  workflow_id  = var.workflow_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `workspace_id` - (Required, String) Specifies the workspace ID.

* `workflow_id` - (Required, String) Specifies the workflow ID.

* `status` - (Optional, String) Specifies the workflow version status.
  The value can be **pending_submit**, **pending_approval**, **not_activated**, **activated** or **rejected**.

activated,pending_approval,not_activated，pending_submit，rejected

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `data` - The list of workflow versions.

  The [data](#data_struct) structure is documented below.

<a name="data_struct"></a>
The `data` block supports:

* `id` - The workflow version ID.

* `name` - The workflow name.

* `description` - The workflow version description.

* `vworkflow_id` - The workflow ID.

* `project_id` - The project ID.

* `owner_id` - The workflow owner ID.

* `creator_id` - The workflow version creator.

* `enabled` - Whether the workflow version enabled.

* `status` - The workflow version status.

* `version` - The workflow version number.

* `taskconfig` - The workflow topology diagram parameter configuration.

* `taskflow` - The workflow topology diagram Base64 encoding.

* `taskflow_type` - The workflow topology diagram type.

* `aop_type` - The workflow type.
  The value can be **NORMAL**, **SURVEY**, **HEMOSTASIS** or **EASE**.
