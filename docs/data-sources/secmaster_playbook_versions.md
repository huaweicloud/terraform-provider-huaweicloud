---
subcategory: "SecMaster"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_secmaster_playbook_versions"
description: |-
  Use this data source to get the list of SecMaster playbook versions.
---

# huaweicloud_secmaster_playbook_versions

Use this data source to get the list of SecMaster playbook versions.

## Example Usage

```hcl
variable "workspace_id" {}
variable "playbook_id" {}

data "huaweicloud_secmaster_playbook_versions" "test" {
  workspace_id = var.workspace_id
  playbook_id  = var.playbook_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `workspace_id` - (Required, String) Specifies the workspace ID.

* `playbook_id` - (Required, String) Specifies the playbook ID.

* `status` - (Optional, String) Specifies the playbook version status.
  The value can be **EDITING**, **APPROVING**, **UNPASSED** or **PUBLISHED**.

* `type` - (Optional, String) Specifies the version type.
  The options are as follows:
  + **0**: indicates draft version;
  + **1**: indicates official version.

* `enabled` - (Optional, String) Specifies whether this version is activated.
  The options are as follows:
  + **0**: indicates false;
  + **1**: indicates true.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `playbook_versions` - The playbook version list.

  The [playbook_versions](#playbook_versions_struct) structure is documented below.

<a name="playbook_versions_struct"></a>
The `playbook_versions` block supports:

* `id` - The playbook version ID.

* `status` - The playbook version status.

* `version` - The playbook version.

* `type` - The playbook version type.

* `enabled` - Whether the playbook version is activated.

* `description` - The description.

* `created_at` - The creation time.

* `dataobject_create` - Whether to trigger a playbook when a data object is created.

* `data_class_id` - The data class ID.

* `playbook_id` - The playbook ID.

* `trigger_type` - The triggering type.

* `modifier_id` - The ID of the user who updated the information.

* `project_id` - The project ID.

* `rule_enabled` - Whether the filtering rule is enabled.

* `data_object_delete` - Whether to trigger a playbook when a data object is deleted.

* `data_object_update` - Whether to trigger a playbook when a data object is updated.

* `rule_id` - The rule ID.

* `data_class_name` - The data class name.

* `updated_at` - The update time.

* `creator_id` - The creator ID.

* `action_strategy` - The execution policy.

* `approve_name` - The reviewer.
