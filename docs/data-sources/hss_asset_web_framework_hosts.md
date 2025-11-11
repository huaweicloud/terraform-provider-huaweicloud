---
subcategory: "Host Security Service (HSS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_hss_asset_web_framework_hosts"
description: |-
  Use this data source to get the list of servers for a specified web framework.
---

# huaweicloud_hss_asset_web_framework_hosts

Use this data source to get the list of servers for a specified web framework.

## Example Usage

```hcl
variable "category" {}
variable "catalogue" {}
variable "name" {}

data "huaweicloud_hss_asset_web_framework_hosts" "test" {
  category  = var.category
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

* `file_name` - (Optional, String) Specifies the file name.

* `host_name` - (Optional, String) Specifies the host name.

* `host_id` - (Optional, String) Specifies the host ID.

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

* `data_list` - The list of servers information.
  The [data_list](#data_list_struct) structure is documented below.

<a name="data_list_struct"></a>
The `data_list` block supports:

* `agent_id` - The agent ID.

* `host_id` - The host ID.

* `host_name` - The host name.

* `host_ip` - The host IP address.

* `name` - The web framework name.

* `version` - The web framework version.

* `path` - The web framework file path.

* `record_time` - The web framework scan time.

* `catalogue` - The software type.

* `file_name` - The web framework file name.

* `file_type` - The web framework file type.

* `gid` - The web framework GID.

* `hash` - The web framework file hash.

* `is_embedded` - Whether the file is compressed.

* `mode` - The file permissions.

* `pid` - The web framework process ID.

* `proc_path` - The web framework process path.

* `uid` - The web framework UID.

* `container_id` - The container ID.

* `container_name` - The container name.
