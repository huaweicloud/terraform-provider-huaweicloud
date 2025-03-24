---
subcategory: "SecMaster"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_secmaster_workspace"
description: |-
  Manages a SecMaster workspace resource within HuaweiCloud.
---

# huaweicloud_secmaster_workspace

Manages a SecMaster workspace resource within HuaweiCloud.

-> Destroying this resource will not change the status of the workspace resource.

## Example Usage

### Basic Example

```hcl
variable "name" {}
variable "project_name" {}

resource "huaweicloud_secmaster_workspace" "test" {
  name         = var.name
  project_name = var.project_name
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.
  Changing this creates a new resource.

* `name` - (Required, String, NonUpdatable) Specifies the name of the workspace.

* `project_name` - (Required, String, NonUpdatable) Specifies name of the project in which to create the workspace.

* `description` - (Optional, String, NonUpdatable) Specifies the description of the workspace.

* `enterprise_project_id` - (Optional, String, NonUpdatable) Specifies ID of the enterprise project
  in which to create the workspace.

* `enterprise_project_name` - (Optional, String, NonUpdatable) Specifies name of the enterprise project
  in which to create the workspace.

* `view_bind_id` - (Optional, String, NonUpdatable) Specifies the space ID bound to the view.

* `is_view` - (Optional, String, NonUpdatable) Specifies whether the view is used.

* `tags` - (Optional, Map, NonUpdatable) Specifies the tags of the workspace in key/pair format.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.
