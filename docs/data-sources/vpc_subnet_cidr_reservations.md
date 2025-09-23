---
subcategory: "Virtual Private Cloud (VPC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_vpc_subnet_cidr_reservations"
description: |-
  Use this data source to get a list of VPC subnet CIDR reservations.
---

# huaweicloud_vpc_subnet_cidr_reservations

Use this data source to get a list of VPC subnet CIDR reservations.

## Example Usage

```hcl
variable "reservation_id_1" {}
variable "reservation_id_2" {}

data "huaweicloud_vpc_subnet_cidr_reservations" "test" {
  reservation_id = [var.reservation_id_1, var.reservation_id_2]
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `reservation_id` - (Optional, List) Specifies the subnet CIDR reservation resource IDs.
  Multiple IDs are supported for filtering.

* `subnet_id` - (Optional, List) Specifies the IDs of the subnets containing the CIDR reservations.
  Multiple IDs supported for filtering.

* `cidr` - (Optional, List) Specifies the CIDRs of the subnet reservations. Multiple CIDRs supported for filtering.

* `ip_version` - (Optional, List) Specifies the IP versions of the subnets. Multiple versions supported for filtering.

* `name` - (Optional, List) Specifies the Names of the subnet CIDR reservations. Multiple names supported for filtering.

* `description` - (Optional, List) Specifies the Descriptions of the subnet CIDR reservations.
  Multiple descriptions supported for filtering.

* `enterprise_project_id` - (Optional, String) Specifies the Enterprise project ID.
  Used to filter reservations within a specific enterprise project.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `reservations` - The list of VPC subnet CIDR reservations.

  The [reservations](#reservations) structure is documented below.

<a name="reservations"></a>
The `reservations` block supports:

* `id` - Indicates the ID of the subnet CIDR reservation.

* `subnet_id` - Indicates the ID of the subnet containing the CIDR reservation.

* `vpc_id` - Indicates the ID of the VPC containing the subnet.

* `ip_version` - Indicates the IP version of the subnet CIDR reservation.

* `cidr` - Indicates the CIDR of the subnet reservation.

* `name` - Indicates the name of the subnet CIDR reservation.

* `description` - Indicates the description of the subnet CIDR reservation.

* `project_id` - Indicates the project ID to which the CIDR reservation belongs.

* `created_at` - Indicates the creation time of the CIDR reservation.

* `updated_at` - Indicates the last update time of the CIDR reservation.
