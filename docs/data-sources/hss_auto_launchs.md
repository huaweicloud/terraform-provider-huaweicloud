---
subcategory: "Host Security Service (HSS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_hss_auto_launchs"
description: |-
  Use this data source to get the list of auto launch items.
---

# huaweicloud_hss_auto_launchs

Use this data source to get the list of auto launch items.

## Example Usage

```hcl
data "huaweicloud_hss_auto_launchs" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `host_id` - (Optional, String) Specifies the host ID.

* `host_name` - (Optional, String) Specifies the host name.

* `host_ip` - (Optional, String) Specifies the host IP address.

* `name` - (Optional, String) Specifies the auto launch item name.

* `type` - (Optional, String) Specifies the auto launch item type.
  The valid values are as follows:
  + **0**: Auto launch service.
  + **1**: Scheduled task.
  + **2**: Preloaded dynamic library.
  + **3**: Run registry key.
  + **4**: Startup folder.

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

* `id` - The data source ID.

* `data_list` - All auto launch items that match the filter parameters.
  The [data_list](#data_list_struct) structure is documented below.

<a name="data_list_struct"></a>
The `data_list` block supports:

* `host_id` - The host ID.

* `host_name` - The host name.

* `host_ip` - The host IP address.

* `agent_id` - The agent ID.

* `name` - The auto launch item name.

* `type` - The auto launch item type.

* `path` - The path of the auto launch item.

* `hash` - The file hash value generated using the SHA256 algorithm.

* `run_user` - The user who starts the execution.

* `recent_scan_time` - The latest scan time.
