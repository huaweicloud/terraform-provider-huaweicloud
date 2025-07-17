---
subcategory: "Server Migration Service (SMS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_sms_task_network_check_info_report"
description: |-
  Manages an SMS update task network check info resource within HuaweiCloud.
---

# huaweicloud_sms_task_network_check_info_report

Manages an SMS update task network check info resource within HuaweiCloud.

~> Deleting update task network check info resource is not supported, it will only be removed from the state.

## Example Usage

```hcl
variable "task_id" {}
variable "evaluation_result" {}

resource "huaweicloud_sms_task_network_check_info_report" "test" {
  task_id           = var.task_id
  network_delay     = 20.0
  network_jitter    = 2.0
  migration_speed   = 100.0
  loss_percentage   = 0.0
  cpu_usage         = 20.0
  mem_usage         = 20.0
  evaluation_result = var.evaluation_result
}
```

## Argument Reference

The following arguments are supported:

* `task_id` - (Required, String, NonUpdatable) Specifies the task ID.

* `network_delay` - (Required, Float, NonUpdatable) Specifies the network latency.

* `network_jitter` - (Required, Float, NonUpdatable) Specifies the network jitter.

* `migration_speed` - (Required, Float, NonUpdatable) Specifies the bandwidth.

* `loss_percentage` - (Required, Float, NonUpdatable) Specifies the packet loss rate.

* `cpu_usage` - (Required, Float, NonUpdatable) Specifies the CPU usage.

* `mem_usage` - (Required, Float, NonUpdatable) Specifies the memory usage.

* `evaluation_result` - (Required, String, NonUpdatable) Specifies the network evaluation result.

* `domain_connectivity` - (Optional, Bool, NonUpdatable) Specifies the connectivity to domain names.

* `destination_connectivity` - (Optional, Bool, NonUpdatable) Specifies the connectivity to the target server.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID, which equals to `task_id`.
