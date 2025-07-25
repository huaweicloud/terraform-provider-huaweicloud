---
subcategory: "Host Security Service (HSS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_hss_antivirus_statistic"
description: |-
  Use this data source to get the list of HSS antivirus statistic within HuaweiCloud.
---
# huaweicloud_hss_antivirus_statistic

Use this data source to get the list of HSS antivirus statistic within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_hss_antivirus_statistic" "test" {}
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

* `total_malware_num` - The total number of viruses.

* `malware_host_num` - Affects the number of hosts.

* `total_task_num` - Accumulated number of scanning tasks.

* `scanning_task_num` - The number of running tasks.

* `latest_scan_time` - The start time in milliseconds.

* `scan_type` - The scan type.  
  The valid values are as follows:
  + **quick**: Quick scan.
  + **full**: Full scan.
  + **custom**: Custom scan.
