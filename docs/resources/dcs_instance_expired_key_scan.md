---
subcategory: "Distributed Cache Service (DCS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dcs_instance_expired_key_scan"
description: |-
  Manages a DCS instance expired key scan resource within HuaweiCloud.
---

# huaweicloud_dcs_instance_expired_key_scan

Manages a DCS instance expired key scan resource within HuaweiCloud.

## Example Usage

```hcl
variable "instance_id" {}

resource "huaweicloud_dcs_instance_expired_key_scan" "test"{
  instance_id = var.instance_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `instance_id` - (Required, String, NonUpdatable) Specifies the ID of the DCS instance.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID which equals the `instance_id`.

* `status` - Indicates the status of the expired key scan task. The value can be:
  + **waiting**: The task is to be processed.
  + **running**: The task is being processed.
  + **success**: The task is successful.
  + **failed**: The task fails.

* `scan_type` - Indicates the can mode. The value can be:
  + **manual**: manual scan.
  + **auto**: automatic scan.

* `num` - Indicates the number of expired keys scanned at a time.

* `created_at` - Indicates the time when a scan task is created. The format is **yyyy-mm-dd hh:mm:ss**.
  The value is in UTC format.

* `started_at` - Indicates the time when a scan task started. The format is **yyyy-mm-dd hh:mm:ss**.
  The value is in UTC format.

* `finished_at` - Indicates the time when a scan task is complete. The format is **yyyy-mm-dd hh:mm:ss**.
  The value is in UTC format.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 30 minutes.
