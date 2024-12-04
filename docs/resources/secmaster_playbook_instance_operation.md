---
subcategory: "SecMaster"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_secmaster_playbook_instance_operation"
description: |-
  Manages a SecMaster playbook instance operation resource within HuaweiCloud.
---

# huaweicloud_secmaster_playbook_instance_operation

Manages a SecMaster playbook instance operation resource within HuaweiCloud.

-> The current resource is a one-time resource, and destroying this resource will not change the current status.

## Example Usage

```hcl
variable "workspace_id" {}
variable "instance_id" {}
variable "operation" {}

resource "huaweicloud_secmaster_playbook_instance_operation" "test" {
  workspace_id = var.workspace_id
  instance_id  = var.instance_id
  operation    = var.operation
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.
  Changing this creates a new resource.

* `workspace_id` - (Required, String, NonUpdatable) Specifies the ID of the workspace to which the playbook instance belongs.

* `instance_id` - (Required, String, NonUpdatable) Specifies the version ID of the playbook.

* `operation` - (Required, String, NonUpdatable) Specifies the operation of the playbook instance.
  The value can be **RETRY** or **TERMINATE**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.
