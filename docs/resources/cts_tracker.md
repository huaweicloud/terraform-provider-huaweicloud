---
subcategory: "Cloud Trace Service (CTS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cts_tracker"
description: |-
  Manages CTS system tracker resource within HuaweiCloud.
---

# huaweicloud_cts_tracker

Manages CTS system tracker resource within HuaweiCloud.

## Example Usage

```hcl
variable "bucket_name" {}

resource "huaweicloud_cts_tracker" "tracker" {
  bucket_name = var.bucket_name
  file_prefix = "cts"
  lts_enabled = true
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to manage the CTS system tracker resource.
  If omitted, the provider-level region will be used. Changing this creates a new resource.

* `bucket_name` - (Optional, String) Specifies the OBS bucket to which traces will be transferred.

* `file_prefix` - (Optional, String) Specifies the file name prefix to mark trace files that need to be stored
  in an OBS bucket. The value contains `0` to `64` characters. Only letters, numbers, hyphens (-), underscores (_),
  and periods (.) are allowed.

* `lts_enabled` - (Optional, Bool) Specifies whether trace analysis is enabled.

* `organization_enabled` - (Optional, Bool) Specifies whether to apply the tracker configuration to the organization.
  If the value is set to **true**, the audit logs of all members in the organization in the current region will be
  transferred to the OBS bucket or LTS log stream configured for the management tracker.

* `validate_file` - (Optional, Bool) Specifies whether trace file verification is enabled during trace transfer.

* `kms_id` - (Optional, String) Specifies the ID of KMS key used for trace file encryption.

* `compress_type` - (Optional, String) Specifies the compression type of trace files. The value can be
  **gzip** or **json**. The default value is **gzip**.

* `is_sort_by_service` - (Optional, Bool) Specifies whether to divide the path of the trace file by cloud service.
  The default value is **true**.

* `exclude_service` - (Optional, List) Specifies the names of the cloud services for which traces don't need to be transferred.
  Currently, only **KMS** is supported.

* `enabled` - (Optional, Bool) Specifies whether tracker is enabled.

* `tags` - (Optional, Map) Specifies the key/value pairs to associate with the CTS tracker.

* `delete_tracker` - (Optional, Bool) Specifies whether the tracker can be deleted.
  
  -> **NOTE:** By default, resource deletion only clears parameters without removing the system tracker.
  To delete the system tracker, set `delete_tracker` to true. Note that this will disable the CTS service.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID in UUID format.
* `name` - The tracker name, only **system** is available.
* `type` - The tracker type, only **system** is available.
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

CTS tracker can be imported using `name`, only **system** is available. e.g.

```bash
$ terraform import huaweicloud_cts_tracker.tracker system
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response, security or some other reason. The missing attributes include: `tags`, `delete_tracker`.

It is generally recommended running `terraform plan` after importing the resource.
You can then decide if changes should be applied to the instance, or the resource definition should be updated to
align with the resource. Also you can ignore changes as below.

```hcl
resource "huaweicloud_cts_tracker" "test" {
    ...

  lifecycle {
    ignore_changes = [
      tags, delete_tracker
    ]
  }
}
```
