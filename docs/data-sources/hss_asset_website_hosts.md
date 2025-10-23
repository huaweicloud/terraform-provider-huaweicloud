---
subcategory: "Host Security Service (HSS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_hss_asset_website_hosts"
description: |-
  Use this data source to get the list of servers for a specified website.
---

# huaweicloud_hss_asset_website_hosts

Use this data source to get the list of servers for a specified website.

## Example Usage

```hcl
variable "category" {}
variable "domain" {}

data "huaweicloud_hss_asset_website_hosts" "test" {
  category = var.category
  domain   = var.domain
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `category` - (Required, String) Specifies the asset category.
  The valid values are as follows:
  + **0**: Host.
  + **1**: Container.

* `domain` - (Required, String) Specifies the domain name.

* `host_name` - (Optional, String) Specifies the host name.

* `host_ip` - (Optional, String) Specifies the host IP address.

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project ID to which the hosts belong.
  This parameter is valid only when the enterprise project is enabled.
  The default value is **0**, indicating the default enterprise project.
  If you need to query data for all enterprise projects, the value is **all_granted_eps**.
  If you only have permissions for a specific enterprise project, you need set the enterprise project ID. Otherwise,
  the operation may fail due to insufficient permissions.

* `part_match` - (Optional, Bool) Specifies whether to use fuzzy matching.
  Defaults to **false**, indicates exact matching.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID in UUID format.

* `data_list` - The list of servers.
  The [data_list](#data_list_struct) structure is documented below.

<a name="data_list_struct"></a>
The `data_list` block supports:

* `agent_id` - The agent ID.

* `host_id` - The host ID.

* `host_name` - The host name.

* `host_ip` - The host IP address.

* `domain` - The domain name.

* `app_name` - The application name.

* `path` - The path.

* `port` - The port.

* `bind_addr` - The IP address to be bound.

* `url_path` - The URL path.

* `uid` - The user ID.

* `gid` - The user group ID.

* `mode` - The file permissions.

* `pid` - The process ID.

* `proc_path` - The process path.

* `is_https` - Whether HTTPS is used.

* `cert_issuer` - The SSL certificate issuer.

* `cert_user` - The SSL certificate user.

* `cert_issue_time` - The SSL certificate issue time.

* `cert_expired_time` - The SSL certificate expiration time.

* `record_time` - The scan time.

* `container_id` - The container ID.

* `container_name` - The container name.
