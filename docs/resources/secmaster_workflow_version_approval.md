---
subcategory: "SecMaster"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_secmaster_workflow_version_approval"
description: |-
  Manages a resource to approval workflow version within HuaweiCloud.
---

# huaweicloud_secmaster_workflow_version_approval

Manages a resource to approval workflow version within HuaweiCloud.

-> This resource is a one-time action resource. Deleting this resource will not clear the corresponding request record,
  but will only remove the resource information from the tf state file.

## Example Usage

```hcl
variable "workspace_id" {}
variable "version_id" {}
variable "content" {}
variable "result" {}

resource "huaweicloud_secmaster_workflow_version_approval" "test" {
  workspace_id = var.workspace_id
  version_id   = var.version_id
  content      = var.content
  result       = var.result
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.
  Changing this creates a new resource.

* `workspace_id` - (Required, String, NonUpdatable) Specifies the workspace ID to which the workflow belongs.

* `version_id` - (Required, String, NonUpdatable) Specifies the workflow version ID.

* `content` - (Required, String, NonUpdatable) Specifies the workflow version approval comments.

* `result` - (Required, String, NonUpdatable) Specifies the workflow version approval result.
  The value can be **PASS** or **UN_PASS**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.
