---
subcategory: "Host Security Service (HSS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_hss_antivirus_handle_history"
description: |-
  Use this data source to get the list of antivirus handle history.
---

# huaweicloud_hss_antivirus_handle_history

Use this data source to get the list of antivirus handle history.

## Example Usage

```hcl
data "huaweicloud_hss_antivirus_handle_history" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `malware_name` - (Optional, String) Specifies the virus name.

* `file_path` - (Optional, String) Specifies the file path.

* `severity_list` - (Optional, List) Specifies the threat level.
  The valid values are as follows:
  + **Low**: Low risk.
  + **Medium**: Medium risk.
  + **High**: High risk.
  + **Critical**: Critical risk.

* `host_name` - (Optional, String) Specifies the server name.

* `private_ip` - (Optional, String) Specifies the server private IP.

* `public_ip` - (Optional, String) Specifies the server public IP.

* `asset_value` - (Optional, String) Specifies the asset importance.
  The valid values are as follows:
  + **important**: Important assets.
  + **common**: Common assets.
  + **test**: Test assets.

* `handle_method` - (Optional, String) Specifies the handling status to filter the results.  
  The valid values are as follows:
  + **mark_as_handled**: Manual handling.
  + **ignore**: Ignore.
  + **add_to_alarm_whitelist**: Add to alarm whitelist.
  + **isolate_and_kill**: Isolate file.
  + **unhandled**: Cancel manual handling.
  + **do_not_ignore**: Unignore.
  + **remove_from_alarm_whitelist**: Remove from the alarm whitelist.
  + **do_not_isolate_or_kill**: Cancel isolation of a file.

* `user_name` - (Optional, String) Specifies the user name.

* `event_type` - (Optional, Int) Specifies the event type.

* `sort_dir` - (Optional, String) Specifies sorting order.
  The valid values are as follows:
  + **asc**: Ascending order.
  + **desc**: Descending order.

* `sort_key` - (Optional, String) Specifies sorting field.
  The valid value is **handle_time**.

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project ID.  
  This parameter is valid only when the enterprise project is enabled.
  The default value is **0**, indicating the default enterprise project.
  If you need to query data for all enterprise projects, the value is **all_granted_eps**.
  If you only have permissions for a specific enterprise project, you need to set the enterprise project ID. Otherwise,
  the operation may fail due to insufficient permissions.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `data_list` - The list of virus handle history.

  The [data_list](#data_list_struct) structure is documented below.

<a name="data_list_struct"></a>
The `data_list` block supports:

* `result_id` - The result ID of virus scanning and removal.

* `malware_type` - The virus type.

* `malware_name` - The virus name.

* `severity` - The threat level.

* `file_path` - The file path.

* `host_name` - The server name.

* `private_ip` - The server private IP address.

* `public_ip` - The server public IP address.

* `asset_value` - The asset importance.

* `occur_time` - The occurrence time, in milliseconds.

* `handle_status` - The handling status.

* `handle_method` - The handling method.

* `notes` - The remarks.

* `handle_time` - The handing time.

* `user_name` - The user name.
