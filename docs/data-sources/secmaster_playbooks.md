---
subcategory: "SecMaster"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_secmaster_playbooks"
description: |-
  Use this data source to get the list of SecMaster playbooks.
---

# huaweicloud_secmaster_playbooks

Use this data source to get the list of SecMaster playbooks.

## Example Usage

```hcl
variable "name" {}
variable "workspace_id" {}

data "huaweicloud_secmaster_playbooks" "test" {
  workspace_id = var.workspace_id
  name         = var.name
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `workspace_id` - (Required, String) Specifies the workspace ID.

* `name` - (Optional, String) Specifies the playbook name.

* `enabled` - (Optional, String) Specifies whether the playbook is enabled. The value can be **true** or **false**.

* `description` - (Optional, String) Specifies the playbook description.

* `data_class_name` - (Optional, String) Specifies the data class name.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `playbooks` - The playbook list.

  The [playbooks](#playbooks_struct) structure is documented below.

<a name="playbooks_struct"></a>
The `playbooks` block supports:

* `id` - The playbook ID.

* `name` - The playbook name.

* `description` - The playbook description.

* `enabled` - Whether the playbook is enabled.

* `workspace_id` - The workspace ID to which the playbook belongs.

* `owner_id` - The owner ID.

* `project_id` - The project ID.

* `version` - The version.

* `version_id` - The playbook version ID.

* `unaudited_version_id` - The ID of the playbook version to be reviewed.

* `reject_version_id` - The ID of the rejected playbook version.

* `user_role` - The user role of the playbook.

* `edit_role` - The edit role of the playbook.

* `data_class_name` - The data class name.

* `created_at` - The playbook creation time.

* `updated_at` - The playbook update time.

* `approve_role` - The approval role of the playbook.

* `data_class_id` - The data class ID.
