---
subcategory: "Host Security Service (HSS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_hss_asset_kernel_module_hosts"
description: |-
  Use this data source to get the list of servers for a specified kernel module.
---

# huaweicloud_hss_asset_kernel_module_hosts

Use this data source to get the list of servers for a specified kernel module.

## Example Usage

```hcl
variable "name" {}

data "huaweicloud_hss_asset_kernel_module_hosts" "test" {
  name = var.name
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `name` - (Required, String) Specifies the kernel module name.

* `host_name` - (Optional, String) Specifies the host name.

* `host_ip` - (Optional, String) Specifies the host IP address.

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project ID to which the hosts belong.
  This parameter is valid only when the enterprise project is enabled.
  The default value is **0**, indicating the default enterprise project.
  If you need to query data for all enterprise projects, the value is **all_granted_eps**.
  If you only have permissions for a specific enterprise project, you need set the enterprise project ID. Otherwise,
  the operation may fail due to insufficient permissions.

* `part_match` - (Optional, Bool) Specifies whether to use fuzzy matching. The default value is **false**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID in UUID format.

* `data_list` - The list of hosts.
  The [data_list](#data_list_struct) structure is documented below.

<a name="data_list_struct"></a>
The `data_list` block supports:

* `agent_id` - The agent ID.

* `host_id` - The host ID.

* `host_name` - The host name.

* `host_ip` - The host IP address.

* `kernel_module_info` - The kernel module information.
  The [kernel_module_info](#kernel_module_info_struct) structure is documented below.

<a name="kernel_module_info_struct"></a>
The `kernel_module_info` block supports:

* `name` - The kernel module name.

* `file_name` - The file name.

* `version` - The kernel module version.

* `srcversion` - The source code version.

* `path` - The file path.

* `size` - The file size.

* `mode` - The file permissions.

* `uid` - The file UID.

* `ctime` - The file creation time.

* `mtime` - The last modify time.

* `hash` - The file hash.

* `desc` - The kernel module description.

* `record_time` - The scan time.
