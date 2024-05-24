---
subcategory: "Cloud Trace Service (CTS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cts_trackers"
description: |-
  Use this data source to get the list of CTS trackers.
---

# huaweicloud_cts_trackers

Use this data source to get the list of CTS trackers.

## Example Usage

```hcl
variable "tracker_type" {}

data "huaweicloud_cts_trackers" "test" {
  type = var.tracker_type
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `tracker_id` - (Optional, String) Specifies the tracker ID.

* `name` - (Optional, String) Specifies the tracker name.
  If this parameter is not specified, all trackers of a tenant will be queried.

* `type` - (Optional, String) Specifies the tracker type.
  The value can be **system** (management tracker) or **data** (data tracker).

* `status` - (Optional, String) Specifies the tracker status.
  The valid values are **enabled**, **disabled** and **error**.

* `data_bucket_name` - (Optional, String) Specifies the data bucket name.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `trackers` - List of tracker information.

  The [trackers](#trackers_struct) structure is documented below.

<a name="trackers_struct"></a>
The `trackers` block supports:

* `id` - The unique tracker ID.

* `name` - The tracker name.

* `type` - The tracker type.

* `create_time` - The time when the tracker was created. The time is in UTC.

* `status` - The tracker status.

* `kms_id` - The ID of the key used for trace file encryption.

* `domain_id` - The ID of the account that the tracker belongs to.

* `project_id` - The project ID.

* `detail` - This parameter is returned only when the tracker status is **error**.

* `group_id` - The ID of the LTS log group.

* `stream_id` - The ID of the LTS log stream.

* `data_bucket` - Information about the bucket tracked by a data tracker.

  The [data_bucket](#trackers_data_bucket_struct) structure is documented below.

* `obs_info` - Information about the bucket to which traces are transferred.

  The [obs_info](#trackers_obs_info_struct) structure is documented below.

* `is_support_validate` - Whether trace file verification is enabled.

* `lts` - The LTS configuration.

  The [lts](#trackers_lts_struct) structure is documented below.

* `is_support_trace_files_encryption` - Whether trace files are encrypted during transfer to an OBS bucket.

* `is_organization_tracker` - Whether system tracker to apply to my organization.

* `management_event_selector` - The management event selector.

  The [management_event_selector](#trackers_management_event_selector_struct) structure is documented below.

* `agency_name` - The name of a cloud service agency.

<a name="trackers_data_bucket_struct"></a>
The `data_bucket` block supports:

* `data_bucket_name` - The OBS bucket name.

* `search_enabled` - Whether the logs of the tracked bucket can be searched.

* `data_event` - The list of the bucket event operation types.

<a name="trackers_obs_info_struct"></a>
The `obs_info` block supports:

* `bucket_name` - The OBS bucket name.

* `file_prefix_name` - File name prefix to mark trace files that need to be stored in an OBS bucket.

* `is_obs_created` - Whether the OBS bucket is automatically created by the tracker.

* `is_authorized_bucket` - Whether CTS has been granted permissions to perform operations on the OBS bucket.

* `bucket_lifecycle` - Duration that traces are stored in the OBS bucket.

* `compress_type` - The compression type.

* `is_sort_by_service` - Whether to sort the path by cloud service.

<a name="trackers_lts_struct"></a>
The `lts` block supports:

* `is_lts_enabled` - Whether traces are synchronized to LTS for trace search and analysis.

* `log_group_name` - The name of the log group that CTS creates in LTS.

* `log_topic_name` - The name of the log stream that CTS creates in LTS.

<a name="trackers_management_event_selector_struct"></a>
The `management_event_selector` block supports:

* `exclude_service` - The cloud service that is not dumped.
