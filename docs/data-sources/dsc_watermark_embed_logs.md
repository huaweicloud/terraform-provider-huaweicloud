---
subcategory: "Data Security Center (DSC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dsc_watermark_embed_logs"
description: |-
  Use this data source to get the list of watermark injection task logs within HuaweiCloud.
---

# huaweicloud_dsc_watermark_embed_logs

Use this data source to get the list of watermark injection task logs within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_dsc_watermark_embed_logs" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `water_mark_embed_logs` - The list of watermark injection task logs.

  The [water_mark_embed_logs](#water_mark_embed_logs_struct) structure is documented below.

<a name="water_mark_embed_logs_struct"></a>
The `water_mark_embed_logs` block supports:

* `blind_watermark` - The blind watermark content.

* `dest_url` - The destination URL.

* `doc_path` - The document path.

* `download_url` - The download URL.

* `file_exist` - Whether the file exists.

* `file_url` - The file URL.

* `finish_num` - The number of finished files.

* `project_id` - The project ID.

* `remark` - The remark.

* `state` - The task status.

* `task_id` - The task ID.

* `task_name` - The task name.

* `total_file_num` - The total number of files.

* `visible_watermark` - The visible watermark content.
