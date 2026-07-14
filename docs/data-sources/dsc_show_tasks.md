---
subcategory: "Data Security Center (DSC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dsc_show_tasks"
description: |-
  Use this data source to get the task execution statistics of the current project within HuaweiCloud.
---

# huaweicloud_dsc_show_tasks

Use this data source to get the task execution statistics of the current project within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_dsc_show_tasks" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `mask_task_num` - The total number of static desensitization tasks.

* `scan_task_num` - The total number of sensitive data identification tasks.

* `watermark_embed_num` - The number of watermark embedding executions.

* `watermark_extract_num` - The number of watermark extraction executions.

* `watermark_task_num` - The total number of watermark-related tasks.
