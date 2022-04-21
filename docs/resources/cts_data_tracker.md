---
subcategory: "Cloud Trace Service (CTS)"
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
  in an OBS bucket. The value contains 0 to 64 characters. Only letters, numbers, hyphens (-), underscores (_),
  and periods (.) are allowed.

* `obs_retention_period` - (Optional, Int) Specifies the retention period that traces are stored in `bucket_name`,
  the value can be **0**(permanent), **30**, **60**, **90**, **180** or **1095**.

* `lts_enabled` - (Optional, Bool) Specifies whether trace analysis is enabled.

* `validate_file` - (Optional, Bool) Specifies whether trace file verification is enabled during trace transfer.

* `enabled` - (Optional, Bool) Specifies whether tracker is enabled.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID which equals the tracker name.
* `type` - The tracker type, only **data** is available.
* `transfer_enabled` - Whether traces will be transferred.
* `status` - The tracker status, the value can be **enabled**, **disabled** or **error**.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 5 minute.
* `delete` - Default is 5 minute.

## Import

CTS data tracker can be imported using `name`, e.g.:

```
$ terraform import huaweicloud_cts_data_tracker.tracker your_tracker_name
```
