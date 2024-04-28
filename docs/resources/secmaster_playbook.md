---
subcategory: "SecMaster"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_secmaster_playbook"
description: ""
---

# huaweicloud_secmaster_playbook

Manages a SecMaster playbook resource within HuaweiCloud.

~> This resource can only be used in region **cn-east-3** for now.

## Example Usage

```hcl
variable "workspace_id" {}

resource "huaweicloud_secmaster_playbook" "test" {
  workspace_id = var.workspace_id
  name         = "test"
  description  = "created by terraform"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `workspace_id` - (Required, String, ForceNew) Specifies the ID of the workspace to which the playbook belongs.

  Changing this parameter will create a new resource.

* `name` - (Required, String) Specifies the playbook name.

* `description` - (Optional, String) Specifies the description of the playbook.

* `enabled` - (Optional, Bool) Specifies whether to enable the playbook.

* `active_version_id` - (Optional, String) Specifies the active version ID.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `approve_role` - Indicates the approve role of the playbook.

* `created_at` - Indicates the created time of the playbook.

* `updated_at` - Indicates the updated time of the playbook.

* `dataclass_id` - Indicates the data class ID.

* `dataclass_name` - Indicates the data class name.

* `edit_role` - Indicates the edit role.

* `owner_id` - Indicates the owner ID.

* `reject_version_id` - Indicates the rejected version ID.

* `unaudited_version_id` - Indicates the unaudited version ID.

* `user_role` - Indicates the user role.

* `version` - Indicates the version.

* `version_id` - Indicates the version ID.

## Import

The playbook can be imported using the workspace ID and the playbook ID, e.g.

```bash
$ terraform import huaweicloud_secmaster_playbook.test <workspace_id>/<id>
```
