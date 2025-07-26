---
subcategory: "SecMaster"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_secmaster_playbook_approval"
description: |-
  Manages a SecMaster playbook approval resource within HuaweiCloud.
---

# huaweicloud_secmaster_playbook_approval

Manages a SecMaster playbook approval resource within HuaweiCloud.

-> Destroying this resource will not change the status of the playbook approval resource.

## Example Usage

### Basic Example

```hcl
variable "workspace_id" {}
variable "version_id" {}

resource "huaweicloud_secmaster_playbook_approval" "test" {
  workspace_id = var.workspace_id
  version_id   = var.version_id
  result       = "PASS"
  content      = "ok"
}
```

### More Examples

For more detailed associated usage see [playbook instructions](/examples/secmaster/playbook/README.md)

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.
  Changing this creates a new resource.

* `workspace_id` - (Required, String, NonUpdatable) Specifies the ID of the workspace to which the playbook belongs.

* `version_id` - (Required, String, NonUpdatable) Specifies the version ID of the playbook.

* `result` - (Optional, String) Specifies the result of playbook approval. The value can be **PASS** and **UN_PASS**.

* `content` - (Optional, String) Specifies the content of the playbook approval.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.
