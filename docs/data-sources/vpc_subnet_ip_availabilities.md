---
subcategory: "Virtual Private Cloud (VPC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_vpc_subnet_ip_availabilities"
description: |-
  Use this data source to get IP availabilities of a subnet.
---

# huaweicloud_vpc_subnet_ip_availabilities

Use this data source to get IP availabilities of a subnet.

## Example Usage

```hcl
variable "network_id" {}
data "huaweicloud_vpc_subnet_ip_availabilities" "test" {
  network_id = var.network_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `network_id` - (Required, String) Specifies the subnet ID.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `network_ip_availability` - The network IP address usage objects.

  The [network_ip_availability](#network_ip_availability_struct) structure is documented below.

<a name="network_ip_availability_struct"></a>
The `network_ip_availability` block supports:

* `network_id` - The network ID.

* `network_name` - The network name.

* `total_ips` - The total number of IP addresses on a network.
  The reserved IP addresses are not included.

* `used_ips` - The number of in-use IP addresses on a network.
  The reserved IP addresses are not included.

* `subnet_ip_availability` - The subnet IP address usage objects.

  The [subnet_ip_availability](#network_ip_availability_subnet_ip_availability_struct) structure is documented below.

<a name="network_ip_availability_subnet_ip_availability_struct"></a>
The `subnet_ip_availability` block supports:

* `used_ips` - The number of in-use IP addresses on a subnet.
  The reserved IP addresses are not included.

* `subnet_id` - The subnet ID.

* `subnet_name` - The subnet name.

* `ip_version` - The IP version of the subnet.
  The value can be **4** or **6**.

* `cidr` - The subnet CIDR block.

* `total_ips` - The total number of IP addresses on a subnet.
  The reserved IP addresses are not included.
