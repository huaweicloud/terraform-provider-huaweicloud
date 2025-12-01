---
subcategory: "Workspace"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_workspace_app_server_batch_migrate"
description: |-
  Use this resource to batch migrate servers to the target cloud office host within HuaweiCloud.
---

# huaweicloud_workspace_app_server_batch_migrate

Use this resource to batch migrate servers to the target cloud office host within HuaweiCloud.

-> This resource is a one-time action resource used to batch migrate servers. Deleting this resource will not
   clear the corresponding request record, but will only remove the resource information from the tfstate file.

## Example Usage

```hcl
variable "server_ids" {
  type = list(string)
}
variable "target_host_id" {}

resource "huaweicloud_workspace_app_server_batch_migrate" "test" {
  server_ids = var.server_ids
  host_id    = var.target_host_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region where the servers to be migrated are located.  
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `server_ids` - (Required, List, NonUpdatable) Specifies the list of server IDs to be migrated.  

* `host_id` - (Required, String, NonUpdatable) Specifies the ID of the target cloud office host.  

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 30 minutes.
