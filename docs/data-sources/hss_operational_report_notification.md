---
subcategory: "Host Security Service (HSS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_hss_operational_report_notification"
description: |-
  Use this data source to get the operational report notification of HSS within HuaweiCloud.
---

# huaweicloud_hss_operational_report_notification

Use this data source to get the operational report notification of HSS within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_hss_operational_report_notification" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID in UUID format.

* `status` - The status of the notification.  
  The valid values are as follows:
  + **close**: Do not display the popup.
  + **open**: Display the popup.

* `title` - The title of the notification.  
  The valid values are as follows:
  + **vul-fix-master**: Patch Master
  + **vul-fix-expert**: Vulnerability Fix Expert
  + **baseline-fix**: Risk Configuration Fix Expert
  + **malware-file**: Anti-virus Pioneer
  + **ransomware-event**: Anti-ransomware Expert
  + **web-tamper-event**: Website Guard

* `report_id` - The unique identifier of the current user, report time, returned as a string: yyyy-MM.

* `current_month` - The current month.
