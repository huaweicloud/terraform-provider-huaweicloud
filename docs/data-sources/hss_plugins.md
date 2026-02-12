---
subcategory: "Host Security Service (HSS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_hss_plugins"
description: |-
  Use this data source to get the list of HSS plugins within HuaweiCloud.
---

# huaweicloud_hss_plugins

Use this data source to get the list of HSS plugins within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_hss_plugins" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project ID.  
  This parameter is valid only when the enterprise project is enabled.
  The default value is **0**, indicating the default enterprise project.
  If you need to query data for all enterprise projects, the value is **all_granted_eps**.
  If you only have permissions for a specific enterprise project, you need set the enterprise project ID. Otherwise,
  the operation may fail due to insufficient permissions.

* `code` - (Optional, String) Specifies the plugin code.

* `name` - (Optional, String) Specifies the plugin name.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `total_num` - The total number of plugins.

* `data_list` - The plugin information list.

The [data_list](#data_list_struct) structure is documented below.

<a name="data_list_struct"></a>
The `data_list` block supports:

* `code` - The plugin code.

* `name` - The plugin name.

* `description` - The plugin description.

* `tags` - The plugin tags.

* `installed_attachment_num` - The number of hosts with installed plugins.

* `uninstall_attachment_num` - The number of hosts without installed plugins, including hosts with plugin status
  as uninstalled, hosts in installation, and hosts with failed installations.

* `max_cpu_limit` - The maximum single core CPU (0-100%) required to run the plugin in this type of plugin package.

* `max_memory_limit` - The maximum amount of memory (MB) required to run the plugin in this type of plugin package.

* `max_size` - The maximum plugin size (in MB) in this type of plugin package.

* `display_mode` - The plugin display mode.  
  The valid values are as follows:
  + **0**: All operation functions of the plugin are supported.
  + **1**: Installation and uninstallation of plugins are not supported.
  + **2**: All operation functions of the plugin are not supported.
