---
subcategory: "Host Security Service (HSS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_hss_asset_app_change_history"
description: |-
  Use this data source to get the historical change records of software information.
---

# huaweicloud_hss_asset_app_change_history

Use this data source to get the historical change records of software information.

## Example Usage

```hcl
data "huaweicloud_hss_asset_app_change_history" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `host_id` - (Optional, String) Specifies the host ID.

* `host_ip` - (Optional, String) Specifies the host IP address.

* `host_name` - (Optional, String) Specifies the host name.

* `app_name` - (Optional, String) Specifies the software name.

* `variation_type` - (Optional, String) Specifies the change type.
  The valid values are as follows:
  + **add**
  + **delete**
  + **modify**

* `start_time` - (Optional, Int) Specifies the query start time.
  The format is 13-digit timestamp in millisecond.

* `end_time` - (Optional, Int) Specifies the query end time.
  The format is 13-digit timestamp in millisecond.

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project ID to which the hosts belong.
  This parameter is valid only when the enterprise project is enabled.
  The default value is **0**, indicating the default enterprise project.
  If you need to query data for all enterprise projects, the value is **all_granted_eps**.
  If you only have permissions for a specific enterprise project, you need set the enterprise project ID. Otherwise,
  the operation may fail due to insufficient permissions.

* `sort_key` - (Optional, String) Specifies the sort key.
  Currently, data can only be sorted by `recent_scan_time`.

* `sort_dir` - (Optional, String) Specifies the sort order. The default value is **desc**.
  The valid values are as follows:
  + **asc**
  + **desc**

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `data_list` - The list of software change history.

  The [data_list](#data_list_struct) structure is documented below.

<a name="data_list_struct"></a>
The `data_list` block supports:

* `agent_id` - The agent ID.

* `variation_type` - The change type.

* `host_id` - The host Id.

* `app_name` - The software name.

* `host_name` - The host name.

* `host_ip` - The host IP address.

* `version` - The software version.

* `update_time` - The software update time.

* `recent_scan_time` - The last scan time.
