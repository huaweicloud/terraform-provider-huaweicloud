---
subcategory: "Virtual Private Cloud (VPC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_vpc_traffic_mirror_session"
description: ""
---

# huaweicloud_vpc_traffic_mirror_session

 Manages a VPC traffic mirror session resource within HuaweiCloud.

## Example Usage

```hcl
variable "traffic_mirror_session_name" {}
variable "traffic_mirror_filter_id" {}

resource "huaweicloud_compute_instance" "test" {
  count     = 3
  flavor_id = "c7t.large.2" // currently only instance of c7t flavors can be used as the traffic mirror source
  ...
}

resource "huaweicloud_vpc_traffic_mirror_session" "test" {
  name                     = var.traffic_mirror_session_name
  description              = "Traffic mirror session created by terraform"
  traffic_mirror_filter_id = var.traffic_mirror_filter_id

  traffic_mirror_sources = [
    huaweicloud_compute_instance.test[1].network[0].port,
    huaweicloud_compute_instance.test[2].network[0].port
  ]

  traffic_mirror_target_id   = huaweicloud_compute_instance.test[0].network[0].port
  traffic_mirror_target_type = "eni"
  priority                   = 10
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.
  Changing this creates a new resource.

* `name` - (Required, String) Specifies the name of the traffic mirror session.

* `traffic_mirror_filter_id` - (Required, String) Specifies the traffic mirror filter ID used in the session.

* `traffic_mirror_sources` - (Required, List) Specifies the mirror source IDs.
  An elastic network interface can be used as a mirror source.

* `traffic_mirror_target_id` - (Required, String) Specifies the mirror target ID.

* `traffic_mirror_target_type` - (Required, String) Specifies the mirror target type. The value can be:
  + **eni**: elastic network interface;
  + **elb**: private network load balancer;

* `priority` - (Required, Int) Specifies the mirror session priority. The value range is **1-32766**.
  A smaller value indicates a higher priority.

* `description` - (Optional, String) Specifies the description of the traffic mirror session.

* `enabled` - (Optional, Bool) Specifies whether the mirror session is enabled. Defaults to **true**.

* `type` - (Optional, String) Specifies the mirror source type. The value can be **eni**(elastic network interface).

* `virtual_network_id` - (Optional, Int) Specifies the VNI, which is used to distinguish mirrored traffic of
  different sessions. The valid value is range from `0` to `16,777,215`, defaults to `1`.

* `packet_length` - (Optional, Int) Specifies the maximum transmission unit (MTU).
 The valid value is range from `1` to `1,460`, defaults to `96`.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `created_at` - The creation time of the traffic mirror session.

* `updated_at` - The latest update time of the traffic mirror session.

## Import

The traffic mirror session can be imported using `id`, e.g.

```bash
$ terraform import huaweicloud_vpc_traffic_mirror_session.test <id>
```
