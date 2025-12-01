---
subcategory: "Host Security Service (HSS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_hss_hosts_risk"
description: |-
  Use this data source to get the list of HSS hosts risk within HuaweiCloud.
---

# huaweicloud_hss_hosts_risk

Use this data source to get the list of HSS hosts risk within HuaweiCloud.

## Example Usage

```hcl
variable "host_id_list" {
  type = list(string)
}

data "huaweicloud_hss_hosts_risk" "test" {
  host_id_list = var.host_id_list
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `host_id_list` - (Required, List) Specifies the host ID list.

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project ID.  
  The default value is **0**, indicating the default enterprise project.
  If you need to query data for all enterprise projects, the value is **all_granted_eps**.
  If you only have permissions for a specific enterprise project, you need set the enterprise project ID. Otherwise,
  the operation may fail due to insufficient permissions.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID in UUID format.

* `total_num` - The total number.

* `data_list` - The list of hosts risk.

  The [data_list](#data_list_struct) structure is documented below.

<a name="data_list_struct"></a>
The `data_list` block supports:

* `host_id` - The host ID.

* `agent_status` - The agent status.  
  The valid values are as follows:
  + **not_installed**
  + **online**
  + **offline**
  + **install_failed**
  + **installing**

* `install_result_code` - The installation result.  
  The valid values are as follows:
  + **install_succeed**
  + **network_access_timeout**
  + **invalid_port**
  + **auth_failed**
  + **permission_denied**
  + **no_available_vpc**
  + **install_exception**
  + **invalid_param**
  + **install_failed**
  + **package_unavailable**
  + **os_type_not_support**
  + **os_arch_not_support**

* `version` - The version.  
  The valid values are as follows:
  + **hss.version.null**
  + **hss.version.basic**
  + **hss.version.advanced**
  + **hss.version.enterprise**
  + **hss.version.premium**
  + **hss.version.wtp**
  + **hss.version.container.enterprise**

* `protect_status` - The protection status.  
  The valid values are as follows:
  + **closed**
  + **opened**

* `detect_result` - The cloud host security testing results.  
  The valid values are as follows:
  + **undetected**
  + **clean**
  + **risk**
  + **scanning**

* `asset` - The asset risk.

* `vulnerability` - The vulnerability risk.

* `baseline` - The baseline risk.

* `intrusion` - The intrusion risk.
