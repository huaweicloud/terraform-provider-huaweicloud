---
subcategory: Content Delivery Network (CDN)
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cdn_cache_url_tasks"
description: |-
  Use this data source to get the list of CDN cache url tasks within HuaweiCloud.
---

# huaweicloud_cdn_cache_url_tasks

Use this data source to get the list of CDN cache url tasks within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_cdn_cache_url_tasks" "test" {}
```

## Argument Reference

The following arguments are supported:

* `start_time` - (Optional, Int) Specifies the start timestamp, in milliseconds.  
  The default value is 00:00 of the current day.

* `end_time` - (Optional, Int) Specifies the end timestamp, in milliseconds.  
  The default value is 00:00 of the next day.

* `url` - (Optional, String) Specifies the refresh or preheat URL.

* `task_type` - (Optional, String) Specifies the task type.  
  The valid values are as follows:
  + **REFRESH**: cache refresh.
  + **PREHEATING**: cache preheat.

* `status` - (Optional, String) Specifies the URL status.  
  The valid values are as follows:
  + **processing**
  + **succeed**
  + **failed**
  + **waiting**
  + **refreshing**
  + **preheating**

* `file_type` - (Optional, String) Specifies the file type.
  The valid values are as follows:
  + **file**
  + **directory**

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `tasks` - The list of URL task information that matched filter parameters.
  The [tasks](#cdn_cache_url_tasks) structure is documented below.

<a name="cdn_cache_url_tasks"></a>
The `tasks` block supports:

* `id` - The URL ID.

* `url` - The URL.

* `status` - The URL status.
  + **processing**
  + **succeed**
  + **failed**
  + **waiting**
  + **refreshing**
  + **preheating**

* `task_type` - The task type.
  + **REFRESH**: cache refresh.
  + **PREHEATING**: cache preheat.

* `mode` - The directory refresh mode.
  + **all**: refresh all resources in the directory.
  + **detect_modify_refresh**: refresh changed resources in the directory.

* `task_id` - The task ID.

* `modify_time` - The modification time, in RFC3339 format.

* `created_at` - The creation time, in RFC3339 format.

* `file_type` - The file type.
  + **file**
  + **directory**
