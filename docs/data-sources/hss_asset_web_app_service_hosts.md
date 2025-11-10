---
subcategory: "Host Security Service (HSS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_hss_asset_web_app_service_hosts"
description: |-
  Use this data source to get the list of servers for a specified web service, web application or database.
---

# huaweicloud_hss_asset_web_app_service_hosts

Use this data source to get the list of servers for a specified web service, web application or database.

## Example Usage

```hcl
variable "category" {}
variable "catalogue" {}
variable "name" {}

data "huaweicloud_hss_asset_web_app_service_hosts" "test" {
  category  = var.category
  catalogue = var.catalogue
  name      = var.name
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `category` - (Required, String) Specifies the asset category.
  The valid values are as follows:
  + **host**: Host.
  + **container**: Container.

* `catalogue` - (Required, String) Specifies the asset type.
  The valid values are as follows:
  + **web-app**: Web application.
  + **web-service**: Web service.
  + **database**: Database.

* `name` - (Required, String) Specifies the web application, web service or database name.

* `host_name` - (Optional, String) Specifies the host name.

* `host_id` - (Optional, String) Specifies the host ID.

* `host_ip` - (Optional, String) Specifies the host IP address.

* `version` - (Optional, String) Specifies the web application, web service or database version.

* `install_dir` - (Optional, String) Specifies the web application, web service or database installation directory.

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

* `data_list` - The list of hosts asset information.
  The [data_list](#data_list_struct) structure is documented below.

<a name="data_list_struct"></a>
The `data_list` block supports:

* `catalogue` - The asset type.

* `name` - The web application, web service or database name.

* `version` - The web application, web service or database version.

* `agent_id` - The agent ID.

* `install_path` - The installation path.

* `config_path` - The configuration file path.

* `uid` - The user ID.

* `gid` - The user group ID.

* `mode` - The file permissions.

* `ctime` - The file last changed time.

* `mtime` - The file last modified time.

* `atime` - The file last accessed time.

* `pid` - The process ID.

* `proc_path` - The process path.

* `container_id` - The container ID.

* `container_name` - The container name.

* `record_time` - The scan time.

* `host_id` - The host ID.

* `host_name` - The host name.

* `host_ip` - The host IP address.
