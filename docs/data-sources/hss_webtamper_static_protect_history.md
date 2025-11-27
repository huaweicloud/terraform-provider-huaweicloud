---
subcategory: "Host Security Service (HSS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_hss_webtamper_static_protect_history"
description: |-
  Use this data source to get the list of static WTP events.
---

# huaweicloud_hss_webtamper_static_protect_history

Use this data source to get the list of static WTP events.

## Example Usage

```hcl
variable "start_time" {}
variable "end_time" {}

data "huaweicloud_hss_webtamper_static_protect_history" "test" {
  start_time = var.start_time
  end_time   = var.end_time
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `start_time` - (Required, Int) Specifies the query start time, in milliseconds.
  The start time cannot be earlier than `30` days ago. If an earlier time is specified,
  the query range is still the past `30` days.

* `end_time` - (Required, Int) Specifies the query end time, in milliseconds.
  The end time cannot be earlier than the start time. The difference between the end time and start time
  cannot exceed `30` days. If it exceeds `30` days, the end time is set to one day later than the start time.

* `host_id` - (Optional, String) Specifies the server ID.

* `host_name` - (Optional, String) Specifies the server name.

* `host_ip` - (Optional, String) Specifies the server EIP.

* `file_path` - (Optional, String) Specifies the protection file path.

* `file_operation` - (Optional, String) Specifies the file operation.
  The valid values are as follows:
  + **add**
  + **delete**
  + **modify**
  + **attribute**

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project ID to which the hosts belong.  
  This parameter is valid only when the enterprise project is enabled.
  The default value is **0**, indicating the default enterprise project.
  If you need to query data for all enterprise projects, the value is **all_granted_eps**.
  If you only have permissions for a specific enterprise project, you need set the enterprise project ID. Otherwise,
  the operation may fail due to insufficient permissions.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `data_list` - The list of static WTP events.
  The [data_list](#static_data_list) structure is documented below.

<a name="static_data_list"></a>
The `data_list` block supports:

* `host_name` - The server name.

* `occur_time` - The detection time, in milliseconds.

* `file_path` - The protection file path.

* `file_operation` - The file operation.

* `host_ip` - The server EIP.

* `process_id` - The process ID.

* `process_name` - The process name.

* `process_cmd` - The process command line.
