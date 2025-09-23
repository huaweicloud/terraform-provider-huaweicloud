---
subcategory: "SecMaster"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_secmaster_dataspace"
description: |-
  Manages a SecMaster dataspace resource within HuaweiCloud.
---

# huaweicloud_secmaster_dataspace

Manages a SecMaster dataspace resource within HuaweiCloud.

-> The current resource is a one-time resource, and destroying this resource will not change the current status.

## Example Usage

```hcl
variable "workspace_id" {}
variable "dataspace_name" {}
variable "description" {}

resource "huaweicloud_secmaster_dataspace" "test" {
  workspace_id   = var.workspace_id
  dataspace_name = var.workflow_id
  description    = var.description
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.
  Changing this creates a new resource.

* `workspace_id` - (Required, String, NonUpdatable) Sepcifies the ID of the workspace to which the dataspace belongs.

* `dataspace_name` - (Required, String, NonUpdatable) Sepcifies the name of the dataspace.

* `description` - (Required, String, NonUpdatable) Sepcifies the description of the dataspace.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.
