---
subcategory: "Host Security Service (HSS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_hss_asset_process_detail"
description: |-
  Use this data source to get the list of servers for the specified process.
---

# huaweicloud_hss_asset_process_detail

Use this data source to the list of servers for the specified process.

## Example Usage

```hcl
data "huaweicloud_hss_asset_process_detail" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `host_name` - (Optional, String) Specifies the host name.

* `host_ip` - (Optional, String) Specifies the host IP address.

* `path` - (Optional, String) Specifies the process executable file path.

* `category` - (Optional, String) Specifies the type.
  The valid values are as follows:
  + **host**: Host.
  + **container**: Container.

* `enterprise_project_id` - (Optional, String) Specifies the ID of the enterprise project to which the resource belongs.
  This parameter is valid only when the enterprise project function is enabled.
  The value **all_granted_eps** indicates all enterprise projects.
  If omitted, the default enterprise project will be used.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `data_list` - The host statistics list.
  The [data_list](#process_detail_struct) structure is documented below.

<a name="process_detail_struct"></a>
The `data_list` block supports:

* `hash` - The SHA256 value corresponding to the path.

* `host_ip` - The host IP address.

* `host_name` - The host name.

* `launch_params` - The startup parameter.

* `launch_time` - The startup time.

* `process_path` - The process executable file path.

* `process_pid` - The process PID.

* `run_permission` - The file permission.

* `container_id` - The container ID.

* `container_name` - The container name.
