---
subcategory: Content Delivery Network (CDN)
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cdn_domain_statistics"
description: |-
  Use this data source to get the statistics of CDN domain within HuaweiCloud.
---

# huaweicloud_cdn_domain_statistics

Use this data source to get the statistics of CDN domain within HuaweiCloud.

## Example Usage

```hcl
variable "domain_name" {}
variable "action" {}
variable "start_time" {}
variable "end_time" {}

data "huaweicloud_cdn_domain_statistics" "test" {
  domain_name = var.domain_name
  action      = var.action
  start_time  = var.start_time
  end_time    = var.end_time
  stat_type   = "req_num"
}
```

## Argument Reference

The following arguments are supported:

* `domain_name` - (Required, String) Specifies the domain name list.
  Domain names are separated by commas (,), for example, `www.test1.com,www.test2.com`.
  The value all indicates that all domain names under your account are queried.

* `stat_type` - (Required, String) Specifies the statistic type.

  For network resource consumption statistics, the value can be:
  + **bw**: bandwidth
  + **flux**: traffic.

  For access statistics, the value can be:
  + **req_num**: total number of requests.

  For HTTP status code statistics (one or more types can be returned), the value can be:
  + **http_code_2xx**: status codes 2xx.
  + **http_code_3xx**: status codes 3xx.
  + **http_code_4xx**: status codes 4xx.
  + **http_code_5xx**: status codes 5xx.
  + **status_code_2xx**: details of status code 2xx.
  + **status_code_3xx**: details of status code 3xx.
  + **status_code_4xx**: details of status code 4xx.
  + **status_code_5xx**: details of status code 5xx.

* `action` - (Required, String) Specifies the action name.  
  The valid values are as follows:
  + **location_summary**
  + **location_detail**

* `start_time` - (Required, Int) Specifies the start timestamp of the query.
  The timestamp must be set to a multiple of 5 minutes.
  + If the value of interval is `300`, set this parameter to a multiple of `5` minutes,
    for example, 1631240100000, which means 2021-09-10 10:15:00.
  + If the value of interval is `3,600`, set this parameter to a multiple of `1` hour,
    for example, 1631239200000, which means 2021-09-10 10:00:00.
  + If the value of interval is `86,400`, set this parameter to 00:00:00 (GMT+08:00),
    for example, 1631203200000, which means 2021-09-10 00:00:00.

* `end_time` - (Required, Int) Specifies the end timestamp of the query.
  The timestamp must be set to a multiple of `5` minutes.
  + If the value of interval is `300`, set this parameter to a multiple of `5` minutes,
    for example, 1631243700000, which means 2021-09-10 11:15:00.
  + If the value of interval is `3,600`, set this parameter to a multiple of `1` hour,
    for example, 1631325600000, which means 2021-09-11 10:00:00.
  + If the value of interval is `86,400`, set this parameter to 00:00:00 (GMT+08:00),
    for example, 1631376000000, which means 2021-09-12 00:00:00.

* `interval` - (Optional, Int) Specifies the query time interval, in seconds.  
  The vaild values are as follows:
  + **300**(`5` minutes): Maximum query span `2` days
  + **3,600**(`1` hour): Maximum query span `7` days
  + **86,400**(`1` day): Maximum query span `31` days

  The default is the minimum interval for the corresponding time span.

* `group_by` - (Optional, String) Specifies the data grouping mode. Use commas (,) to separate multiple groups.
  Available data groups are **domain**, **country**, **province**, and **isp**. By default, data is not grouped.

* `country` - (Optional, String) Specifies the country or region code. Use commas (,) to separate multiple codes.
  The value all indicates all country/region codes.
  See the [country and region](https://support.huaweicloud.com/intl/en-us/api-cdn/cdn_02_0089.html) for values.

* `province` - (Optional, String) Specifies the province code. This parameter is valid when country is set to **cn**.
  Use commas (,) to separate multiple codes. The value all indicates all provinces.
  See the [areas](https://support.huaweicloud.com/intl/en-us/api-cdn/cdn_02_0074.html) for values.

* `isp` - (Optional, String) Specifies the carrier code. Use commas (,) to separate multiple codes.
  The value all indicates all carrier codes.
  See the [carriers](https://support.huaweicloud.com/intl/en-us/api-cdn/cdn_02_0075.html) for values.

* `enterprise_project_id` - (Optional, String) Specifies the ID of the enterprise project to which the resource belongs.
  This parameter is valid only when the enterprise project function is enabled.
  The value all indicates all projects. This parameter is mandatory when you use an IAM user.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `result` - The data organized according to the specified grouping mode.
