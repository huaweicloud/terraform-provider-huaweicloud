---
subcategory: "Host Security Service (HSS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_hss_setting_plugins"
description: |-
  Use this data source to get the list of plug-ins.
---

# huaweicloud_hss_setting_plugins

Use this data source to get the list of plug-ins.

## Example Usage

```hcl
variable "name" {}

data "huaweicloud_hss_setting_plugins" "test" {
  name = var.name
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `name` - (Required, String) Specifies the plug-in name.
  The valid value is **opa-docker-authz**.

* `host_name` - (Optional, String) Specifies the server name.

* `host_id` - (Optional, String) Specifies the server ID.

* `private_ip` - (Optional, String) Specifies the server private IP address.

* `public_ip` - (Optional, String) Specifies the server EIP.

* `group_id` - (Optional, String) Specifies the server group ID.

* `asset_value` - (Optional, String) Specifies the asset importance.
  The valid values are as follows:
  + **important**
  + **common**
  + **test**

* `agent_status` - (Optional, String) Specifies the agent status.

* `detect_result` - (Optional, String) Specifies the detection result.

* `host_status` - (Optional, String) Specifies the host status.

* `os_type` - (Optional, String) Specifies the OS type.
  The valid values are as follows:
  + **Linux**
  + **Windows**

* `ip_addr` - (Optional, String) Specifies the private IP address or EIP.

* `protect_status` - (Optional, String) Specifies the protection status.

* `group_name` - (Optional, String) Specifies the server group name.

* `policy_group_id` - (Optional, String) Specifies the policy group ID.

* `policy_group_name` - (Optional, String) Specifies the policy group name.

* `label` - (Optional, String) Specifies the asset tag.

* `charging_mode` - (Optional, String) Specifies the billing mode.

* `refresh` - (Optional, Bool) Specifies whether to forcibly synchronize servers from ECS.
  The valid values are as follows:
  + **true**
  + **false** (Default value)

* `above_version` - (Optional, Bool) Specifies whether to return all the versions later than the current version.

* `version` - (Optional, String) Specifies the enabled HSS edition.

* `plugin` - (Optional, String) Specifies the plug-in type.
  The valid value is **opa-docker-authz**.

* `outside_host` - (Optional, Bool) Specifies whether a server is a non-Huawei Cloud server.
  The valid values are as follows:
  + **true**
  + **false** (Default value)

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project ID.
  This parameter is valid only when the enterprise project is enabled.
  The default value is **0**, indicating the default enterprise project.
  If you need to query data for all enterprise projects, the value is **all_granted_eps**.
  If you only have permissions for a specific enterprise project, you need set the enterprise project ID. Otherwise,
  the operation may fail due to insufficient permissions.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `data_list` - The plug-ins information.
  The [data_list](#data_list_struct) structure is documented below.

<a name="data_list_struct"></a>
The `hosts` block supports:

* `host_name` - The server name.

* `host_id` - The server ID.

* `private_ip` - The server private IP address.

* `public_ip` - The server EIP.

* `os_type` - The OS type.

* `plugin_name` - The plug-in name.

* `plugin_version` - The plug-in version.

* `plugin_status` - The plug-in status.

* `upgrade_status` - The plug-in upgrade status.
