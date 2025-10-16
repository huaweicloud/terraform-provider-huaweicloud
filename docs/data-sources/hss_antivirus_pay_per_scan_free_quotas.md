---
subcategory: "Host Security Service (HSS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_hss_antivirus_pay_per_scan_free_quotas"
description: |-
  Use this data source to get the HSS antivirus pay per scan free quotas within HuaweiCloud.
---

# huaweicloud_hss_antivirus_pay_per_scan_free_quotas

Use this data source to get the HSS antivirus pay per scan free quotas within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_hss_antivirus_pay_per_scan_free_quotas" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `enterprise_project_id` - (Optional, String, NonUpdatable) Specifies the enterprise project ID.  
  This parameter is valid only when the enterprise project is enabled.
  The default value is **0**, indicating the default enterprise project.
  If it is necessary to operate the asset under all enterprise projects, the value is **all_granted_eps**.
  If you only have permissions for a specific enterprise project, you need set the enterprise project ID. Otherwise,
  the operation may fail due to insufficient permissions.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID in UUID format.

* `free_quota` - The free quota for antivirus scanning.

* `occupied_free_quota` - The occupied free quota for current scanning tasks.

* `antivirus_already_given` - Whether the antivirus interface has already given free quota.

* `antivirus_free_quota` - The free quota that should be given by the antivirus interface.

* `report_already_given` - Whether the monthly report interface has already given free quota.

* `report_free_quota` - The free quota that should be given by the monthly report interface.
