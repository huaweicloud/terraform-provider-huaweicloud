---
subcategory: "Host Security Service (HSS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_hss_honeypot_port_support_list"
description: |-
  Use this data source to get the list of HSS honeypot port support host list within HuaweiCloud.
---

# huaweicloud_hss_honeypot_port_support_list

Use this data source to get the list of HSS honeypot port support host list within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_hss_honeypot_port_support_list" "test" {
  os_type = "Linux"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `os_type` - (Required, String) Specifies the operating system type. Valid values are:
  + **Linux**: Linux.
  + **Windows**: Windows.

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project ID.  
  This parameter is valid only when the enterprise project is enabled.
  The default value is **0**, indicating the default enterprise project.
  If you need to query data for all enterprise projects, the value is **all_granted_eps**.
  If you only have permissions for a specific enterprise project, you need set the enterprise project ID. Otherwise,
  the operation may fail due to insufficient permissions.

* `ports` - (Optional, String) Specifies the port numbers to set, separated by commas.

* `policy_id` - (Optional, String) Specifies the policy ID.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID in UUID format.

* `total_num` - The total number.

* `data_list` - The list of honeypot port support host.

  The [data_list](#data_list_struct) structure is documented below.

<a name="data_list_struct"></a>
The `data_list` block supports:

* `host_id` - The server (host) unique identifier ID.

* `host_name` - The server name.

* `private_ip` - The server private IP.

* `agent_id` - The agent ID.

* `conflict_port` - The conflicting ports.

* `os_type` - The operating system type.

* `group_name` - The group name.

* `group_id` - The group ID.
