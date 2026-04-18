---
subcategory: "DataArts Studio"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dataarts_security_workspace_queue_associate"
description: |-
  Manages a DataArts Security workspace queue associate resource within HuaweiCloud.
---

# huaweicloud_dataarts_security_workspace_queue_associate

Manages a DataArts Security workspace queue associate resource within HuaweiCloud.

## Example Usage

### Create a DLI associate

```hcl
variable "workspace_id" {}
variable "connection_id" {}
variable "queue_name" {}

resource "huaweicloud_dataarts_security_workspace_queue_associate" "test" {
  workspace_id  = var.workspace_id
  connection_id = var.connection_id
  queue_name    = var.queue_name
  source_type   = "dli"
}
```

### Create an MRS associate

```hcl
variable "workspace_id" {}
variable "connection_id" {}
variable "cluster_id" {}
variable "queue_name" {}

resource "huaweicloud_dataarts_security_workspace_queue_associate" "test" {
  workspace_id  = var.workspace_id
  connection_id = var.connection_id
  queue_name    = var.queue_name
  cluster_id    = var.cluster_id
  source_type   = "mrs"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region where the workspace queue associate is located.  
  If omitted, the provider-level region will be used.  
  Changing this creates a new resource.

* `workspace_id` - (Required, String, NonUpdatable) Specifies the ID of the workspace to which the queue resource
  is assigned.

* `source_type` - (Required, String, NonUpdatable) Specifies the queue resource service type.  
  The valid values are as follows:
  + **mrs**
  + **dli**

* `queue_name` - (Required, String, NonUpdatable) Specifies the queue name.

* `connection_id` - (Required, String, NonUpdatable) Specifies the data connection ID.

* `cluster_id` - (Optional, String, NonUpdatable) Specifies the ID of MRS cluster.  
  This parameter is **required** when the `source_type` is **mrs**.

* `description` - (Optional, String, NonUpdatable) Specifies the description of the queue resource assigned to
  the workspace.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `queue_type` - The type of the queue.

* `queue_attribute` - The attribute of the queue.

* `connection_name` - The data connection name.

* `cluster_name` - The cluster name.

* `created_at` - The time when the queue was added to the workspace, in RFC3339 format.

* `create_user` - The user who added the queue to the workspace.

* `updated_at` - The time when the queue resource under the workspace was updated.

* `update_user` - The user who updated the queue resource under the workspace, in RFC3339 format.

## Import

The workspace queue associate can be imported using `workspace_id` and `queue_name`, separated by a slash (/), e.g.

```bash
$ terraform import huaweicloud_dataarts_security_workspace_queue_associate.test <workspace_id>/<queue_name>
```
