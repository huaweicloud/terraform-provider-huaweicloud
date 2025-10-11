---
subcategory: "Host Security Service (HSS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_hss_files_statistic"
description: |-
  Use this data source to get the server files statistic.
---

# huaweicloud_hss_files_statistic

Use this data source to get the server files statistic.

## Example Usage

```hcl
data "huaweicloud_hss_files_statistic" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `begin_time` - (Optional, Int) Specifies the query start time.
  The format is 13-digit timestamp in millisecond.

* `end_time` - (Optional, Int) Specifies the query end time.
  The format is 13-digit timestamp in millisecond.

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project ID to which the hosts belong.
  This parameter is valid only when the enterprise project is enabled.
  The default value is **0**, indicating the default enterprise project.
  If you need to query data for all enterprise projects, the value is **all_granted_eps**.
  If you only have permissions for a specific enterprise project, you need set the enterprise project ID. Otherwise,
  the operation may fail due to insufficient permissions.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `host_num` - The total number of servers.

* `change_total_num` - The total number of changes.

* `change_file_num` - The change the number of files.

* `change_registry_num` - The change the number of registry entries.

* `modify_num` - The modification quantity.

* `add_num` - The add quantity.

* `delete_num` - The deletion quantity.
