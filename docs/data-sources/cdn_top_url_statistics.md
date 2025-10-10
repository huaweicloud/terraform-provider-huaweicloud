---
subcategory: "Content Delivery Network (CDN)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cdn_top_url_statistics"
description: |-
  Use this data source to get the TOP100 URL statistics of CDN domain within HuaweiCloud.
---

# huaweicloud_cdn_top_url_statistics

Use this data source to get the TOP100 URL statistics of CDN domain within HuaweiCloud.

## Example Usage

```hcl
variable "domain_name" {}
variable "start_timestamp" {}
variable "end_timestamp" {}

data "huaweicloud_cdn_top_url_statistics" "test" {
  domain_name  = var.domain_name
  start_time   = var.start_timestamp
  end_time     = var.end_timestamp
  stat_type    = "req_num"
  service_area = "mainland_china"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region where the CDN service is located.

* `start_time` - (Required, String) Specifies the start time of the query, in UTC format.  
  The timestamp must be set to 0:00:00 (GMT+08:00), for example, 2022-10-29 00:00:00.

* `end_time` - (Required, String) Specifies the end time of the query, in UTC format.  
  The timestamp must be set to 0:00:00 (GMT+08:00), for example, 2022-10-30 00:00:00.

* `domain_name` - (Required, String) Specifies the list of queried domain names.  
  Domain names are separated by commas (,), for example, `www.test1.com,www.test2.com`.
  The value all indicates that all domain names under your account are queried.

* `stat_type` - (Required, String) Specifies the query type.  
  The valid values are as follows:
  + **flux**: traffic (unit: Byte)
  + **req_num**: total number of requests

* `service_area` - (Optional, String) Specifies the service area of the query.  
  The valid values are as follows:
  + **mainland_china**: mainland China
  + **outside_mainland_china**: outside mainland China
  + **global**: global (default)

* `enterprise_project_id` - (Optional, String) Specifies the ID of the enterprise project to which the resources
  belong.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `service_area` - The service area.

* `top_url_summary` - The list of TOP100 URL statistics that matched filter parameters.  
  The [top_url_summary](#cdn_top_url_summary) structure is documented below.

<a name="cdn_top_url_summary"></a>
The `top_url_summary` block supports:

* `url` - The URL name.

* `value` - The value corresponding to the query type.

* `start_time` - The start time of the query, in UTC format.

* `end_time` - The end time of the query, in UTC format.

* `stat_type` - The query type.
