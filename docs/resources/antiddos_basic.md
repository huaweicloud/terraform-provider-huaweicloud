---
subcategory: "Anti-DDoS"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_antiddos_basic"
description: ""
---

# huaweicloud_antiddos_basic

Manages Cloud Native Anti-DDos Basic resource within HuaweiCloud.

-> The Cloud Native Anti-DDos Basic resource will be set to the default traffic cleaning threshold when destroyed,
  instead of deleting it.

## Example Usage

```hcl
variable "eip_id" {}

resource "huaweicloud_antiddos_basic" "antiddos_1" {
  eip_id            = var.eip_id
  traffic_threshold = 150
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to manage the Cloud Native Anti-DDos Basic resource.
  If omitted, the provider-level region will be used. Changing this creates a new resource.

* `eip_id` - (Required, String, ForceNew) Specifies the ID of an EIP. Changing this creates a new resource.

* `traffic_threshold` - (Required, Int) Specifies the traffic cleaning threshold in Mbps.
  The value can be 10, 30, 50, 70, 100, 120, 150, 200, 250, 300, 1000 Mbps.

* `topic_urn` - (Optional, String) Specifies the SMN topic URN. When the value is not empty, it means turning on the alarm
  notification. When the value is empty, it means turning off the alarm notification.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID in UUID format.
* `public_ip` - The public address of the EIP.
* `status` - The Anti-DDos status.

## Timeouts

This resource provides the following timeouts configuration options:

* `update` - Default is 5 minutes.
* `delete` - Default is 5 minutes.

## Import

Cloud Native Anti-DDos Basic resources can be imported using `eip_id`. e.g.

```bash
$ terraform import huaweicloud_antiddos_basic.antiddos_1 c5256d47-8f9e-4ae7-9943-6e77e3d8bd2d
```
