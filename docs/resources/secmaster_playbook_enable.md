---
subcategory: "SecMaster"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_secmaster_playbook_enable"
description: |-
  Use this resource to enable a SecMaster playbook within HuaweiCloud.
---

# huaweicloud_secmaster_playbook_enable

Use this resource to enable a SecMaster playbook within HuaweiCloud.

## Example Usage

### Basic Example

```hcl
variable "workspace_id" {}
variable "playbook_id" {}
variable "playbook_name" {}
variable "active_version_id" {}

resource "huaweicloud_secmaster_playbook_enable" "test" {
  workspace_id      = var.workspace_id
  playbook_id       = var.playbook_id
  playbook_name     = var.playbook_name
  active_version_id = var.active_version_id
}
```

### More Examples

For more detailed associated usage see [playbook instructions](/examples/secmaster/playbook/README.md)

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.
  Changing this creates a new resource.

* `workspace_id` - (Required, String, NonUpdatable) Specifies the ID of the workspace to which the playbook version belongs.

* `playbook_id` - (Required, String, NonUpdatable) Specifies the playbook ID.

* `playbook_name` - (Required, String, NonUpdatable) Specifies the playbook name.

* `active_version_id` - (Required, String, NonUpdatable) Specifies the actived playbook version ID.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.
