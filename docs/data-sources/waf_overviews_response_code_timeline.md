---
subcategory: "Web Application Firewall (WAF)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_waf_overviews_response_code_timeline"
description: |-
  Use this data source to query the WAF overviews response code timeline.
---

# huaweicloud_waf_overviews_response_code_timeline

Use this data source to query the WAF overviews response code timeline.

## Example Usage

```hcl
variable "start_time" {}
variable "end_time" {}

data "huaweicloud_waf_overviews_response_code_timeline" "test" {
  from = var.start_time
  to   = var.end_time
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `from` - (Required, Int) Specifies the query start time.

* `to` - (Required, Int) Specifies the query end time.

* `hosts` - (Optional, List) Specifies the ID list of the domain.

* `instances` - (Optional, List) Specifies the ID list of the dedicated WAF instances.

* `response_source` - (Optional, String) Specifies the response source.
  The valid values are **WAF** and **UPSTREAM**.

* `group_by` - (Optional, String) Specifies the display dimension.
  The value can be **DAY**, indicates data is displayed by the day.

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project ID.
  This parameter is only valid for enterprise users.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID in UUID format.

* `response_codes` - The response code timeline data for security statistics.

  The [response_codes](#response_codes_struct) structure is documented below.

<a name="response_codes_struct"></a>
The `response_codes` block supports:

* `key` - The response code.
  The valid values are as follows:
  + **ACCESS**: Total number of requests.
  + **CRAWLER**: Bot attack protection.
  + **ATTACK**: Total attack count.
  + **WEB_ATTACK**: Web basic protection.
  + **PRECISE**: Precision protection.
  + **CC**: CC attack protection.

* `timeline` - The statistics data over time for the corresponding response code.

  The [timeline](#response_codes_timeline_struct) structure is documented below.

<a name="response_codes_timeline_struct"></a>
The `timeline` block supports:

* `time` - The time point.

* `num` - The statistics data for the time range from the previous time point to the point specified by `time`.
