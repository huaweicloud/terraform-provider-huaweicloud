---
subcategory: "Host Security Service (HSS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_hss_host_statistics"
description: |-
  Use this data source to get the server statistics.
---

# huaweicloud_hss_host_statistics

Use this data source to get the server statistics.

## Example Usage

```hcl
data "huaweicloud_hss_host_statistics" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project ID.
  This parameter is valid only when the enterprise project is enabled.
  The default value is **0**, indicating the default enterprise project.
  If you need to query data for all enterprise projects, the value is **all_granted_eps**.
  If you only have permissions for a specific enterprise project, you need set the enterprise project ID. Otherwise,
  the operation may fail due to insufficient permissions.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `total_num` - The total number of servers.

* `risk_num` - The unsafe servers.

* `unprotected_num` - The number of unprotected servers.

* `not_installed_num` - The number of servers without agent installed.

* `installed_failed_num` - The number of servers where agent installation failed.

* `not_online_num` - The number of servers without online agent.

* `version_basic_num` - The number of servers that have enabled basic version protection.

* `version_advanced_num` - The number of servers that have enabled professional edition protection.

* `version_enterprise_num` - The number of servers that have enabled enterprise edition protection.

* `version_premium_num` - The number of servers that have enabled premium edition protection.

* `version_wtp_num` - The number of servers that have enabled WTP edition protection.

* `version_container_num` - The number of servers that have enabled container edition protection.

* `host_group_num` - The total number of the server groups.

* `server_group_num` - The total number of the asset server groups.

* `asset_value_list` - The asset importance list.

  The [asset_value_list](#asset_value_list_struct) structure is documented below.

* `server_group_list` - The server group list.

  The [server_group_list](#server_group_list_struct) structure is documented below.

* `ignore_host_num` - The number of ignored servers.

* `protected_num` - The number of protected servers.

* `protect_interrupt_num` - The number of servers with protection interruption.

* `idle_num` - The number of idle quota.

* `premium_non_sp_num` - The number of servers with agent self-protection disabled of the premium edition.

<a name="asset_value_list_struct"></a>
The `asset_value_list` block supports:

* `value_type` - The asset importance type.

* `host_num` - The number of servers.

<a name="server_group_list_struct"></a>
The `server_group_list` block supports:

* `server_group_id` - The server group ID.

* `server_group_name` - The server group name.

* `host_num` - The number of servers allocated to the server group.
