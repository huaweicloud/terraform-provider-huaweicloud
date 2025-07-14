---
subcategory: "Workspace"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_workspace_desktop_connections"
description: |-
  Use this data source to get the list of the Workspace desktop connections within HuaweiCloud.
---

# huaweicloud_workspace_desktop_connections

Use this data source to get the list of the Workspace desktop connections within HuaweiCloud.

## Example Usage

### Basic Usage

```hcl
data "huaweicloud_workspace_desktop_connections" "test" {}
```

### Filter desktops by user name

```hcl
variable "user_names" {
  type = list(string)
}

data "huaweicloud_workspace_desktop_connections" "test" {
  user_names = var.user_names
}
```

### Filter desktops by connection status

```hcl
variable "connect_status" {}

data "huaweicloud_workspace_desktop_connections" "test" {
  connect_status = var.connect_status
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region where the desktop connections are located.

* `user_names` - (Optional, List) Specifies the list of desktop users to be queried.
  The user_names don't support fuzzy match.

* `connect_status` - (Optional, String) Specifies the connection status of the desktop.  
  The valid values are as follows:
  + **UNREGISTER**: The desktop is not registered or powered off.
  + **REGISTERED**: The desktop is registered and waiting for user connection.
  + **CONNECTED**: The user has successfully connected and is using the desktop.
  + **DISCONNECTED**: The desktop is disconnected from the client.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `desktop_connections` - The list of desktop connections that match the query parameters.  
  The [desktop_connections](#workspace_desktop_connections_attr) structure is documented below.

<a name="workspace_desktop_connections_attr"></a>
The `desktop_connections` block supports:

* `id` - The ID of the desktop.

* `name` - The name of the desktop.

* `connect_status` - The connection status of the desktop.

* `attach_users` - The list of users or user groups attached to the desktop.  
  The [attach_users](#workspace_desktop_connections_attach_users) structure is documented below.

<a name="workspace_desktop_connections_attach_users"></a>
The `attach_users` block supports:

* `id` - The ID of the user or user group.

* `name` - The name of the user or user group.

* `user_group` - The user group of the desktop user.  
  The valid values are as follows:
  + **sudo**: Linux administrator group.
  + **default**: Linux default user group.
  + **administrators**: Windows administrator group.
  + **users**: Windows standard user group.

* `type` - The type of the user or user group.  
  The valid values are as follows:
  + **USER**: User.
  + **GROUP**: User group.
