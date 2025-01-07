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

* `start_time` - (Optional, Int) Specifies the start timestamp, in milliseconds. The default value is 00:00 of the
  current day.

* `end_time` - (Optional, Int) Specifies the end timestamp, in milliseconds. The default value is 00:00 of the next day.

* `url` - (Optional, String) Specifies the refresh or preheat URL.

* `task_type` - (Optional, String) Specifies the task type. Possible values: **REFRESH** (cache refresh) and
  **PREHEATING** (cache preheat).

* `status` - (Optional, String) Specifies the URL status. Possible values: **processing**, **succeed**, **failed**,
  **waiting**, **refreshing**, and **preheating**.

* `file_type` - (Optional, String) Specifies the file type. Possible values: **file** and **directory**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `tasks` - The list of URL task information. The [tasks](#CacheUrlTasks_tasks) structure is documented below.

<a name="CacheUrlTasks_tasks"></a>
The `tasks` block supports:

* `id` - Indicates the URL ID.

* `url` - Indicates the URL.

* `status` - Indicates the URL status. Possible values: **processing**, **succeed**, **failed**, **waiting**,
  **refreshing**, and **preheating**.

* `task_type` - Indicates the task type. Possible values: **REFRESH** (cache refresh) and **PREHEATING** (cache preheat).

* `mode` - Indicates the directory refresh mode. Possible values: **all** (refresh all resources in the directory) and
  **detect_modify_refresh** (refresh changed resources in the directory).

* `task_id` - Indicates the task ID.

* `modify_time` - Indicates the modification time.

* `created_at` - Indicates the creation time.

* `file_type` - Indicates the file type. Possible values: **file** and **directory**.
