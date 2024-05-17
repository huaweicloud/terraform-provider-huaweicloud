---
subcategory: "Virtual Private Cloud (VPC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_vpc_traffic_mirror_filters"
description: |-
  Use this data source to get the traffic mirror filters.
---

# huaweicloud_vpc_traffic_mirror_filters

Use this data source to get the traffic mirror filters.

## Example Usage

### query traffic mirror filter list

```hcl
data "huaweicloud_vpc_traffic_mirror_filters" "filter_test1" {
}
```

### query traffic mirror filter by name

```hcl
data "huaweicloud_vpc_traffic_mirror_filters" "filter_test2" {
  name = "test-filter"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `traffic_mirror_filter_id` - (Optional, String) Specifies the ID of the traffic mirror filter.

* `name` - (Optional, String) Specifies the name of the traffic mirror filter.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `traffic_mirror_filters` - List of traffic mirror filters.

  The [traffic_mirror_filters](#traffic_mirror_filters_struct) structure is documented below.

<a name="traffic_mirror_filters_struct"></a>
The `traffic_mirror_filters` block supports:

* `id` - Traffic mirror filter ID.

* `project_id` - Project ID.

* `description` - Description of a traffic mirror filter.

* `name` - Traffic mirror filter name.

* `updated_at` - Time when a traffic mirror filter is updated.

* `created_at` - Time when a traffic mirror filter is created.

* `egress_rules` - Outbound mirror filter rules.

  The [egress_rules](#traffic_mirror_filters_egress_rules_struct) structure is documented below.

* `ingress_rules` - Inbound mirror filter rules.

  The [ingress_rules](#traffic_mirror_filters_ingress_rules_struct) structure is documented below.

<a name="traffic_mirror_filters_egress_rules_struct"></a>
The `egress_rules` block supports:

* `id` - Traffic mirror filter rule ID.

* `project_id` - Project ID.

* `direction` - Traffic direction.

* `description` - Description of a traffic mirror filter rule.

* `action` - Whether to accept or reject traffic.

* `traffic_mirror_filter_id` - Traffic mirror filter ID.

* `source_cidr_block` - Source CIDR block of the mirrored traffic.

* `source_port_range` - Source port range.

* `destination_cidr_block` - Destination CIDR block of the mirrored traffic.

* `destination_port_range` - Source port range.

* `ethertype` - IP address version of the mirrored traffic.

* `protocol` - Protocol of the mirrored traffic.

* `priority` - Mirror filter rule priority.

* `created_at` - Time when a traffic mirror filter rule is created.

* `updated_at` - Time when a traffic mirror filter rule is updated.

<a name="traffic_mirror_filters_ingress_rules_struct"></a>
The `ingress_rules` block supports:

* `id` - Traffic mirror filter rule ID.

* `project_id` - Project ID.

* `description` - Description of a traffic mirror filter rule.

* `traffic_mirror_filter_id` - Traffic mirror filter ID.

* `direction` - Traffic direction.

* `source_cidr_block` - Source CIDR block of the mirrored traffic.

* `source_port_range` - Source port range.

* `destination_cidr_block` - Destination CIDR block of the mirrored traffic.

* `destination_port_range` - Source port range.

* `ethertype` - IP address version of the mirrored traffic.

* `protocol` - Protocol of the mirrored traffic.

* `action` - Whether to accept or reject traffic.

* `priority` - Mirror filter rule priority.

* `created_at` - Time when a traffic mirror filter rule is created.

* `updated_at` - Time when a traffic mirror filter rule is updated.
