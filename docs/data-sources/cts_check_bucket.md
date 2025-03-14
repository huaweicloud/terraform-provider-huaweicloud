---
subcategory: "Cloud Trace Service (CTS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cts_check_bucket"
description: |-
  Use this data source to check whether data can be transferred to the OBS bucket.
---

# huaweicloud_cts_check_bucket

Use this data source to check whether data can be transferred to the OBS bucket.

## Example Usage

```hcl
data "huaweicloud_cts_check_bucket" "test" {
  bucket_name     = "my_bucket"
  bucket_location = "cn-north-4"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `bucket_name` - (Required, String) Specifies the OBS bucket name.

* `bucket_location` - (Required, String) Specifies the OBS bucket location.

* `is_support_trace_files_encryption` - (Optional, Bool) Specifies whether trace files are encrypted during transfer to
  an OBS bucket.

* `kms_id` - (Optional, String) Specifies the Key ID used for encrypting transferred trace files.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `buckets` - The OBS bucket information.
  The [buckets](#BucketsAttr) structure is documented below.

<a name="BucketsAttr"></a>
The `buckets` block supports:

* `bucket_name` - The OBS bucket name.

* `bucket_location` - The OBS bucket location.

* `is_support_trace_files_encryption` - Whether trace files are encrypted during transfer to an OBS bucket.

* `kms_id` - The Key ID used for encrypting transferred trace files.

* `check_bucket_response` - The check result of the OBS bucket.
  The [check_bucket_response](#CheckBucketResponse) structure is documented below.

<a name="CheckBucketResponse"></a>
The `check_bucket_response` block supports:

* `error_code` - The error code.

* `error_message` - The error message.

* `response_code` - The returned HTTP status code.

* `success` - Whether the transfer is successful.
