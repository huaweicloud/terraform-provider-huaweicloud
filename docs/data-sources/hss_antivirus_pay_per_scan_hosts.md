---
subcategory: "Host Security Service (HSS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_hss_antivirus_pay_per_scan_hosts"
description: |-
  Use this data source to get the list of HSS antivirus pay-per-scan hosts within HuaweiCloud.
---

# huaweicloud_hss_antivirus_pay_per_scan_hosts

Use this data source to get the list of HSS antivirus pay-per-scan hosts within HuaweiCloud.

## Example Usage

```hcl
variable "scan_type" {}
variable "start_type" {}

data "huaweicloud_hss_antivirus_pay_per_scan_hosts" "test" {
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

* `start_type` - (Required, String) Specifies the start type.  
  The valid values are as follows:
  + **now**: Start immediately.
  + **period**: Periodic start.

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project ID.  
  This parameter is valid only when the enterprise project is enabled.
  The default value is **0**, indicating the default enterprise project.
  If you need to query data for all enterprise projects, the value is **all_granted_eps**.
  If you only have permissions for a specific enterprise project, you need to set the enterprise project ID. Otherwise,
  the operation may fail due to insufficient permissions.

* `host_id` - (Optional, String) Specifies the host ID.

* `host_name` - (Optional, String) Specifies the host name.

* `private_ip` - (Optional, String) Specifies the private IP.

* `public_ip` - (Optional, String) Specifies the public IP.

* `group_id` - (Optional, String) Specifies the group ID.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID in UUID format.

* `data_list` - The list of hosts.

  The [data_list](#data_list_struct) structure is documented below.

<a name="data_list_struct"></a>
The `data_list` block supports:

* `host_id` - The host ID.

* `host_name` - The host name.

* `group_id` - The group ID.

* `public_ip` - The public IP address.

* `private_ip` - The private IP address.

* `agent_id` - The agent ID.

* `os_type` - The operating system type. The valid values are **Linux** and **Windows**.
