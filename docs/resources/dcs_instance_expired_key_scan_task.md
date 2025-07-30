---
subcategory: "Distributed Cache Service (DCS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dcs_instance_expired_key_scan_task"
description: |-
  Manages a DCS instance expired key scan task resource within HuaweiCloud.
---

# huaweicloud_dcs_instance_expired_key_scan_task

Manages a DCS instance expired key scan task resource within HuaweiCloud.

## Example Usage

```hcl
variable "instance_id" {}

resource "huaweicloud_dcs_instance_expired_key_scan_task" "test" {
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

* `id` - The resource ID.

* `status` - Indicates the status of the expired key scan task.

* `scan_type` - Indicates the scan mode.

* `num` - Indicates the number of expired keys scanned at a time.

* `created_at` - Indicates the time when a scan task is created.

* `started_at` - Indicates the time when a scan task started.

* `finished_at` - Indicates the time when a scan task is complete.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 30 minutes.
