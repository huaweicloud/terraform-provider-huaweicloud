---
subcategory: "Host Security Service (HSS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_hss_baseline_scan_status"
description: |-
  Use this data source to get the manual baseline scan results.
---

# huaweicloud_hss_baseline_scan_status

Use this data source to get the manual baseline scan results.

## Example Usage

```hcl
data "huaweicloud_hss_baseline_scan_status" "test" {}
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

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `scan_status` - The scan status.
  The valid values are as follows:
  + **neverscan**
  + **scanning**
  + **scanned**
  + **failed**

* `scanned_time` - The scan end time, in milliseconds.
