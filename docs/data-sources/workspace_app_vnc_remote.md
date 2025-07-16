---
subcategory: "Workspace"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_workspace_app_vnc_remote"
description: |-
  Use this data source to get the VNC remote information of a Workspace APP server within HuaweiCloud.
---

# huaweicloud_workspace_app_vnc_remote

Use this data source to get the VNC remote information of a Workspace APP server within HuaweiCloud.

## Example Usage

```hcl
variable "server_id" {}

data "huaweicloud_workspace_app_vnc_remote" "test" {
  server_id = var.server_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region where the APP server is located.
  If omitted, the provider-level region will be used.

* `server_id` - (Required, String) Specifies the ID of the APP server to get VNC remote information.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `url` - The remote login console address.

* `type` - The login type.
