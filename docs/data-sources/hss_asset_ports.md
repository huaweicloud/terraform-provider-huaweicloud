---
subcategory: "Host Security Service (HSS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_hss_asset_ports"
description: |-
  Use this data source to get the list of open ports for a specifies host.
---

# huaweicloud_hss_asset_ports

Use this data source to get the list of open ports for a specifies host.

## Example Usage

```hcl
variable "host_id" {}

data "huaweicloud_hss_asset_ports" "test" {
  host_id = var.host_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `host_id` - (Required, String) Specifies the host ID.

* `host_name` - (Optional, String) Specifies the host name.

* `host_ip` - (Optional, String) Specifies the host IP address.

* `port` - (Optional, Int) Specifies the port number.

* `type` - (Optional, String) Specifies the port type.
  The valid values are as follows:
  + **TCP**
  + **UDP**

* `status` - (Optional, String) Specifies the port status.
  The valid values are as follows:
  + **normal**: Normal.
  + **danger**: Dangerous.
  + **unknow**: Unknown.

* `category` - (Optional, String) Specifies the category.
  The valid values are as follows:
  + **host**: Host (default).
  + **container**: Container.

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project ID to which the hosts belong.
  This parameter is valid only when the enterprise project is enabled.
  The default value is **0**, indicating the default enterprise project.
  If you need to query data for all enterprise projects, the value is **all_granted_eps**.
  If you only have permissions for a specific enterprise project, you need set the enterprise project ID. Otherwise,
  the operation may fail due to insufficient permissions.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID in UUID format.

* `data_list` - All ports that match the filter parameters.
  The [data_list](#data_list_struct) structure is documented below.

<a name="data_list_struct"></a>
The `data_list` block supports:

* `host_id` - The host ID.

* `laddr` - The listening IP address.

* `status` - The port status.

* `port` - The port number.

* `type` - The port type.

* `pid` - The process ID.

* `path` - The executable file path of the process.

* `agent_id` - The agent ID.

* `container_id` - The container ID.
