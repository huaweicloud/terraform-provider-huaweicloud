---
subcategory: "Host Security Service (HSS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_hss_honeypot_port_hosts"
description: |-
  Use this data source to get the list of protected hosts.
---

# huaweicloud_hss_honeypot_port_hosts

Use this data source to get the list of protected hosts.

## Example Usage

```hcl
data "huaweicloud_hss_honeypot_port_hosts" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `group_id` - (Optional, String) Specifies the server group ID.

* `host_name` - (Optional, String) Specifies the server name.

* `private_ip` - (Optional, String) Specifies the server private IP address.

* `policy_id` - (Optional, String) Specifies the policy ID.

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project ID.
  This parameter is valid only when the enterprise project is enabled.
  The default value is **0**, indicating the default enterprise project.
  If you need to query data for all enterprise projects, the value is **all_granted_eps**.
  If you only have permissions for a specific enterprise project, you need set the enterprise project ID. Otherwise,
  the operation may fail due to insufficient permissions.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `data_list` - The list of hosts.
  The [data_list](#hosts_data_list) structure is documented below.

<a name="hosts_data_list"></a>
The `data_list` block supports:

* `host_id` - The host ID.

* `host_name` - The host name.

* `private_ip` - The host private IP address.

* `agent_id` - The agent ID.

* `conflict_port` - The conflicting ports.

* `applied_port` - The application ports.
