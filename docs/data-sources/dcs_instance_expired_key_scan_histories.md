---
subcategory: "Distributed Cache Service (DCS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dcs_instance_expired_key_scan_histories"
description: |-
  Use this data source to get the list of expired key scan records.
---

# huaweicloud_dcs_instance_expired_key_scan_histories

Use this data source to get the list of expired key scan records.

## Example Usage

```hcl
variable "instance_id"{}

data "huaweicloud_dcs_instance_expire_key_scan_histories" "test" {
  instance_id = var.instance_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the ID of the DCS instance.

* `status` - (Optional, String) Specifies the status of the expired key scan task.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `records` - Indicates the expired key record.

  The [records](#records_struct) structure is documented below.

<a name="records_struct"></a>
The `records` block supports:

* `instance_id` - Indicates the instance ID.

* `id` - Indicates the ID of the expired key scan task.

* `status` - Indicates the status of the expired key scan task.
  The value can be:
  + **waiting**: The task is to be processed.
  + **running**: The task is being processed.
  + **success**: The task is successful.
  + **failed**: The task fails.

* `scan_type` - Indicates the scan mode.
  The value can be:
  + **auto**: automatic scan
  + **manual**: manual scan

* `num` - Indicates the number of expired keys scanned at a time.

* `created_at` - Indicates the time when a scan task is created.

* `started_at` - Indicates the time when a scan task started.

* `finished_at` - Indicates the time when a scan task is complete.
