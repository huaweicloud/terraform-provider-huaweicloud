---
subcategory: "Host Security Service (HSS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_hss_host_manual_detection_status"
description: |-
  Use this data source to get the HSS host manual detection status within HuaweiCloud.
---

# huaweicloud_hss_host_manual_detection_status

Use this data source to get the HSS host manual detection status within HuaweiCloud.

## Example Usage

```hcl
variable "host_id" {}
variable "type" {}

data "huaweicloud_hss_host_manual_detection_status" "test" {
  host_id = var.host_id
  type    = var.type
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `host_id` - (Required, String) Specifies the host ID.

* `type` - (Required, String) Specifies the type of detection.  
  The valid values are as follows:
  + **pwd**: Weak password detection.

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project ID.  
  This parameter is valid only when the enterprise project is enabled.
  The default value is **0**, indicating the default enterprise project.
  If you need to query data for all enterprise projects, the value is **all_granted_eps**.
  If you only have permissions for a specific enterprise project, you need set the enterprise project ID. Otherwise,
  the operation may fail due to insufficient permissions.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID in UUID format.

* `scan_status` - The manual detection status.  
  The valid values are as follows:
  + **neverscan**: No scanning task.
  + **scanning**: Scanning.
  + **scaned**: Scan completed.
  + **failed**: Scan failed.
  + **over_time**: Scan timeout (scan time exceeds `30` minutes).
  + **longscanning**: Scan timeout (scan time exceeds `60` minutes).

* `scaned_time` - The detection completion time.
