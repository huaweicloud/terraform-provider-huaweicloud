---
subcategory: "Host Security Service (HSS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_hss_asset_port_detail"
description: |-
  Use this data source to get the list of servers for a specifies open port.
---

# huaweicloud_hss_asset_port_detail

Use this data source to get the list of servers for a specifies open port.

## Example Usage

```hcl
variable "port" {}

data "huaweicloud_hss_asset_port_detail" "test" {
  port = var.port
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `port` - (Required, Int) Specifies the port number.

* `host_name` - (Optional, String) Specifies the host name.

* `host_ip` - (Optional, String) Specifies the host IP address.

* `type` - (Optional, String) Specifies the port type.
  The valid values are as follows:
  + **TCP**
  + **UDP**

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

* `id` - The data source ID.

* `data_list` - The list of severs information.
  The [data_list](#data_list_struct) structure is documented below.

<a name="data_list_struct"></a>
The `data_list` block supports:

* `container_id` - The container ID.

* `host_id` - The host ID.

* `host_ip` - The host IP address.

* `host_name` - The host name.

* `laddr` - The listening IP address.

* `path` - The executable file path of the process.

* `pid` - The process ID.

* `port` - The port number.

* `status` - The port status.

* `type` - The port type.

* `container_name` - The container name.

* `agent_id` - The agent ID.
