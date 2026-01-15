---
subcategory: "SecMaster"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_secmaster_pipe_consumption"
description: |-
  Manages a pipe consumption resource of SecMaster within HuaweiCloud.
---

# huaweicloud_secmaster_pipe_consumption

Manages a pipe consumption resource of SecMaster within HuaweiCloud.

-> This resource is used to enable pipe consumption. Destroying this resource will disable pipe consumption.

## Example Usage

```hcl
variable "workspace_id" {}
variable "pipe_id" {}

resource "huaweicloud_secmaster_pipe_consumption" "test" {
  workspace_id = var.workspace_id
  pipe_id      = var.pipe_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.
  Changing this parameter will create a new resource.

* `workspace_id` - (Required, String, NonUpdatable) Specifies the ID of the workspace.

* `pipe_id` - (Required, String, NonUpdatable) Specifies the pipe ID.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID (pipe ID).

* `access_point` - The access point.

* `pipe_name` - The pipe name.

* `status` - The status.

* `subscription_name` - The subscription name.

* `table_id` - The table ID.

## Import

The pipe consumption can be imported using the `workspace_id` and their `pipe_id`, separated by a slash (/), e.g.

```bash
$ terraform import huaweicloud_secmaster_pipe_consumption.test <workspace_id>/<pipe_id>
```
