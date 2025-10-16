---
subcategory: "Host Security Service (HSS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_hss_app_whitelist_optional_hosts"
description: |-
  Use this data source to get the list of optional servers for the process whitelist.
---

# huaweicloud_hss_app_whitelist_optional_hosts

Use this data source to get the list of optional servers for the process whitelist.

## Example Usage

```hcl
data "huaweicloud_hss_app_whitelist_optional_hosts" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `host_id` - (Optional, String) Specifies the host ID.

* `host_name` - (Optional, String) Specifies the host name.

* `version` - (Optional, String) Specifies the host protection version.

* `private_ip` - (Optional, String) Specifies the private IP address of the host.

* `public_ip` - (Optional, String) Specifies the public IP address of the host.

* `policy_id` - (Optional, String) Specifies the policy ID.

* `group_id` - (Optional, String) Specifies the server group ID.

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project ID to which the hosts belong.
  This parameter is valid only when the enterprise project is enabled.
  The default value is **0**, indicating the default enterprise project.
  If you need to query data for all enterprise projects, the value is **all_granted_eps**.
  If you only have permissions for a specific enterprise project, you need set the enterprise project ID. Otherwise,
  the operation may fail due to insufficient permissions.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `data_list` - The list of process whitelist policies.

  The [data_list](#data_list_struct) structure is documented below.

<a name="data_list_struct"></a>
The `data_list` block supports:

* `agent_id` - The agent ID.

* `host_id` - The host Id.

* `host_name` - The host name.

* `public_ip` - The public IP address of the host.

* `private_ip` - The private IP address of the host.

* `asset_value` - The asset importance.

* `os_type` - The OS type.
