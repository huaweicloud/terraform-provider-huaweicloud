---
subcategory: "Host Security Service (HSS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_hss_overview_protection_statistics"
description: |-
  Use this data source to get the overview protection statistics of HSS within HuaweiCloud.
---

# huaweicloud_hss_overview_protection_statistics

Use this data source to get the overview protection statistics of HSS within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_hss_overview_protection_statistics" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project ID.
  The default value is **0**, indicating the default enterprise project.
  If you need to query data for all enterprise projects, the value is **all_granted_eps**.
  If you only have permissions for a specific enterprise project, you need set the enterprise project ID. Otherwise,
  the operation may fail due to insufficient permissions.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID in UUID format.

* `vul_library_update_time` - The vulnerability library update time.

* `protect_days` - The protection days.

* `threat_library_update_time` - The threat library update time.

* `vul_detected_total_num` - The total number of detected vulnerabilities.

* `baseline_detected_total_num` - The total number of detected baseline items.

* `finger_scan_total_num` - The total number of scanned fingerprints.

* `alarm_detected_total_num` - The total number of detected alarms.

* `ransomware_alarm_detected_total_num` - The total number of detected ransomware alarms.

* `file_alarm_detected_total_num` - The total number of detected file alarms.

* `rasp_alarm_detected_total_num` - The total number of detected RASP alarms.

* `wtp_alarm_detected_total_num` - The total number of detected WTP alarms.

* `image_risk_total_num` - The total number of image risks.

* `container_alarm_total_num` - The total number of container alarms.

* `container_firewall_policy_total_num` - The total number of container firewall policies.

* `auto_kill_virus_status` - The status of automatic virus killing. The value can be **true** or **false**.
