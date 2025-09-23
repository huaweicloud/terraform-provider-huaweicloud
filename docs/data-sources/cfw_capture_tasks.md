---
subcategory: "Cloud Firewall (CFW)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cfw_capture_tasks"
description: |-
  Use this data source to get the list of CFW capture tasks.
---

# huaweicloud_cfw_capture_tasks

Use this data source to get the list of CFW capture tasks.

## Example Usage

```hcl
variable "fw_instance_id" {}

data "huaweicloud_cfw_capture_tasks" "test" {
  fw_instance_id = var.fw_instance_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `fw_instance_id` - (Required, String) Specifies the ID of the firewall instance.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `records` - List of capture task information.

  The [records](#data_records_struct) structure is documented below.

<a name="data_records_struct"></a>
The `records` block supports:

* `task_id` - The capture task ID.

* `name` - The capture task name.

* `status` - The capture task status.

* `source_address_type` - The source address type.

* `source_address` - The source address.

* `dest_address_type` - The destination address type.

* `dest_address` - The destination address.

* `protocol` - The protocol type.

* `source_port` - The source port.

* `dest_port` - The destination port.

* `duration` - The capture task duration.

* `remaining_days` - The remaining days.

* `is_deleted` - Whether is deleted.

* `max_packets` - The max packets.

* `capture_size` - The capture task size.

* `created_at` - The creation time of the capture task.

* `updated_at` - The update time of the capture task.
