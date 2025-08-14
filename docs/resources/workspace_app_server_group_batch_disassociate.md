---
subcategory: "Workspace"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_workspace_app_server_group_batch_disassociate"
description: |-
  Use this resource to disassociate all application groups from Workspace APP server group within HuaweiCloud.
---

# huaweicloud_workspace_app_server_group_batch_disassociate

Use this resource to disassociate all application groups from Workspace APP server group within HuaweiCloud.

-> This resource is a one-time action resource used to disassociate all application groups from specified server group.
   Deleting resource will not clear the corresponding request record, but will only remove the resource information from
   the tfstate file.

## Example Usage

```hcl
variable "server_group_id" {}

resource "huaweicloud_workspace_app_server_group_batch_disassociate" "test" {
  server_group_id = var.server_group_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region where the server group is located.  
  Changing this creates a new resource.

* `server_group_id` - (Required, String) Specifies the ID of the server group to disassociate all application groups.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.
