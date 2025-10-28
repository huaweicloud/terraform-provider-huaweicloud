---
subcategory: "Content Delivery Network (CDN)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cdn_cache_preheat"
description: |-
  Manages a CDN cache preheat resource within HuaweiCloud.
---

# huaweicloud_cdn_cache_preheat

Manages a CDN cache preheat resource within HuaweiCloud.

## Example Usage

```hcl
variable "urls" {
  type = list(string)
}

resource "huaweicloud_cdn_cache_preheat" "test" {
  urls = var.urls
}
```

## Argument Reference

The following arguments are supported:

* `urls` - (Required, List, NonUpdatable) Specifies the URLs that need to be preheated.
  A URL must start with `http://` or `https://` and must contain the accelerated domain name.
  A URL can contain up to `4,096` characters. Enter up to `1,000` URLs and separate them by commas (,).

* `enterprise_project_id` - (Optional, String, NonUpdatable) Specifies the ID of the enterprise project to which the
  resource belongs.  
  This parameter is only valid for enterprise users and is required when using sub-account.
  The value **all** represents all enterprise projects.

* `zh_url_encode` - (Optional, Bool, NonUpdatable) Specifies whether to encode Chinese characters in URLs before cache
  preheat.  
  The value **false** indicates disabled, and **true** indicates enabled. After enabled, cache is preheated only for
  transcode URLs. Defaults to **false**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `status` - The task execution result.
  + **task_done**: successful.
  + **task_inprocess**: processing.

* `created_at` - The creation time, in RFC3339 format.

* `processing` - The number of URLs that are being processed.

* `succeed` - The number of URLs processed.

* `failed` - The number of URLs that failed to be processed.

* `total` - The total number of URLs in historical tasks.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 5 minutes.
