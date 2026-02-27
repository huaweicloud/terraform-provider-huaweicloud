---
subcategory: "Host Security Service (HSS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_hss_plugin_attachments"
description: |-
  Use this data source to get the list of HSS plugin attachments within HuaweiCloud.
---

# huaweicloud_hss_plugin_attachments

Use this data source to get the list of HSS plugin attachments within HuaweiCloud.

## Example Usage

```hcl
variable "plugin_code" {}

data "huaweicloud_hss_plugin_attachments" "test" {
  plugin_code = var.plugin_code
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `plugin_code` - (Required, String) Specifies the plugin code.

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project ID.  
  This parameter is valid only when the enterprise project is enabled.
  The default value is **0**, indicating the default enterprise project.
  If you need to query data for all enterprise projects, the value is **all_granted_eps**.
  If you only have permissions for a specific enterprise project, you need set the enterprise project ID. Otherwise,
  the operation may fail due to insufficient permissions.

* `plugin_version` - (Optional, String) Specifies the plugin version.

* `plugin_status` - (Optional, String) Specifies the plugin status.  
  The valid values are as follows:
  + **not_installed**
  + **installing**
  + **install_fail**
  + **starting**
  + **running**
  + **start_fail**
  + **offline**
  + **stopping**
  + **stopped**
  + **updating**
  + **update_fail**
  + **uninstalling**
  + **uninstall_fail**

* `host_name` - (Optional, String) Specifies the host name.

* `host_ids` - (Optional, List) Specifies the host IDs.

* `host_status` - (Optional, List) Specifies the host status.  
  The valid values are as follows:
  + **ACTIVE**
  + **BUILDING**
  + **ERROR**
  + **SHUTOFF**

* `agent_status` - (Optional, String) Specifies the agent status.  
  The valid values are as follows:
  + **not_installed**
  + **online**
  + **offline**
  + **install_failed**
  + **installing**

* `os_type` - (Optional, String) Specifies the host operating system.  
  The valid values are as follows:
  + **Linux**
  + **Windows**

* `os_arch` - (Optional, String) Specifies the system architecture.  
  The valid values are as follows:
  + **x86_64**
  + **arm**

* `host_type` - (Optional, String) Specifies the server type.  
  The valid values are as follows:
  + **host**
  + **container**

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `total_num` - The total number.

* `data_list` - The plugin status information list.

The [data_list](#data_list_struct) structure is documented below.

<a name="data_list_struct"></a>
The `data_list` block supports:

* `host_id` - The host ID.

* `host_name` - The host name.

* `host_type` - The host type.

* `private_ip` - The host private network IP address.

* `public_ip` - The host public IP address.

* `host_status` - The host status.

* `agent_status` - The agent status.

* `agent_version` - The agent version.

* `asset_value` - The importance of server assets.  
  The valid values are as follows:
  + **important**
  + **common**
  + **test**

* `os_type` - The operating system type.

* `os_arch` - The system architecture.

* `os_name` - The system name.

* `os_version` - The operating system version.

* `plugin_status` - The plugin status.

* `plugin_version` - The plugin version.

* `status_detail` - The reasons for plugin operation failure, including installation failure, startup failure, offline,
  stop failure, update failure, and uninstallation failure.

* `install_progress` - The plugin installation progress, percentage.

* `remaining_time` - The remaining installation time of plugin, in minutes.

* `protect_status` - The host protection status.  
  The valid values are as follows:
  + **closed**
  + **opened**
