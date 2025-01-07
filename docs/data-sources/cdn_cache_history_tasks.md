---
subcategory: Content Delivery Network (CDN)
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cdn_cache_history_tasks"
description: |-
  Use this data source to get the list of CDN cache history tasks within HuaweiCloud.
---

# huaweicloud_cdn_cache_history_tasks

Use this data source to get the list of CDN cache history tasks within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_cdn_cache_history_tasks" "test" {}
```

## Argument Reference

The following arguments are supported:

* `enterprise_project_id` - (Optional, String) Specifies the ID of the enterprise project to which the resource belongs.
  This parameter is valid only when the enterprise project function is enabled. The value **all** indicates all projects.
  This parameter is mandatory when you are an IAM user.
  For enterprise users, if omitted, default enterprise project will be used.

* `status` - (Optional, String) Specifies the task status. Valid values are as follows:
  + **task_inprocess**: The task is being processed.
  + **task_done**: The task is completed.

* `start_date` - (Optional, Int) Specifies the query start time. The value is the number of milliseconds since
  the UNIX epoch (Jan 1, 1970).

* `end_date` - (Optional, Int) Specifies the query end time. The value is the number of milliseconds since
  the UNIX epoch (Jan 1, 1970).

* `order_field` - (Optional, String) Specifies the field used for sorting. Supported fields include
  **task_type** (task type), **total** (total number of URLs), **processing** (number of URLs that are being processed),
  **succeed** (number of processed URLs), **failed** (number of URLs that fail to be processed),
  and **create_time** (task creation time). Both `order_field` and `order_type` must be set together.
  Otherwise, the default values **create_time** and **desc** are used.

* `order_type` - (Optional, String) Specifies the sorting type. Valid values are as follows:
  + **desc**: Descending order.
  + **asc**: Ascending order.

  Defaults to **desc**.

* `file_type` - (Optional, String) Specifies the file type. Possible values: **file** and **directory**.

* `task_type` - (Optional, String) Specifies the task type. Possible values: **refresh** (cache refresh) and
  **preheating** (cache preheat).

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `tasks` - The history task list.
  The [tasks](#CacheHistoryTasks_tasks) structure is documented below.

<a name="CacheHistoryTasks_tasks"></a>
The `tasks` block supports:

* `id` - Indicates the task ID.

* `status` - Indicates the task result. Possible values: **task_done** (task is completed) and
  **task_inprocess** (task is being processed).

* `processing` - Indicates the number of URLs that are being processed.

* `succeed` - Indicates the number of URLs processed.

* `failed` - Indicates the number of URLs that failed to be processed.

* `total` - Indicates the total number of URLs in the task.

* `task_type` - Indicates the task type. Possible values: **refresh** (cache refresh) and **preheating** (cache preheat).

* `created_at` - Indicates the time when the task was created.

* `file_type` - Indicates the file type. Possible values: **file** and **directory**.
