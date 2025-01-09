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

* `urls` - (Required, List, ForceNew) Specifies the URLs that need to be preheated.
  A URL must start with `http://` or `https://` and must contain the accelerated domain name.
  A URL can contain up to `4,096` characters. Enter up to `1,000` URLs and separate them by commas (,).

  Changing this parameter will create a new resource.

* `enterprise_project_id` - (Optional, String, ForceNew) Specifies the enterprise project ID to which the accelerated
  domain name belongs. This parameter is only valid for enterprise users and is required when using Sub-account.
  The value **all** represents all enterprise projects.

  Changing this parameter will create a new resource.

* `zh_url_encode` - (Optional, Bool, ForceNew) Specifies whether to encode Chinese characters in URLs before cache preheat.
  The value **false** indicates disabled, and **true** indicates enabled. After enabled, cache is preheated only for
  transcode URLs. Defaults to **false**.

  Changing this parameter will create a new resource.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `status` - The task execution result. Possible values: **task_done** (successful) and **task_inprocess** (processing).

* `created_at` - The creation time.

* `processing` - The number of URLs that are being processed.

* `succeed` - The number of URLs processed.

* `failed` - The number of URLs that failed to be processed.

* `total` - The total number of URLs in historical tasks.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 5 minutes.
