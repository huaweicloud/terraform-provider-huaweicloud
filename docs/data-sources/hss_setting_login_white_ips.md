---
subcategory: "Host Security Service (HSS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_hss_setting_login_white_ips"
description: |-
  Use this data source to get the list of login IP whitelist information.
---

# huaweicloud_hss_setting_login_white_ips

Use this data source to get the list of login IP whitelist information.

## Example Usage

```hcl
data "huaweicloud_hss_setting_login_white_ips" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `white_ip` - (Optional, String) Specifies the whitelist IP address or IP segment.

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project ID.
  This parameter is valid only when the enterprise project is enabled.
  The default value is **0**, indicating the default enterprise project.
  If you need to query data for all enterprise projects, the value is **all_granted_eps**.
  If you only have permissions for a specific enterprise project, you need set the enterprise project ID. Otherwise,
  the operation may fail due to insufficient permissions.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `total_num` - The total number of login IP whitelist.

* `data_list` - The list of login IP whitelist information.
  The [data_list](#login_white_ips_data_list) structure is documented below.

<a name="login_white_ips_data_list"></a>
The `data_list` block supports:

* `enabled` - The login IP whitelist enabling status.
  The valid values are as follows:
  + **true**: Indicates enabled.
  + **false**: Indicates disabled.

* `white_ip` - The whitelist IP address or IP segment.

* `total_num` - The total number of servers.

* `host_id_list` - The list of server IDs.
