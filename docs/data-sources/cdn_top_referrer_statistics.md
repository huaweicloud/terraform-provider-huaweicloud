---
subcategory: "Content Delivery Network (CDN)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cdn_top_referrer_statistics"
description: |-
  Use this data source to get the TOP100 referrer statistics of CDN domain within HuaweiCloud.
---

# huaweicloud_cdn_top_referrer_statistics

Use this data source to get the TOP100 referrer statistics of CDN domain within HuaweiCloud.

-> The statistic data is obtained by scanning the service's offline logs and is subject
   to a delay of at least `6` hours.

## Example Usage

```hcl
variable "domain_name" {}
variable "start_time" {}
variable "end_time" {}

data "huaweicloud_cdn_top_referrer_statistics" "test" {
  domain_name = var.domain_name
  start_time  = var.start_time
  end_time    = var.end_time
  stat_type   = "req_num"
}
```

## Argument Reference

The following arguments are supported:

* `start_time` - (Required, String) Specifies the start time of the query range, in RFC3339 format.  
  The time must be set to twelve o'clock in the evening, for example, **2022-10-29T00:00:00Z**.

* `end_time` - (Required, String) Specifies the end time of the query range, in RFC3339 format.  
  The time must be set to twelve o'clock in the evening, for example, **2022-10-30T00:00:00Z**.

* `domain_name` - (Required, String) Specifies the list of queried domain names.  
  Domain names are separated by the comma (,) character, for example, **"www.test1.com,www.test2.com"**.
  The value all indicates that all domain names under your account are queried.

* `stat_type` - (Required, String) Specifies the statistical type of the query.  
  The valid values are as follows:
  + **flux**: traffic (unit: Byte)
  + **req_num**: total number of requests

* `service_area` - (Optional, String) Specifies the service area of the query.  
  The valid values are as follows:
  + **mainland_china**: mainland China
  + **outside_mainland_china**: outside mainland China
  + **global**: global

  Defaults to **global**.

* `enterprise_project_id` - (Optional, String) Specifies the ID of the enterprise project to which the resource
  belongs.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `statistics` - The list of TOP100 referrer statistics that matched filter parameters.  
  The [statistics](#cdn_top_refer_statistics) structure is documented below.

<a name="cdn_top_refer_statistics"></a>
The `statistics` block supports:

* `refer` - The referrer value.

* `value` - The value corresponding to the query type.
