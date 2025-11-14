---
subcategory: "Host Security Service (HSS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_hss_auto_launch_change_history"
description: |-
  Use this data source to get the list of historical change records of auto launch items.
---

# huaweicloud_hss_auto_launch_change_history

Use this data source to get the list of historical change records of auto launch items.

## Example Usage

```hcl
data "huaweicloud_hss_auto_launch_change_history" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `host_id` - (Optional, String) Specifies the host ID.

* `host_ip` - (Optional, String) Specifies the host IP address.

* `host_name` - (Optional, String) Specifies the host name.

* `auto_launch_name` - (Optional, String) Specifies the auto launch item name.

* `type` - (Optional, String) Specifies the auto launch item type.
  The valid values are as follows:
  + **0**: Auto launch service.
  + **1**: Scheduled task.
  + **2**: Preloaded dynamic library.
  + **3**: Run registry key.
  + **4**: Startup folder.

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
  The valid value is **recent_scan_time**.

* `sort_dir` - (Optional, String) Specifies the sort direction. The default value is **desc**.
  The valid values are as follows:
  + **asc**: Ascending order.
  + **desc**: Descending order.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `data_list` - The list of change records.
  The [data_list](#data_list_struct) structure is documented below.

<a name="data_list_struct"></a>
The `data_list` block supports:

* `agent_id` - The agent ID.

* `variation_type` - The change type.

* `type` - The auto launch item type.

* `host_id` - The host ID.

* `host_name` - The host name.

* `host_ip` - The host IP address.

* `path` - The path of the auto launch item.

* `hash` - The file hash value generated using the SHA256 algorithm.

* `run_user` - The user who starts the execution.

* `name` - The auto launch item name.

* `recent_scan_time` - The latest scan time.
