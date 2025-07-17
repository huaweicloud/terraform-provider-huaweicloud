---
subcategory: "Workspace"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_workspace_desktop_remote_console"
description: |-
  Use this data source to get the remote console information of a Workspace desktop within HuaweiCloud.
---

# huaweicloud_workspace_desktop_remote_console

Use this data source to get the remote console information of a Workspace desktop within HuaweiCloud.

## Example Usage

```hcl
variable "desktop_id" {}

data "huaweicloud_workspace_desktop_remote_console" "test" {
  desktop_id = var.desktop_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region where the desktop is located.  
  If omitted, the provider-level region will be used.

* `desktop_id` - (Required, String) Specifies the ID of the desktop.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `remote_console` - The remote console information.  
  The [remote_console](#workspace_desktop_remote_console) structure is documented below.

<a name="workspace_desktop_remote_console"></a>
The `remote_console` block supports:

* `type` - The login type of console.

* `url` - The remote login URL of console.

* `protocol` - The login protocol of console.
