---
subcategory: "Cloud Firewall (CFW)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cfw_capture_task"
description: |-
  Manages a CFW capture task resource within HuaweiCloud.
---

# huaweicloud_cfw_capture_task

Manages a CFW capture task resource within HuaweiCloud.

-> **NOTE:** For the Cloud Firewall service, you can only initiate up to 20 packet capture tasks per day.
Beyond this limit, no additional packet capture tasks can be initiated. Furthermore, only one packet capture task can be
in progress at any given time.

## Example Usage

```hcl
variable "fw_instance_id" {}
variable "name" {}
variable "duration" {}
variable "max_packets" {}

resource huaweicloud_cfw_capture_task "test" {
  fw_instance_id = var.fw_instance_id
  name           = var.name
  duration       = var.duration
  max_packets    = var.max_packets
  
  destination {
    address      = "1.1.1.1"
    address_type = 0
  }

  source {
    address      = "2.2.2.2"
    address_type = 0
  }

  service {
    protocol = -1
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.
  Changing this creates a new resource.

* `fw_instance_id` - (Required, String, NonUpdatable) Specifies the ID of the firewall instance.

* `name` - (Required, String, NonUpdatable) Specifies the capture task name.

* `duration` - (Required, Int, NonUpdatable) Specifies the capture task duration.

* `max_packets` - (Required, Int, NonUpdatable) Specifies the maximum number of packets captured.
  The Maximum value is `1,000,000`.

* `destination` - (Required, List, NonUpdatable) Specifies the destination configuration.
  The [destination](#Address) structure is documented below.

* `source` - (Required, List, NonUpdatable) Specifies the source configuration.
  The [source](#Address) structure is documented below.

* `service` - (Required, List, NonUpdatable) Specifies the service configuration.
  The [service](#Service) structure is documented below.

* `stop_capture` - (Optional, Bool) Specifies whether to stop the capture task.

<a name="Address"></a>
The `destination` or `source` block supports:

* `address` - (Required, String, NonUpdatable) Specifies the address.

* `address_type` - (Required, Int, NonUpdatable) Specifies the address type.
  The valid values are:
  + **0**: indicates IPv4;
  + **1**: indicates IPv6.

<a name="Service"></a>
The `service` block supports:

* `protocol` - (Required, Int, NonUpdatable) Specifies the protocol type.
  The valid values are:
  + **6**: indicates TCP;
  + **17**: indicates UDP;
  + **1**: indicates ICMP;
  + **58**: indicates ICMPv6;
  + **-1**: indicates any protocol.

* `dest_port` - (Optional, String, NonUpdatable) Specifies the destination port.

* `source_port` - (Optional, String, NonUpdatable) Specifies the source port.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID. The value is the name of the capture task.

* `created_at` - The creation time of the capture task.

* `status` - The status of the capture task.

* `updated_at` - The update time of the capture task.

* `task_id` - The ID of the capture task.

## Timeouts

This resource provides the following timeouts configuration options:

* `update` - Default is 5 minutes.
* `delete` - Default is 5 minutes.

## Import

The capture task can be imported using `fw_instance_id`, `name`, separated by a slash, e.g.

```bash
$ terraform import huaweicloud_cfw_capture_task.test <fw_instance_id>/<name>
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response, security or some other reason.
The missing attributes is `stop_capture`. It is generally recommended running `terraform plan` after importing the resource.
You can then decide if changes should be applied to the capture task, or the resource definition should be updated to
align with the capture task. Also you can ignore changes as below.

```hcl
resource "huaweicloud_cfw_capture_task" "test" {
    ...

  lifecycle {
    ignore_changes = [
      stop_capture,
    ]
  }
}
```
