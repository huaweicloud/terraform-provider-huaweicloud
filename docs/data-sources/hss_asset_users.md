---
subcategory: "Host Security Service (HSS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_hss_asset_users"
description: |-
  Use this data source to get the list of HSS asset users within HuaweiCloud.
---

# huaweicloud_hss_asset_users

Use this data source to get the list of HSS asset users within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_hss_asset_users" "test" {
  user_name = "daemon"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `host_id` - (Optional, String) Specifies the host ID.

* `user_name` - (Optional, String) Specifies the account name.

* `host_name` - (Optional, String) Specifies the host name.

* `private_ip` - (Optional, String) Specifies the private IP of the server.

* `login_permission` - (Optional, Bool) Specifies whether the user has the login permission.

* `root_permission` - (Optional, Bool) Specifies whether the user has root permissions.

* `user_group` - (Optional, String) Specifies the server user group.

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project ID.
  This parameter is valid only when the enterprise project function is enabled.
  The value **all_granted_eps** indicates all enterprise projects.
  If omitted, the default enterprise project will be used.

* `category` - (Optional, String) Specifies the type. The default value is **host**.
  The valid values are as follows:
  + **host**: Host
  + **container**: Container

* `part_match` - (Optional, Bool) Specifies whether to use fuzzy matching. The default value is **false**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `data_list` - The list of account information.
  The [data_list](#user_info) structure is documented below.

<a name="user_info"></a>
The `data_list` block supports:

* `agent_id` - The agent ID.

* `host_id` - The host ID.

* `host_name` - The host name.

* `host_ip` - The host IP.

* `user_name` - The user name.

* `login_permission` - Whether the user has the login permission.

* `root_permission` - Whether the user has root permissions.

* `user_group_name` - The user group name.

* `user_home_dir` - The user home directory.

* `shell` - The user startup shell.

* `recent_scan_time` - The latest scan time.

* `container_id` - The container ID.

* `container_name` - The container name.
