---
subcategory: "SecMaster"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_secmaster_playbook_version"
description: |-
  Manages a SecMaster playbook version resource within HuaweiCloud.
---

# huaweicloud_secmaster_playbook_version

Manages a SecMaster playbook version resource within HuaweiCloud.

-> The playbook version can only be updated after matching a workflow (playbook action).

## Example Usage

### Basic Example

```hcl
variable "workspace_id" {}
variable "playbook_id" {}
variable "dataclass_id" {}

resource "huaweicloud_secmaster_playbook_version" "test" {
  workspace_id = var.workspace_id
  playbook_id  = var.playbook_id
  dataclass_id = var.dataclass_id
  description  = "created by terraform"
}
```

### More Examples

For more detailed associated usage see [playbook instructions](/examples/secmaster/playbook/README.md)

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `workspace_id` - (Required, String, ForceNew) Specifies the ID of the workspace to which the playbook version belongs.

  Changing this parameter will create a new resource.

* `playbook_id` - (Required, String, ForceNew) Specifies playbook ID of the playbook version.

  Changing this parameter will create a new resource.

* `dataclass_id` - (Required, String) Specifies the data class ID of the playbook version.

* `description` - (Optional, String) Specifies the description of the playbook version.

* `rule_enable` - (Optional, Bool) Specifies whether to enable playbook rule.

* `rule_id` - (Optional, String) Specifies the playbook rule ID.

* `trigger_type` - (Optional, String) Specifies the trigger type.
  The value can be **EVENT** and **TIMER**.

* `dataobject_create` - (Optional, Bool) Specifies whether to trigger the playbook when data object is created.

* `dataobject_delete` - (Optional, Bool) Specifies whether to trigger the playbook when data object is deleted.

* `dataobject_update` - (Optional, Bool) Specifies whether to trigger the playbook when data object is updated.

* `action_strategy` - (Optional, String) Specifies the action strategy of the playbook version.
  The value can be **ASYNC**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `created_at` - Indicates the created time of the playbook version.

* `updated_at` - Indicates the updated time of the playbook version.

* `approve_name` - Indicates the approver name.

* `creator_id` - Indicates the creator ID.

* `dataclass_name` - Indicates the data class name.

* `enabled` - Indicates whether the playbook version is enabled.

* `modifier_id` - Indicates the modifier ID.

* `status` - Indicates the status of the playbook version.
  The value can be **EDITING**, **APPROVING**, **UNPASSED** or **PUBLISHED**.

* `version` - Indicates the version number.

* `version_type` - Indicates the playbook version type.
  The value can be **0**(draft version) or **1**(official version).

## Import

The playbook version can be imported using the workspace ID, the playbook ID and the playbook version ID,
separated by slashes, e.g.

```bash
$ terraform import huaweicloud_secmaster_playbook_version.test <workspace_id>/<playbook_id>/<playbook_version_id>
```
