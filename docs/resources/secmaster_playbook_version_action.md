---
subcategory: "SecMaster"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_secmaster_playbook_version_action"
description: |-
  Manages a SecMaster playbook version action resource within HuaweiCloud.
---

# huaweicloud_secmaster_playbook_version_action

Manages a SecMaster playbook version action resource within HuaweiCloud.

-> The current resource is a one-time resource, and destroying this resource will not change the current status.

## Example Usage

### Submit draft version example

```hcl
variable "workspace_id" {}
variable "version_id" {}

resource "huaweicloud_secmaster_playbook_version_action" "submit" {
  workspace_id = var.workspace_id
  version_id   = var.version_id
  status       = "APPROVING"   
}
```

### Enable or disable version example

```hcl
variable "workspace_id" {}
variable "version_id" {}
variable "enabled" {}

resource "huaweicloud_secmaster_playbook_version_action" "enabled" {
  workspace_id = var.workspace_id
  version_id   = var.version_id
  enabled      = var.enabled
}
```

### More Examples

For more detailed associated usage see [playbook instructions](/examples/secmaster/playbook/README.md)

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `workspace_id` - (Required, String, NonUpdatable) Specifies the ID of the workspace to which the playbook version belongs.

* `version_id` - (Required, String, NonUpdatable) Specifies the playbook version ID.

* `status` - (Optional, String, NonUpdatable) Specifies the playbook version status. The value can only be **APPROVING**.

* `enabled` - (Optional, Bool, NonUpdatable) Specifies whether the playbook version is enabled.
  The value can be **true**(enable version) or **false**(disable version).

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `version` - The playbook version.

* `type` - The playbook version type.

* `description` - The description.

* `created_at` - The creation time.

* `data_object_create` - Whether to trigger a playbook when a data object is created.

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
