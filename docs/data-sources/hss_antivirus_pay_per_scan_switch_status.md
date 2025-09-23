---
subcategory: "Host Security Service (HSS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_hss_antivirus_pay_per_scan_switch_status"
description: |-
  Use this data source to get the HSS antivirus pay per scan switch status within HuaweiCloud.
---
# huaweicloud_hss_antivirus_pay_per_scan_switch_status

Use this data source to get the HSS antivirus pay per scan switch status within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_hss_antivirus_pay_per_scan_switch_status" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project ID.  
  This parameter is valid only when the enterprise project is enabled.
  The default value is **0**, indicating the default enterprise project.
  If you need to query data for all enterprise projects, the value is **all_granted_eps**.
  If you only have permissions for a specific enterprise project, you need set the enterprise project ID. Otherwise,
  the operation may fail due to insufficient permissions.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID in UUID format.

* `enabled` - Is it available. The valid value can be **true** or **false**.
