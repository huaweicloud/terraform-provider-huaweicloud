---
subcategory: "Host Security Service (HSS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_hss_antivirus_available_hosts"
description: |-
  Use this data source to get the list of HSS antivirus available hosts within HuaweiCloud.
---
# huaweicloud_hss_antivirus_available_hosts

Use this data source to get the list of HSS antivirus available hosts within HuaweiCloud.

## Example Usage

```hcl
variable "scan_type" {}
variable "start_type" {}

data "huaweicloud_hss_antivirus_available_hosts" "test" {
  scan_type  = var.scan_type
  start_type = var.start_type
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `scan_type` - (Required, String) Specifies the scan type.  
  The valid values are as follows:
  + **quick**: Quick scan.
  + **full**: Full scan.
  + **custom**: Custom scan.

* `start_type` - (Required, String) Specifies the startup type.  
  The valid values are as follows:
  + **now**
  + **period**

* `host_id` - (Optional, String) Specifies the host ID.

* `host_name` - (Optional, String) Specifies the host name.

* `private_ip` - (Optional, String) Specifies the host private IP address.

* `public_ip` - (Optional, String) Specifies the host public IP address.

* `group_id` - (Optional, String) Specifies the host group ID.

* `policy_id` - (Optional, String) Specifies the policy ID.

* `next_start_time` - (Optional, Int) Specifies the next startup time in milliseconds.

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project ID.  
  This parameter is valid only when the enterprise project is enabled.
  The default value is **0**, indicating the default enterprise project.
  If you need to query data for all enterprise projects, the value is **all_granted_eps**.
  If you only have permissions for a specific enterprise project, you need set the enterprise project ID. Otherwise,
  the operation may fail due to insufficient permissions.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID in UUID format.

* `total_num` - The total number of available hosts.

* `data_list` - The list of available hosts details.
  The [data_list](#data_list_struct) structure is documented below.

<a name="data_list_struct"></a>
The `data_list` block supports:

* `host_id` - The host ID.

* `host_name` - The host name.

* `public_ip` - The host public IP address.

* `private_ip` - The host private IP address.

* `agent_id` - The host agent ID.

* `os_type` - The host operating system type. The valid value can be **Linux** or **Windows**.

* `group_id` - The host group ID.
