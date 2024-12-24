---
subcategory: "CodeArts Deploy"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_codearts_deploy_application_group_move"
description: |-
  Manages a CodeArts deploy application group move resource within HuaweiCloud.
---

# huaweicloud_codearts_deploy_application_group_move

Manages a CodeArts deploy application group move resource within HuaweiCloud.

## Example Usage

```hcl
variable "project_id" {}
variable "group_id" {}
variable "movement" {}

resource "huaweicloud_codearts_deploy_application_group" "test" {
  project_id = var.project_id
  group_id   = var.group_id
  movement   = var.movement
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.
  Changing this creates a new resource.

* `project_id` - (Required, String, ForceNew) Specifies the project ID for CodeArts service.
  Changing this creates a new resource.

* `group_id` - (Required, String, ForceNew) Specifies the application group ID.
  Changing this creates a new resource.

* `movement` - (Required, Int, ForceNew) Specifies the moving direction.
  Valid values are as follows:
  + **1**: Upward.
  + **-1**: Downward.

  Changing this creates a new resource.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.
