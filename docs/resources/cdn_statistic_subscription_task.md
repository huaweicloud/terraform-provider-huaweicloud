---
subcategory: "Content Delivery Network (CDN)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cdn_statistic_subscription_task"
description: |-
  Manages a CDN statistic subscription task resource within HuaweiCloud.
---

# huaweicloud_cdn_statistic_subscription_task

Manages a CDN statistic subscription task resource within HuaweiCloud.

## Example Usage

```hcl
variable "task_name" {}

resource "huaweicloud_cdn_statistic_subscription_task" "test" {
  name        = var.task_name
  period_type = 0
  emails      = "user1@example.com,user2@example.com"
  domain_name = "www.example.com,www.example2.com"
  report_type = "0,1,2,3,4,5"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required, String) Specifies the name of the subscription task.  
  The name can contain word characters, hyphens, and Chinese characters, with a maximum length of 32 characters.

* `period_type` - (Required, Int) Specifies the type of the subscription task.  
  The valid values are as follows:
  + **0**: Daily report
  + **1**: Weekly report
  + **2**: Monthly report

* `emails` - (Required, String) Specifies the email addresses to receive the operation reports.  
  Multiple email addresses are separated by commas (,).

* `domain_name` - (Required, String) Specifies the list of domain names to subscribe.  
  Multiple domain names are separated by commas (,).  
  If set to **all**, all domain names under the account will be subscribed.

* `report_type` - (Required, String) Specifies the type of the operation report.  
  Multiple report types are separated by commas (,).  
  The valid values are as follows:
  + **0**: Access area distribution
  + **1**: Country distribution
  + **2**: Carrier distribution
  + **3**: Domain ranking (by traffic)
  + **4**: Popular URLs (by traffic)
  + **5**: Popular URLs (by request count)
  + **6**: Popular Referer (by traffic)
  + **7**: Popular Referer (by request count)
  + **10**: Origin popular URLs (by traffic)
  + **11**: Origin popular URLs (by request count)
  + **13**: Popular UA (by traffic)
  + **14**: Popular UA (by request count)

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `create_time` - The creation time of the subscription task, in RFC3339 format.

* `update_time` - The last update time of the subscription task, in RFC3339 format.

## Import

The subscription task can be imported using the `id` or `name`, e.g.

```bash
$ terraform import huaweicloud_cdn_statistic_subscription_task.test <id>
```

or

```bash
$ terraform import huaweicloud_cdn_statistic_subscription_task.test <name>
```
