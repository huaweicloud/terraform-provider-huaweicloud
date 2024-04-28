---
subcategory: "DataArts Studio"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dataarts_architecture_process"
description: ""
---

# huaweicloud_dataarts_architecture_process

Manages a DataArts Architecture process resource within HuaweiCloud.

## Example Usage

### Create a root process

```hcl
variable "workspace_id" {}
variable "process_name" {}
variable "process_owner" {}

resource "huaweicloud_dataarts_architecture_process" "test-root" {
  workspace_id = var.workspace_id
  name         = var.process_name
  owner        = var.process_owner
}
```

### Create a subordinate process

```hcl
variable "workspace_id" {}
variable "process_name" {}
variable "process_owner" {}
variable "parent_process_id" {}

resource "huaweicloud_dataarts_architecture_process" "test-sub" {
  workspace_id = var.workspace_id
  name         = var.process_name
  owner        = var.process_owner
  parent_id    = var.parent_process_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `workspace_id` - (Required, String, ForceNew) Specifies the ID of DataArts Studio workspace.
  Changing this parameter will create a new resource.

* `name` - (Required, String) Specifies the name of process.

* `owner` - (Required, String) Specifies the name of person responsible for process. The responsible person must exist
  in the system.

* `description` - (Optional, String) Specifies the description of process.

* `parent_id` - (Optional, String) Specifies the parent catalog ID of process.  
  It's **Required** when you create a subordinate process, and directory level cannot exceed `3` layer.

* `prev_id` - (Optional, String) Specifies the ID of the previous node in the process.  
  When querying the process, if the field is null, it means that is the first node in the current process directory.

* `next_id` - (Optional, String) Specify the ID of the next node in the process.  
  When querying the process, if the field is null, it means that is the last node in the current process directory.

-> The `prev_id` and `next_id` will change with the position of the process in the current directory.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `qualified_id` - The ID of all superior processes. Format is `<root_process_id>.<sub_process_id1>.<sub_process_id2>`

* `created_at` - The creation time of the process.

* `updated_at` - The latest update time of the process.

* `created_by` - The creator of the process.

* `updated_by` - The last editor of the process.

* `children` - The name list of subordinate process.

## Import

The DataArts architecture process can be imported using `workspace_id` and `qualified_id`, separated by a slash, e.g.

```bash
$ terraform import huaweicloud_dataarts_architecture_process.test <workspace_id>/<qualified_id>
```
