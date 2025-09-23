---
subcategory: "Workspace"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_workspace_app_hda_configurations"
description: |-
  Use this data source to get the list of HDA configurations within HuaweiCloud.
---

# huaweicloud_workspace_app_hda_configurations

Use this data source to get the list of HDA configurations within HuaweiCloud.

## Example Usage

```hcl
variable "server_name" {}

data "huaweicloud_workspace_app_hda_configurations" "test" {
  server_name = var.server_name
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region where the HDA configurations are located.  
  If omitted, the provider-level region will be used.

* `server_group_id` - (Optional, String) Specifies the ID of the server group to be queried.

* `server_name` - (Optional, String) Specifies the name of the server to be queried.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `configurations` - The list of HDA configurations that match the query parameters.  
  The [configurations](#workspace_app_hda_configurations) structure is documented below.

<a name="workspace_app_hda_configurations"></a>
The `configurations` block supports:

* `server_id` - The ID of the server.

* `machine_name` - The machine name of the server.

* `maintain_status` - Whether the server is in maintenance status.

* `server_name` - The name of the server.

* `server_group_id` - The ID of the server group.

* `server_group_name` - The name of the server group.

* `sid` - The SID of the server.

* `session_count` - The number of sessions.

* `status` - The status of the server.
  + **UNREGISTER** - Not registered
  + **REGISTERED** - Registered and ready
  + **MAINTAINING** - Under maintenance
  + **FREEZE** - Frozen
  + **STOPPED** - Stopped
  + **NONE** - Abnormal status

* `current_version` - The current version of the access agent.
