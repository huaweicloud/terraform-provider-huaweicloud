---
subcategory: "Virtual Private Cloud (VPC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_vpc_traffic_mirror_filter_rule"
description: ""
---

# huaweicloud_vpc_traffic_mirror_filter_rule

Manages a VPC traffic mirror filter rule resource within HuaweiCloud.

## Example Usage

```hcl
variable "traffic_mirror_filter_id" {}

resource "huaweicloud_vpc_traffic_mirror_filter_rule" "test" {
  traffic_mirror_filter_id = var.traffic_mirror_filter_id
  direction                = "ingress"
  protocol                 = "tcp"
  ethertype                = "IPv4"
  action                   = "accept"
  priority                 = 1
  source_cidr_block        = "10.0.0.0/8"
  source_port_range        = "80-90"
  destination_cidr_block   = "192.168.1.0/24"
  destination_port_range   = "10-65535"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.
  Changing this creates a new resource.

* `traffic_mirror_filter_id` - (Required, String, ForceNew) Specifies an ID of the traffic mirror filter to which
  the rule belongs. Changing this creates a new resource.

* `direction` - (Required, String, ForceNew) Specifies the direction of the traffic mirror filter rule.
  Valid values are **ingress** or **egress**. Changing this creates a new resource.

* `protocol` - (Required, String) Specifies the protocol of the traffic mirror filter rule.
  Valid value are **tcp**, **udp**, **icmp**, **icmpv6**, **all**.

* `ethertype` - (Required, String) Specifies the IP address protocol type of the traffic mirror filter rule.
  Valid values are **IPv4** or **IPv6**.

* `action` - (Required, String) Specifies the policy of in the traffic mirror filter rule.
  Valid values are **accept** or **reject**.

* `priority` - (Required, Int) Specifies the priority number of the traffic mirror filter rule.
  Valid value ranges from `1` to `65,535`.
  The smaller the priority number, the higher the priority.

* `source_cidr_block` - (Optional, String) Specifies the source IP address of the traffic traffic mirror filter rule.

* `source_port_range` - (Optional, String) Specifies the source port number range of the traffic mirror filter rule.
  The value ranges from `1` to `65,535`, enter two port numbers connected by a hyphen (-). For example, **80-200**.

* `destination_cidr_block` - (Optional, String) Specifies the destination IP address of the traffic traffic mirror filter
  rule.

* `destination_port_range` - (Optional, String) Specifies the destination port number range of the traffic mirror filter
  rule. The value ranges from `1` to `65,535`, enter two port numbers connected by a hyphen (-). For example, **80-200**.

* `description` - (Optional, String) Specifies the description of the traffic mirror filter rule.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `created_at` - The creation time of the traffic mirror filter rule.

* `updated_at` - The latest update time of the traffic mirror filter rule.

## Import

The traffic mirror filter rule can be imported using `id`, e.g.

```bash
$ terraform import huaweicloud_vpc_traffic_mirror_filter_rule.test <id>
```
