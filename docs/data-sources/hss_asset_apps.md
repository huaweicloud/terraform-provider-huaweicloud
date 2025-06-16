---
subcategory: "Host Security Service (HSS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_hss_asset_apps"
description: |-
  Use this data source to query the HSS asset app list within HuaweiCloud.
---

# huaweicloud_hss_asset_apps

Use this data source to query the HSS asset app list within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_hss_asset_apps" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `enterprise_project_id` - (Optional, String) Specifies the ID of the enterprise project to which the hosts belong.
  This parameter is valid only when the enterprise project function is enabled.
  The value **all_granted_eps** indicates all enterprise projects.
  If omitted, the default enterprise project will be used.

* `host_id` - (Optional, String) Specifies the host ID.

* `host_name` - (Optional, String) Specifies the host name.

* `app_name` - (Optional, String) Specifies the software name.

* `host_ip` - (Optional, String) Specifies the host IP address.

* `version` - (Optional, String) Specifies the software version.

* `install_dir` - (Optional, String) Specifies the installation directory.

* `category` - (Optional, String) Specifies the type. The default value is **host**.
  The valid values are as follows:
  + **host**: Host.
  + **container**: Container.

* `part_match` - (Optional, Bool) Specifies whether fuzzy match is used. The default value is **false**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `data_list` - The software list.
  The [data_list](#app_structure) structure is documented below.

<a name="app_structure"></a>
The `data_list` block supports:

* `agent_id` - The agent ID.

* `host_id` - The host ID.

* `host_name` - The host name.

* `host_ip` - The host IP address.

* `app_name` - The software name.

* `version` - The version number.

* `update_time` - The latest update time, in milliseconds.

* `recent_scan_time` - The latest scanning time, in milliseconds.

* `container_id` - The container ID.

* `container_name` - The container name.
