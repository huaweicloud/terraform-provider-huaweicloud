---
subcategory: "Workspace"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_workspace_app_server_group_status"
description: |-
  Use this data source to get the status of a server group within HuaweiCloud Workspace.
---

# huaweicloud_workspace_app_server_group_status

Use this data source to get the status of a server group within HuaweiCloud Workspace.

## Example Usage

```hcl
variable "server_group_id" {}

data "huaweicloud_workspace_app_server_group_status" "test" {
  server_group_id = var.server_group_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `server_group_id` - (Required, String) Specifies the unique identifier of the server group.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `aps_status` - The number of servers in each status.
