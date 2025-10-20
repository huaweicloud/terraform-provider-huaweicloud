---
subcategory: "Host Security Service (HSS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_hss_files_change_hosts"
description: |-
  Use this data source to get the modify server list.
---

# huaweicloud_hss_files_change_hosts

Use this data source to get the modify server list.

## Example Usage

```hcl
data "huaweicloud_hss_files_change_hosts" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `begin_time` - (Optional, Int) Specifies the query start time.
  The format is 13-digit timestamp in millisecond.

* `end_time` - (Optional, Int) Specifies the query end time.
  The format is 13-digit timestamp in millisecond.

* `host_name` - (Optional, String) Specifies the server name.

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project ID to which the hosts belong.
  This parameter is valid only when the enterprise project is enabled.
  The default value is **0**, indicating the default enterprise project.
  If you need to query data for all enterprise projects, the value is **all_granted_eps**.
  If you only have permissions for a specific enterprise project, you need set the enterprise project ID. Otherwise,
  the operation may fail due to insufficient permissions.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `data_list` - The list of change servers.
  The [data_list](#data_list_struct) structure is documented below.

<a name="data_list_struct"></a>
The `data_list` block supports:

* `host_name` - The server name.

* `host_id` - The server ID.

* `change_total_num` - The total number of changes.

* `change_file_num` - The change the number of files.

* `change_registry_num` - The change the number of registry entries.

* `latest_time` - The last modified time.
