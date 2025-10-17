---
subcategory: "Host Security Service (HSS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_hss_webtamper_host_management_hosts"
description: |-
  Use this data source to get the list of HSS web tamper host management hosts within HuaweiCloud.
---

# huaweicloud_hss_webtamper_host_management_hosts

Use this data source to get the list of HSS web tamper host management hosts within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_hss_webtamper_host_management_hosts" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `enterprise_project_id` - (Optional, String, NonUpdatable) Specifies the enterprise project ID.  
  This parameter is valid only when the enterprise project is enabled.
  The default value is **0**, indicating the default enterprise project.
  If it is necessary to operate the asset under all enterprise projects, the value is **all_granted_eps**.
  If you only have permissions for a specific enterprise project, you need set the enterprise project ID. Otherwise,
  the operation may fail due to insufficient permissions.

* `host_id` - (Optional, String) Specifies the host ID.

* `host_name` - (Optional, String) Specifies the host name.

* `private_ip` - (Optional, String) Specifies the private IP of the host.

* `public_ip` - (Optional, String) Specifies the public IP of the host.

* `group_id` - (Optional, String) Specifies the group ID.

* `os_type` - (Optional, String) Specifies the OS type. Valid values are:
  + **Linux**
  + **Windows**

* `web_app_name` - (Optional, String) Specifies the web application name.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID in UUID format.

* `data_list` - The list of web tamper host management hosts.

  The [data_list](#data_list_struct) structure is documented below.

<a name="data_list_struct"></a>
The `data_list` block supports:

* `host_id` - The host ID.

* `host_name` - The host name.

* `public_ip` - The public IP address.

* `private_ip` - The private IP address.

* `agent_id` - The agent ID.

* `os_type` - The OS type.

* `asset_value` - The importance of assets. Valid values are:
  + **important**: Important assets.
  + **common**: Common assets.
  + **test**: Test assets.

* `web_app_list` - The web application list.
