---
subcategory: "SecMaster"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_secmaster_playbook_actions"
description: |-
  Use this data source to get the list of playbook actions.
---

# huaweicloud_secmaster_playbook_actions

Use this data source to get the list of playbook actions.

## Example Usage

```hcl
variable "workspace_id" {}
variable "version_id" {}

data "huaweicloud_secmaster_playbook_actions" "test" {
  workspace_id = var.workspace_id
  version_id   = var.version_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `workspace_id` - (Required, String) Specifies the workspace ID.

* `version_id` - (Required, String) Specifies the playbook version ID.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `data` - The list of the playbook workflows.

  The [data](#data_struct) structure is documented below.

<a name="data_struct"></a>
The `data` block supports:

* `id` - The playbook workflow action ID.

* `name` - The workflow name.

* `description` - The workflow action description.

* `action_id` - The workflow ID.

* `action_type` - The workflow action type.

* `playbook_id` - The playbook ID.

* `playbook_version_id` - The playbook version ID.

* `project_id` - The project ID.
