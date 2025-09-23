---
subcategory: "SecMaster"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_secmaster_playbook_action"
description: |-
  Manages a SecMaster playbook action resource within HuaweiCloud.
---

# huaweicloud_secmaster_playbook_action

Manages a SecMaster playbook action resource within HuaweiCloud.

## Example Usage

### Basic Example

```hcl
variable "workspace_id" {}
variable "version_id" {}
variable "action_id" {}

resource "huaweicloud_secmaster_playbook_action" "name" {
  workspace_id = var.workspace_id
  version_id   = var.version_id
  action_id    = var.action_id
  name         = "test"
  description  = "created by terraform"
}
```

### More Examples

For more detailed associated usage see [playbook instructions](/examples/secmaster/playbook/README.md)

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `workspace_id` - (Required, String, ForceNew) Specifies the ID of the workspace to which the playbook action belongs.

  Changing this parameter will create a new resource.

* `version_id` - (Required, String, ForceNew) Specifies playbook version ID of the action.

  Changing this parameter will create a new resource.

* `action_id` - (Required, String) Specifies the workflow ID.

* `name` - (Optional, String) Specifies playbook action name.

* `description` - (Optional, String) Specifies the description of the playbook action.

* `action_type` - (Optional, String) Specifies the playbook action type.
  The value can be **AOP_WORKFLOW**.

* `sort_order` - (Optional, Float) Specifies the sort order of the playbook action.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `playbook_id` - Indicates the playbook ID of the action.

* `created_at` - Indicates the created time of the playbook action.

* `updated_at` - Indicates the updated time of the playbook action.

## Import

The playbookaction can be imported using the workspace ID, the playbook version ID and the playbook action ID,
separated by slashes, e.g.

```bash
$ terraform import huaweicloud_secmaster_playbook_action.test <workspace_id>/<playbook_version_id>/<playbook_action_id>
```
