---
subcategory: "Host Security Service (HSS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_hss_plugin_info"
description: |-
  Use this data source to get the information of HSS plugin within HuaweiCloud.
---

# huaweicloud_hss_plugin_info

Use this data source to get the information of HSS plugin within HuaweiCloud.

## Example Usage

```hcl
variable "code" {}

data "huaweicloud_hss_plugin_info" "test" {
  code = var.code
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `code` - (Required, String) Specifies the plugin code.

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project ID.  
  This parameter is valid only when the enterprise project is enabled.
  The default value is **0**, indicating the default enterprise project.
  If you need to query data for all enterprise projects, the value is **all_granted_eps**.
  If you only have permissions for a specific enterprise project, you need set the enterprise project ID. Otherwise,
  the operation may fail due to insufficient permissions.

* `plugin_version` - (Optional, String) Specifies the plugin version.

* `agent_version` - (Optional, String) Specifies the agent version.

* `plugin_arch` - (Optional, String) Specifies the plugin architecture.  
  The valid values are as follows:
  + **x86_64**
  + **arm**

* `plugin_os_type` - (Optional, String) Specifies the type of operating system supported by plugin.  
  The valid values are as follows:
  + **Linux**
  + **Windows**

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `total_num` - The total number.

* `data_list` - The plugin information list.

The [data_list](#data_list_struct) structure is documented below.

<a name="data_list_struct"></a>
The `data_list` block supports:

* `name` - The plugin name.

* `id` - The plugin ID.

* `version` - The plugin version.

* `agent_version` - The minimum agent version supported by the plugin.

* `arch` - The plugin architecture.

* `os_type` - The type of operating system supported by plugin.

* `version_description` - The plugin version information description.

* `size` - The plugin installation package size (MB).

* `cpu_limit` - The single core CPU required to run plugins (0-100%).

* `memory_limit` - The memory required to run plugin (MB).

* `update_time` - The plugin update time.
