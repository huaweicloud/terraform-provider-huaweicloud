---
subcategory: "Content Delivery Network (CDN)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cdn_cache_refresh"
description: |-
  Manages a CDN cache refresh resource within HuaweiCloud.
---

# huaweicloud_cdn_cache_refresh

Manages a CDN cache refresh resource within HuaweiCloud.

## Example Usage

```hcl
variable "urls" {
  type = list(string)
}

resource "huaweicloud_cdn_cache_refresh" "test" {
  urls = var.urls
}
```

## Argument Reference

The following arguments are supported:

* `urls` - (Required, List, NonUpdatable) Specifies the URLs that need to be refreshed.
  A URL must start with `http://` or `https://` and must contain the accelerated domain name.
  A URL can contain up to `4,096` characters. Enter up to `1,000` URLs or `100` directories and separate them by
  commas(,).
  + When `type` is set to **file**, the value should be file path. Example: `http://www.example.com/file01.html` or
    `http://www.example.com/`.
  + When `type` is set to **directory**, the value should be directory path. The URL must end with a slash (/).
    Example: `http://www.example.com/tt/ee/`.

* `enterprise_project_id` - (Optional, String, NonUpdatable) Specifies the ID of the enterprise project to which the
  resource belongs.  
  This parameter is only valid for enterprise users and is required when using sub-account.
  The value **all** represents all enterprise projects.

* `type` - (Optional, String, NonUpdatable) Specifies the refresh type.  
  The valid values are as follows:
  + **file**
  + **directory**

  Defaults to **file**.

* `mode` - (Optional, String, NonUpdatable) Specifies the directory refresh mode.  
  The valid values are as follows:
  + **all**: Refresh all resources in the directory.
  + **detect_modify_refresh**: Refresh changed resources in the directory.

  This field is valid only when `type` is set to **directory**. Defaults to **all**.

* `zh_url_encode` - (Optional, Bool, NonUpdatable) Specifies whether to encode Chinese characters in URLs before cache
  refresh.  
  The value **false** indicates disabled, and **true** indicates enabled. After enabled, cache is refreshed only for
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
