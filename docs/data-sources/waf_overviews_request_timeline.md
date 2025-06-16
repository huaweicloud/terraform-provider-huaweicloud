---
subcategory: "Web Application Firewall (WAF)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_waf_overviews_request_timeline"
description: |-
  Use this data source to query the website requests.
---

# huaweicloud_waf_overviews_request_timeline

Use this data source to query the website requests.

## Example Usage

```hcl
variable "start_time" {}
variable "end_time" {}

data "huaweicloud_waf_overviews_request_timeline" "test" {
  from = var.start_time
  to   = var.end_time
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `from` - (Required, Int) Specifies the query start time.
  The format is 13-digit timestamp in millisecond.

* `to` - (Required, Int) Specifies the query end time.
  The format is 13-digit timestamp in millisecond.

* `group_by` - (Optional, String) Specifies the display dimension.
  The value can be **DAY**, indicates data is displayed by the day. Defaults display by minutes.

* `hosts` - (Optional, List) Specifies the ID list of the domain.

* `instances` - (Optional, List) Specifies the ID list of the dedicated WAF instances.

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project ID.
  If you want to query resources under all enterprise projects, set this parameter to **all_granted_eps**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `requests` - The request timeline data for security statistics.

  The [requests](#requests_struct) structure is documented below.

<a name="requests_struct"></a>
The `requests` block supports:

* `key` - The key type.
  The options are **ACCESS** for total requests, **CRAWLER** for bot mitigation, **ATTACK** for total attacks,
  **WEB_ATTACK** for basic web protection, **PRECISE** for precise protection, and **CC** for CC attack protection.

* `timeline` - The statistics data over time for the corresponding key.

  The [timeline](#requests_timeline_struct) structure is documented below.

<a name="requests_timeline_struct"></a>
The `timeline` block supports:

* `time` - The time point.

* `num` - The statistics data for the time range from the previous time point to the point specified by `time`.
