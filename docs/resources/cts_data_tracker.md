---
subcategory: "Cloud Trace Service (CTS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cts_data_tracker"
description: |-
  Manages CTS data tracker resource within HuaweiCloud.
---

# huaweicloud_cts_data_tracker

Manages CTS **data** tracker resource within HuaweiCloud.

## Example Usage

```hcl
variable "data_bucket" {}
variable "transfer_bucket" {}

resource "huaweicloud_cts_data_tracker" "tracker" {
  name        = "data-tracker"
  data_bucket = var.data_bucket
  bucket_name = var.transfer_bucket
  file_prefix = "cloudTrace"
  lts_enabled = true
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to manage the CTS data tracker resource.
  If omitted, the provider-level region will be used. Changing this creates a new resource.

* `name` - (Required, String, ForceNew) Specifies the data tracker name. The name cannot be system or ststem-trace.
  Changing this creates a new resource.

* `data_bucket` - (Required, String, ForceNew) Specifies the OBS bucket tracked by the data tracker.
  Changing this creates a new resource.

* `data_operation` - (Optional, List) Specifies an array of operation types tracked by the data tracker,
  the value of operation can be **WRITE** or **READ**.

* `bucket_name` - (Optional, String) Specifies the OBS bucket to which traces will be transferred.

* `file_prefix` - (Optional, String) Specifies the file name prefix to mark trace files that need to be stored
  in an OBS bucket. The value contains `0` to `64` characters. Only letters, numbers, hyphens (-), underscores (_),
  and periods (.) are allowed.

* `obs_retention_period` - (Optional, Int) Specifies the retention period that traces are stored in `bucket_name`,
  the value can be `0`(permanent), `30`, `60`, `90`, `180` or `1,095`.

* `compress_type` - (Optional, String) Specifies the compression type of trace files. The value can be **gzip**
  or **json**. The default value is **gzip**.

* `is_sort_by_service` - (Optional, Bool) Specifies whether to divide the path of the trace file by cloud service.
  The default value is **true**.

* `lts_enabled` - (Optional, Bool) Specifies whether trace analysis is enabled.

* `validate_file` - (Optional, Bool) Specifies whether trace file verification is enabled during trace transfer.

* `enabled` - (Optional, Bool) Specifies whether tracker is enabled.

* `tags` - (Optional, Map) Specifies the key/value pairs to associate with the CTS data tracker.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID in UUID format.
* `type` - The tracker type, only **data** is available.
* `transfer_enabled` - Whether traces will be transferred.
* `status` - The tracker status, the value can be **enabled**, **disabled** or **error**.
* `agency_name` - The cloud service delegation name.
* `create_time` - The creation time of the tracker.
* `detail` - It indicates the cause of the abnormal status.
* `domain_id` - The Account ID.
* `group_id` - The LTS log group ID.
* `stream_id` - The LTS log stream ID.
* `log_group_name` - The name of the log group that CTS creates in LTS.
* `log_topic_name` - The name of the log topic that CTS creates in LTS.
* `is_authorized_bucket` - Whether CTS has been granted permissions to perform operations on the OBS bucket.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 5 minutes.
* `delete` - Default is 5 minutes.

## Import

CTS data tracker can be imported using `name`, e.g.:

```bash
$ terraform import huaweicloud_cts_data_tracker.tracker your_tracker_name
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response, security or some other reason. The missing attribute is `tags`.

It is generally recommended running `terraform plan` after importing the resource.
You can then decide if changes should be applied to the instance, or the resource definition should be updated to
align with the resource. Also you can ignore changes as below.

```hcl
resource "huaweicloud_cts_data_tracker" "test" {
    ...

  lifecycle {
    ignore_changes = [tags]
  }
}
```
