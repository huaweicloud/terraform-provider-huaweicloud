---
subcategory: "DataArts Studio"
layout: "huaweicloud"
page_title: "HuaweiCLoud: huaweicloud_dataarts_security_workspace_associated_queues"
description: |-
  Use this data source to get the list of workspace associated queues within HuaweiCloud.
---

# huaweicloud_dataarts_security_workspace_associated_queues

Use this data source to get the list of workspace associated queues within HuaweiCloud.

## Example Usage

```hcl
variable "workspace_id" {}

data "huaweicloud_dataarts_security_workspace_associated_queues" "test" {
  workspace_id = var.workspace_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region where the workspace associated queues are located.  
  If omitted, the provider-level region will be used.

* `workspace_id` - (Required, String) Specifies the ID of the workspace.

* `type` - (Optional, String) Specifies the type of the queue.

* `cluster_id` - (Optional, String) Specifies the ID of the cluster.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `queues` - The list of workspace associated queues that matched filter parameters.
  The [queues](#security_workspace_associated_queues) structure is documented below.

<a name="security_workspace_associated_queues"></a>
The `queues` block supports:

* `id` - The ID of the queue.

* `source_type` - The service type of the queue.

* `name` - The name of the queue.

* `type` - The type of the queue.

* `attribute` - The attribute of the queue.

* `connection_id` - The ID of the data connection.

* `connection_name` - The name of the data connection.

* `cluster_id` - The ID of the cluster.

* `cluster_name` - The name of the cluster.

* `created_at` - The time when the queue was added to the workspace, in RFC3339 format.

* `create_user` - The operator who added the queue to the workspace.

* `updated_at` - The update time of the queue in the workspace, in RFC3339 format.

* `update_user` - The updater of the queue in the workspace.

* `project_id` - The ID of the project.

* `description` - The description of the workspace associated queue.
