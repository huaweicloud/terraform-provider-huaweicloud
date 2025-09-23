---
subcategory: "Host Security Service (HSS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_hss_setting_login_common_ips"
description: |-
  Use this data source to get the list of common login IP addresses information.
---

# huaweicloud_hss_setting_login_common_ips

Use this data source to get the list of common login IP addresses information.

## Example Usage

```hcl
data "huaweicloud_hss_setting_login_common_ips" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `ip_addr` - (Optional, String) Specifies the login IP address.
  The IP address can a specify IP or a network segment. e.g. **192.68.78.3** or **192.78.10.0/24**.

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project ID.
  This parameter is valid only when the enterprise project is enabled.
  The default value is **0**, indicating the default enterprise project.
  If you need to query data for all enterprise projects, the value is **all_granted_eps**.
  If you only have permissions for a specific enterprise project, you need set the enterprise project ID. Otherwise,
  the operation may fail due to insufficient permissions.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `total_num` - The total number of common login IPs that have been set.

* `data_list` - The list of common login IP information.
  The [data_list](#login_common_ips_data_list) structure is documented below.

<a name="login_common_ips_data_list"></a>
The `data_list` block supports:

* `ip_addr` - The IP address.

* `total_num` - The total number of hosts associated with the IP address.

* `host_id_list` - The list of server IDs associated with the IP address.
