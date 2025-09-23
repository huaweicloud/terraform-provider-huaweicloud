---
subcategory: "SecMaster"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_secmaster_clone_playbook_version"
description: |-
  Manages a SecMaster playbook and playbook version clone resource within HuaweiCloud.
---

# huaweicloud_secmaster_clone_playbook_version

Manages a SecMaster playbook and playbook version clone resource within HuaweiCloud.

-> The current resource is a one-time resource, and destroying this resource will not change the current status.

## Example Usage

```hcl
variable "workspace_id" {}
variable "version_id" {}
variable "name" {}

resource "huaweicloud_secmaster_clone_playbook_version" "test" {
  workspace_id = var.workspace_id
  version_id   = var.version_id
  name         = var.name
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.
  Changing this creates a new resource.

* `workspace_id` - (Required, String, NonUpdatable) Sepcifies the ID of the workspace to which the playbook belongs.

* `version_id` - (Required, String, NonUpdatable) Sepcifies the ID of the playbook version.

* `name` - (Required, String, NonUpdatable) Sepcifies the name of the clone generated playbook.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.
