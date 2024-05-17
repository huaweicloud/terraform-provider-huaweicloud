---
subcategory: "Virtual Private Cloud (VPC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_vpc_traffic_mirror_filter_rules"
description: |-
  Use this data source to get the traffic mirror filter rules.
---

# huaweicloud_vpc_traffic_mirror_filter_rules

Use this data source to get the traffic mirror filter rules.

## Example Usage

```hcl
data "huaweicloud_vpc_traffic_mirror_filter_rules" "test1" {
  protocol = "udp"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `traffic_mirror_filter_rule_id` - (Optional, String) The traffic mirror filter rule ID used as the query filter.

* `traffic_mirror_filter_id` - (Optional, String) The traffic mirror filter ID used as the query filter.

* `direction` - (Optional, String) The direction of the traffic mirror filter rule.
  Valid values are **ingress** or **egress**.

* `protocol` - (Optional, String) The protocol of the traffic mirror filter rule.
  Valid value are **tcp**, **udp**, **icmp**, **icmpv6**, **all**.

* `action` - (Optional, String) The policy of in the traffic mirror filter rule.
  Valid values are **accept** or **reject**.

* `priority` - (Optional, String) The priority number of the traffic mirror filter rule.
  Valid value ranges from `1` to `65,535`.

* `source_cidr_block` - (Optional, String) The source IP address of the traffic mirror filter rule.

* `source_port_range` - (Optional, String) The source port number range of the traffic mirror filter rule.
  The value ranges from `1` to `65,535`, enter two port numbers connected by a hyphen (-). For example, **80-200**.

* `destination_cidr_block` - (Optional, String) The destination IP address of the traffic mirror filter rule.

* `destination_port_range` - (Optional, String) The destination port number range of the traffic mirror filter rule.
  The value ranges from `1` to `65,535`, enter two port numbers connected by a hyphen (-). For example, **80-200**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `traffic_mirror_filter_rules` - List of traffic mirror filter rules.

  The [traffic_mirror_filter_rules](#traffic_mirror_filter_rules_struct) structure is documented below.

<a name="traffic_mirror_filter_rules_struct"></a>
The `traffic_mirror_filter_rules` block supports:

* `created_at` - Time when a traffic mirror filter rule is created.

* `source_cidr_block` - Source CIDR block of the mirrored traffic.

* `project_id` - Project ID.

* `id` - Traffic mirror filter rule ID.

* `action` - Whether to accept or reject traffic.

* `source_port_range` - Source port range.

* `destination_cidr_block` - Destination CIDR block of the mirrored traffic.

* `traffic_mirror_filter_id` - Traffic mirror filter ID.

* `description` - Description of a traffic mirror filter rule.

* `updated_at` - Time when a traffic mirror filter rule is updated.

* `priority` - Mirror filter rule priority.

* `ethertype` - IP address version of the mirrored traffic.

* `destination_port_range` - Source port range.

* `protocol` - Protocol of the mirrored traffic.

* `direction` - Traffic direction.
