---
subcategory: "Host Security Service (HSS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_hss_ransomware_backup_operation_logs"
description: |-
  Use this data source to get the list of HSS ransomware backup operation logs within HuaweiCloud.
---

# huaweicloud_hss_ransomware_backup_operation_logs

Use this data source to get the list of HSS ransomware backup operation logs within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_hss_ransomware_backup_operation_logs" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project ID.  
  The default value is **0**, indicating the default enterprise project.
  If it is necessary to operate the asset under all enterprise projects, the value is **all_granted_eps**.
  If you only have permissions for a specific enterprise project, you need set the enterprise project ID. Otherwise,
  the operation may fail due to insufficient permissions.

* `status` - (Optional, String) Specifies the restore status.  
  The valid values are as follows:
  + **success**: Success.
  + **skipped**: Skipped.
  + **failed**: Failed.
  + **running**: Running.
  + **timeout**: Timeout.
  + **waiting**: Waiting.

* `resource_name` - (Optional, String) Specifies the server name.

* `last_days` - (Optional, Int) Specifies the query time range.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID in UUID format.

* `total_num` - The total number.

* `data_list` - The list of operation logs.

  The [data_list](#data_list_struct) structure is documented below.

<a name="data_list_struct"></a>
The `data_list` block supports:

* `host_id` - The host ID.

* `host_name` - The host name.

* `backup_name` - The backup name.

* `process` - The restore progress (percentage).

* `status` - The restore status.  
  The valid values are as follows:
  + **success**: Success.
  + **skipped**: Skipped.
  + **failed**: Failed.
  + **running**: Running.
  + **timeout**: Timeout.
  + **waiting**: Waiting.

* `started_at` - The task start time.

* `ended_at` - The task end time.

* `error_info` - The failure information.

  The [error_info](#error_info_struct) structure is documented below.

<a name="error_info_struct"></a>
The `error_info` block supports:

* `code` - The error code.

* `message` - The error message.
